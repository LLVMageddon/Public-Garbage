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
)

var (
	ErrNotSize      = errors.New("not enough space")
	ErrInvalidSlot  = errors.New("invalid slot")
	ErrDeleteRecord = errors.New("deleted record")
)

type SlottedPage struct {
	Header PageHeader
	Slots  []Slot
	Data   []byte
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

func (p *SlottedPage) Insert(record []byte) (uint16, error) {
	recordSize := uint16(len(record))
	slotSize := uint16(binary.Size(Slot{}))

	available := p.Header.FreeEnd - p.Header.FreeStart
	if available < recordSize+slotSize {
		return 0, ErrNotSize
	}

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
