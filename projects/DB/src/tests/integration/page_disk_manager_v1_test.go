package integration

import (
	"dbp/src/internal/storage"
	"dbp/src/internal/storage/disk"
	"encoding/binary"
	"fmt"
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
		t.Fatalf("TestingSinglePageToFile: Open Disk Manager error: %v", err)
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
		t.Fatalf("TestingSinglePageToFile: Allocating Page error: %v", err)
	}

	page.Header.PageID = uint64(pid)

	buf, err := page.Serialize()
	if err != nil {
		t.Fatalf("TestingSinglePageToFile: Serialize error: %v", err)
	}

	err = dm.WritePage(pid, buf)
	if err != nil {
		t.Fatalf("TestingSinglePageToFile: Write Page error: %v", err)
	}

	err = dm.Sync()
	if err != nil {
		t.Fatalf("TestingSinglePageToFile: Sync error: %v", err)
	}

	readBuf := make([]byte, storage.PageSize)
	err = dm.ReadPage(pid, readBuf)
	if err != nil {
		t.Fatalf("TestingSinglePageToFile: Read Page error: %v", err)
	}

	readPage, err := storage.Deserialize(readBuf)
	if err != nil {
		t.Fatalf("TestingSinglePageToFile: Deserialize error: %v", err)
	}

	for i, r := range rows {
		start := i * RowSize
		got := decodeRow(readPage.Data[start : start+RowSize])
		if r != got {
			t.Fatalf("TestingSinglePageToFile: Row Values error")
		}

	}

}

func TestMultiplePageSingleFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "dbfile")

	dm, err := disk.OpenDiskManager(path, storage.PageSize)
	if err != nil {
		t.Fatalf("TestMultiplePageSingleFile: Open Disk Manager error: %v", err)
	}

	defer dm.Close()

	pageCount := 3
	userCounter := 1

	for p := range pageCount {

		rows := []Row{}

		rows = append(rows, Row{
			ID:       uint64((p+1)*10 + 1),
			Username: fmt.Sprintf("User %d", userCounter),
			Password: "Password1"},
		)
		rows = append(rows, Row{
			ID:       uint64((p+1)*10 + 2),
			Username: fmt.Sprintf("User %d", userCounter+1),
			Password: "Password2"},
		)
		rows = append(rows, Row{
			ID:       uint64((p+1)*10 + 3),
			Username: fmt.Sprintf("User %d", userCounter+2),
			Password: "Password3"},
		)

		userCounter += len(rows)

		page := makeTestPage(rows)
		pid, err := dm.AllocatePage()
		if err != nil {
			t.Fatalf("TestMultiplePageSingleFile: Allocate error: %v", err)
		}

		page.Header.PageID = uint64(pid)

		buf, err := page.Serialize()

		if err != nil {
			t.Fatalf("TestMultiplePageSingleFile: Serialize error: %v", err)
		}

		err = dm.WritePage(pid, buf)
		if err != nil {
			t.Fatalf("TestMultiplePageSingleFile: Write Page error: %v", err)
		}

	}

	err = dm.Sync()
	if err != nil {
		t.Fatalf("TestMultiplePageSingleFile: Sync error: %v", err)
	}

	//Read file
	userCounter = 1
	for pid := uint64(0); pid < uint64(pageCount); pid++ {
		buf := make([]byte, storage.PageSize)

		err = dm.ReadPage(disk.PageID(pid), buf)
		if err != nil {
			t.Fatalf("TestMultiplePageSingleFile: Read Page error: %v", err)
		}

		page, err := storage.Deserialize(buf)
		if err != nil {
			t.Fatalf("TestMultiplePageSingleFile: Deserialize error: %v", err)
		}

		expectedRows := []Row{
			{ID: uint64((pid+1)*10 + 1), Username: fmt.Sprintf("User %d", userCounter), Password: "Password1"},
			{ID: uint64((pid+1)*10 + 2), Username: fmt.Sprintf("User %d", userCounter+1), Password: "Password2"},
			{ID: uint64((pid+1)*10 + 3), Username: fmt.Sprintf("User %d", userCounter+2), Password: "Password3"},
		}
		userCounter += len(expectedRows)

		rowStartOffset := 0
		rowEndOffset := RowSize

		actualRows := []Row{}

		for row := 0; row < len(expectedRows); row++ {
			actualRows = append(actualRows, decodeRow(page.Data[rowStartOffset:rowEndOffset]))
			rowStartOffset = rowEndOffset
			rowEndOffset += RowSize

		}

		for count := 0; count < len(expectedRows); count++ {
			if expectedRows[count] != actualRows[count] {
				t.Fatalf("TestMultiplePageSingleFile: Validation test error:\nExpected: %v\nGot: %v", expectedRows[count], actualRows[count])
			}
		}

	}

}

func TestMultiplePageMultipleFiles(t *testing.T) {
	dir := t.TempDir()

	for f := 0; f < 2; f++ {
		path := filepath.Join(dir, fmt.Sprintf("dbfile_%d", f))

		dm, err := disk.OpenDiskManager(path, storage.PageSize)
		if err != nil {
			t.Fatalf("TestMultiplePageMultipleFiles: Open Disk Manager error: %v", err)
		}

		rows := []Row{
			{ID: uint64(f*100 + 1), Username: "root", Password: "RootP@ssword"},
		}

		page := makeTestPage(rows)
		pid, err := dm.AllocatePage()
		if err != nil {
			t.Fatalf("TestMultiplePageMultipleFiles: Allocate error: %v", err)
		}
		page.Header.PageID = uint64(pid)

		buf, err := page.Serialize()
		if err != nil {
			t.Fatalf("TestMultiplePageMultipleFiles: Serialize error: %v", err)
		}

		err = dm.WritePage(pid, buf)
		if err != nil {

			t.Fatalf("TestMultiplePageMultipleFiles: Write Page error: %v", err)
		}

		err = dm.Sync()
		if err != nil {
			t.Fatalf("TestMultiplePageMultipleFiles: Sync error: %v", err)
		}

		err = dm.Close()
		if err != nil {
			t.Fatalf("TestMultiplePageMultipleFiles: Close error: %v", err)
		}

		//Reopen file

		dm, err = disk.OpenDiskManager(path, storage.PageSize)
		if err != nil {
			t.Fatalf("TestMultiplePageMultipleFiles: Open Disk Manager error: %v", err)
		}

		readBuf := make([]byte, storage.PageSize)
		err = dm.ReadPage(pid, readBuf)
		if err != nil {
			t.Fatalf("TestMultiplePageMultipleFiles: Read Page error: %v", err)
		}

		readPage, err := storage.Deserialize(readBuf)
		if err != nil {
			t.Fatalf("TestMultiplePageMultipleFiles: Deserialize error: %v", err)
		}

		row := decodeRow(readPage.Data[:RowSize])

		if rows[0] != row {
			t.Fatalf("TestMultiplePageMultipleFiles: Decode Row error:\nExpected: %v\nGot: %v", rows[0], row)
		}

		err = dm.Close()
		if err != nil {
			t.Fatalf("TestMultiplePageMultipleFiles: Close error: %v", err)
		}
	}

}
