package storage

import (
	"encoding/binary"
	"errors"
)

type Slot struct {
	Offset uint16 // 02 bytes
	Length uint16 // 02 bytes
}

const (
	OffsetOffset = 2
	OffsetLenght = 2
	SlotSize     = OffsetOffset + OffsetLenght
)

var (
	ErrNotSize      = errors.New("not enough space")
	ErrInvalidSlot  = errors.New("invalid slot")
	ErrDeleteRecord = errors.New("deleted record")
)

type SlottedPage struct {
	Header PageHeader
	Data []byte
}

func (p *SlottedPage) numSlots() int {
	return int(p.Header.SlotCount)
}

func (p *SlottedPage) readSlot(slotID uint16) Slot {
	pos := int(slotID) * SlotSize

	return Slot{
		Offset: binary.LittleEndian.Uint16(p.Data[OffsetData+pos:]),
		Length: binary.LittleEndian.Uint16(p.Data[OffsetData+pos+2:]),
	}
}

func (p *SlottedPage) writeSlot(slotID uint16, s Slot) {
	pos := int(slotID) * SlotSize
	binary.LittleEndian.PutUint16(p.Data[OffsetData+pos:], s.Offset)
	binary.LittleEndian.PutUint16(p.Data[OffsetData+pos+2:], s.Length)
}

func NewSlottedPage(pageID uint64) *SlottedPage {
	slottedPage := &SlottedPage{
		Header: PageHeader{
			Magic:     PageMagic,
			Version:   1,
			PageID:    pageID,
			FreeStart: PageHeaderSize,
			FreeEnd:   PageSize,
		},
		Data: make([]byte, PageSize),
	}
	return slottedPage
}

func (p *SlottedPage) compact() {

	newUpper := uint16(len(p.Data))
	newSlotId := 0

	for id := 0; id < p.numSlots(); id++ {
		slot := p.readSlot(uint16(id))
		if slot.Length == 0 {
			continue
		}
		newUpper -= slot.Length

		copy(
			p.Data[newUpper:newUpper+slot.Length],
			p.Data[slot.Offset:slot.Offset+slot.Length],
			)

		slot.Offset = newUpper
		p.writeSlot(uint16(newSlotId), slot)
		newSlotId += 1

	}
	p.Header.FreeEnd = newUpper
	p.Header.SlotCount = uint16(newSlotId)
	p.Header.FreeStart = PageHeaderSize + (p.Header.SlotCount * SlotSize)
}

func (p *SlottedPage) Insert(record []byte) (uint16, error) {
	recordSize := uint16(len(record))
	slotSize := uint16(binary.Size(Slot{}))
	prevSlot := Slot{
		Offset: PageSize,
	}

	for id := 0; id < p.numSlots(); id++ {
		readSlot := p.readSlot(uint16(id))
		if prevSlot.Offset-readSlot.Offset >= recordSize && readSlot.Length == 0 {
			copy(p.Data[readSlot.Offset:], record)
			slot := Slot{
				Offset: readSlot.Offset,
				// Length: s.Length,
				Length: recordSize,
			}
			p.writeSlot(uint16(id), slot)
			return uint16(id), nil
		}
		prevSlot = readSlot
	}

	available := p.Header.FreeEnd - p.Header.FreeStart
	if recordSize < available {
		p.Header.FreeEnd -= recordSize
		copy(p.Data[p.Header.FreeEnd:], record)

		slot := Slot{
			Offset: p.Header.FreeEnd,
			Length: recordSize,
		}

		slotOffset := p.Header.FreeStart
		binary.LittleEndian.PutUint16(p.Data[slotOffset:], slot.Offset)
		binary.LittleEndian.PutUint16(p.Data[slotOffset+2:], slot.Length)

		p.Header.FreeStart += slotSize
		slotID := p.Header.SlotCount
		p.Header.SlotCount++

		return slotID, nil
	}

	if recordSize > available {
		p.compact()
		if recordSize > (p.Header.FreeEnd - p.Header.FreeStart) {
			return 0, errors.New("page full")
		}
		sloteID, error := p.Insert(record)
		return sloteID, error
	}

	return 0, errors.New("page full")
}

func (p *SlottedPage) Get(slotID uint16) ([]byte, error) {
	if slotID >= p.Header.SlotCount {
		return nil, ErrInvalidSlot
	}

	slotOffset := PageHeaderSize + slotID*4
	offset := binary.LittleEndian.Uint16(p.Data[slotOffset:])
	length := binary.LittleEndian.Uint16(p.Data[slotOffset+2:])

	if length == 0 {
		return nil, ErrDeleteRecord
	}

	record := make([]byte, length)
	copy(record, p.Data[offset:offset+length])

	return record, nil
}

func (p *SlottedPage) Delete(slotID uint16) error {
	if slotID >= p.Header.SlotCount {
		return ErrInvalidSlot
	}

	slotOffset := PageHeaderSize + slotID*4
	binary.LittleEndian.PutUint16(p.Data[slotOffset+2:], 0)

	return nil
}
