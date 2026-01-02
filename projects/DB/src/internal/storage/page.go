package storage

import (
	"bytes"
	"encoding/binary"
	"errors"	
	"hash/crc32"
)

const(
	PageSize = 8192
	// PageSize = 4096
	PageHeaderSize = 32

	PageMagic uint32 = 0xDBDBDBDB 
)

var(
	ErrInvalidPageSize = errors.New("invalid page size")
	ErrInvalidMagic = errors.New("invalid page magic")
	ErrChecksumMismatch = errors.New("checksum mismatch")
)

type PageHeader struct{
	Magic uint32	//	Magic Number at begginning of page header
	Version uint16	//	Version for forward and backwards compatiablity
	Flags uint16	// 	Page Tyep flag
	PageID uint64	// 	Page Id
	DataSize uint32 //	Bytes Used
	Checksum uint32 //	CRC error detection
	Reserved uint64 //	Special Space
}

type Page struct{
	Header PageHeader
	Data	[]byte
}

func NewPage(pageId uint64) *Page{
	newPage := &Page{
		Header: PageHeader{
			Magic: PageMagic,
			Version: 1,
			PageID : pageId,
		},
		Data: make([]byte, PageSize-PageHeaderSize),
	}
	return newPage
}

func (page *Page) Serialize()([]byte, error){
	if len(page.Data) != PageSize-PageHeaderSize{
		return nil, ErrInvalidPageSize
	}

	buf := make([]byte, PageSize)

	page.Header.Checksum = 0 

	writer := bytes.NewBuffer(buf[:0])
	if err := binary.Write(writer, binary.LittleEndian, &page.Header); err != nil{
		return nil, err
	}

	copy(buf[PageHeaderSize:], page.Data)

	checksum := crc32.ChecksumIEEE(buf)
	page.Header.Checksum = checksum

	writer.Reset()
	if err:= binary.Write(writer, binary.LittleEndian, &page.Header); err!= nil{
		return nil, err
	}
	return buf, nil

}

func Deserialize(raw []byte) (*Page, error){
	if len(raw) != PageSize{
		return nil, ErrInvalidPageSize
	}

	var header PageHeader
	reader := bytes.NewReader(raw[:PageHeaderSize])

	if err := binary.Read(reader, binary.LittleEndian, &header); err != nil{
		return nil, err
	}

	if header.Magic != PageMagic{
		return nil, ErrInvalidMagic
	}

	expected := header.Checksum
	// header.Checksum = 0
	binary.LittleEndian.PutUint32(
		raw[20:24],
		0,
		)

	calculated := crc32.ChecksumIEEE(raw)

	if expected != calculated{
		return nil, ErrChecksumMismatch
	}

	binary.LittleEndian.PutUint32(
		raw[20:24],
		expected,
		)

	page := &Page{
		Header: header,
		Data: make([]byte, PageSize - PageHeaderSize),
	}

	copy(page.Data, raw[PageHeaderSize:])

	return page, nil

}
