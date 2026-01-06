package storage

import (
	"encoding/binary"
	"errors"
	"hash/crc32"
)

// These datatypes and sizes are copied from PostgreSQL
const (
	PageSize = 8192
	// PageSize = 4096
	PageHeaderSize = 32

	PageMagic uint32 = 0xDBDBDBDB
)

var (
	ErrInvalidPageSize  = errors.New("invalid page size")
	ErrInvalidMagic     = errors.New("invalid page magic")
	ErrChecksumMismatch = errors.New("checksum mismatch")
)

type PageHeader struct {
	Magic     uint32 // 04 bytes	//	Magic Number at begginning of page header
	Version   uint16 // 02 bytes //	Version for forward and backwards compatiablity
	Flags     uint16 // 02 bytes // 	Page Tyep flag
	PageID    uint64 // 08 bytes // 	Page Id
	FreeStart uint16 // 02 bytes // Start of free space
	FreeEnd   uint16 // 02 bytes // End of free space
	SlotCount uint16 // 02 bytes
	Checksum  uint32 // 04 bytes //	CRC error detection
	Reserved  uint32 // 04 bytes //	Special Space
}

const (
	OffsetMagic     = 0
	OffsetVersion   = 4
	OffsetFlags     = 6
	OffsetPageID    = 8
	OffsetFreeStart = 16
	OffsetFreeEnd   = 18
	OffsetSlotCount = 20
	OffsetChecksum  = 22
	OffsetReserved  = 26
	OffsetData      = 32
)

type Page struct {
	Header PageHeader // 32 bytes
	Data   []byte     // 8160 bytes
}

func NewPage(pageId uint64) *Page {
	newPage := &Page{
		Header: PageHeader{
			Magic:   PageMagic,
			Version: 1,
			PageID:  pageId,
		},
		Data: make([]byte, PageSize-PageHeaderSize),
	}
	return newPage
}

func serializeHelper(page *Page, checksum uint32) ([]byte, error) {
	buf := make([]byte, PageSize)

	binary.LittleEndian.PutUint32(buf[OffsetMagic:], page.Header.Magic)
	binary.LittleEndian.PutUint16(buf[OffsetVersion:], page.Header.Version)
	binary.LittleEndian.PutUint16(buf[OffsetFlags:], page.Header.Flags)
	binary.LittleEndian.PutUint64(buf[OffsetPageID:], page.Header.PageID)
	binary.LittleEndian.PutUint16(buf[OffsetFreeStart:], page.Header.FreeStart)
	binary.LittleEndian.PutUint16(buf[OffsetFreeEnd:], page.Header.FreeEnd)
	binary.LittleEndian.PutUint16(buf[OffsetSlotCount:], page.Header.SlotCount)
	binary.LittleEndian.PutUint32(buf[OffsetChecksum:], checksum)
	binary.LittleEndian.PutUint32(buf[OffsetReserved:], page.Header.Reserved)

	copy(buf[OffsetData:], page.Data)

	return buf, nil
}

func (page *Page) Serialize() ([]byte, error) {
	if len(page.Data) != PageSize-PageHeaderSize {
		return nil, ErrInvalidPageSize
	}

	buf, _ := serializeHelper(page, 0) //Use zero as a temp checksum
	buf, _ = serializeHelper(page, crc32.ChecksumIEEE(buf))
	return buf, nil
}

func Deserialize(raw []byte) (*Page, error) {
	if len(raw) != PageSize {
		return nil, ErrInvalidPageSize
	}

	magic := binary.LittleEndian.Uint32(raw[OffsetMagic:])
	if magic != PageMagic {
		return nil, ErrInvalidMagic
	}

	checksum := binary.LittleEndian.Uint32(raw[OffsetChecksum:])
	binary.LittleEndian.PutUint32(raw[OffsetChecksum:], 0)
	calculatedChecksum := crc32.ChecksumIEEE(raw)

	if checksum != calculatedChecksum {
		return nil, ErrChecksumMismatch
	}

	header := PageHeader{
		Magic:     magic,
		Version:   binary.LittleEndian.Uint16(raw[OffsetVersion:]),
		Flags:     binary.LittleEndian.Uint16(raw[OffsetFlags:]),
		PageID:    binary.LittleEndian.Uint64(raw[OffsetPageID:]),
		FreeStart: binary.LittleEndian.Uint16(raw[OffsetFreeStart:]),
		FreeEnd:   binary.LittleEndian.Uint16(raw[OffsetFreeEnd:]),
		SlotCount: binary.LittleEndian.Uint16(raw[OffsetSlotCount:]),
		Checksum:  checksum,
		Reserved:  binary.LittleEndian.Uint32(raw[OffsetReserved:]),
	}

	data := make([]byte, PageSize-PageHeaderSize)
	copy(data, raw[OffsetData:])
	page := &Page{
		Header: header,
		Data:   data,
	}

	return page, nil

}
