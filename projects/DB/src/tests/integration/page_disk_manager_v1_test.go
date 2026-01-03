package integration

import (
	"dbp/src/internal/storage"
	"dbp/src/internal/storage/disk"
	"encoding/binary"
	"path/filepath"
	"strings"
	"testing"
)

const (
	UserIdSize   = 8  // int
	UsernameSize = 25 // varchar[25]
	PasswordSize = 25 // varchar[25]
	RowSize      = UserIdSize + UsernameSize + PasswordSize
)

type Row struct {
	ID       uint64
	Username string
	Password string
}

const (
	UserIdOffset   = 0
	UsernameOffset = 8
	PasswordOffset = 33
)

func encodeRow(r Row) []byte {
	buf := make([]byte, RowSize)

	binary.LittleEndian.PutUint64(buf[UserIdOffset:], r.ID)
	copy(buf[UsernameOffset:PasswordOffset], []byte(r.Username))
	copy(buf[PasswordOffset:], []byte(r.Password))

	return buf

}

func decodeRow(buf []byte) Row {
	row := Row{

		ID:       binary.LittleEndian.Uint64(buf[UserIdOffset:]),
		Username: strings.TrimRight(string(buf[UsernameOffset:PasswordOffset]), "\x00"),
		Password: strings.TrimRight(string(buf[PasswordOffset:]), "\x00"),
	}

	return row

}

func makeTestPage(rows []Row) *storage.Page {
	data := make([]byte, storage.PageSize-storage.PageHeaderSize)

	offset := 0

	for _, r := range rows {
		rowBytes := encodeRow(r)
		copy(data[offset:], rowBytes)
		offset += RowSize
	}

	// page := storage.NewPage(1)
	header := storage.PageHeader{
		Magic:   0xDBDBDBDB,
		Version: 1,
		Flags:   1,
		PageID:  0,
	}

	page := &storage.Page{
		Header: header,
		Data:   data,
	}

	return page

}

func TestSinglePageToFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "dbfile")

	dm, err := disk.OpenDiskManager(path, storage.PageSize)
	if err != nil {
		t.Fatalf("TestingSinglePageToFile error: %v", err)
	}

	defer dm.Close()

	rows := []Row{
		{ID: 1, Username: "Leon Trotsky", Password: "Revolutionary"},
		{ID: 2, Username: "Anwar Sadat", Password: "Pragmatist"},
		{ID: 3, Username: "Jean-Paul Sartre", Password: "Existentialism"},
		{ID: 4, Username: "Miles Davis", Password: "Indomitable"},
	}

	page := makeTestPage(rows)
	pid, err := dm.AllocatePage()
	if err != nil {
		t.Fatalf("TestingSinglePageToFile error: %v", err)
	}

	page.Header.PageID = uint64(pid)

	buf, err := page.Serialize()
	if err != nil {
		t.Fatalf("TestingSinglePageToFile error: %v", err)
	}

	err = dm.WritePage(pid, buf)
	if err != nil {
		t.Fatalf("TestingSinglePageToFile error: %v", err)
	}

	err = dm.Sync()
	if err != nil {
		t.Fatalf("TestingSinglePageToFile error: %v", err)
	}

	readBuf := make([]byte, storage.PageSize)
	err = dm.ReadPage(pid, readBuf)
	if err != nil {
		t.Fatalf("TestingSinglePageToFile error: %v", err)
	}

	readPage, err := storage.Deserialize(readBuf)
	if err != nil {
		t.Fatalf("TestingSinglePageToFile error: %v", err)
	}

	for i, r := range rows {
		start := i * RowSize
		got := decodeRow(readPage.Data[start : start+RowSize])

		if r.ID != got.ID || r.Username != got.Username || r.Password != got.Password {

			// if r != got{
			t.Fatalf("TestingSinglePageToFile error: %v", err)
		}

	}

}
