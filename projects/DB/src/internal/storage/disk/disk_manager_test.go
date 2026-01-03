package disk

import (
	"dbp/src/internal/storage"
	"path/filepath"
	"testing"
)

const testPageSize = 8192

func newTestDM(t *testing.T) *DiskManager {
	path := filepath.Join(t.TempDir(), "test.db")
	dm, err := OpenDiskManager(path, testPageSize)

	if err != nil {
		t.Fatalf("Open Disk Manager error: %v", err)
	}

	return dm
}

func TestAllocateAndReadWritePage(t *testing.T) {
	dm := newTestDM(t)
	defer dm.Close()

	pid, err := dm.AllocatePage()
	if err != nil {
		t.Fatalf("Allocate Page error: %v", err)
	}

	data := make([]byte, testPageSize)
	copy(data, []byte("FRPS"))

	page := storage.NewPage(67)
	copy(page.Data, data)

	// serialized_page, err := page.Serialize()

	if err := dm.WritePage(pid, data); err != nil {
		// if err := dm.WritePage(pid, serialized_page); err != nil{
		t.Fatalf("Write Page error: %v", err)
	}

	read := make([]byte, testPageSize)
	if err := dm.ReadPage(pid, read); err != nil {
		t.Fatalf("Read Page error: %v", err)
	}

	// deserialize_page, err := storage.Deserialize(read)

	// if string(deserialize_page.Data[:4]) != "FRPS"{
	// 	t.Fatalf("data mismatch")
	// }

	if string(read[:4]) != "FRPS" {
		t.Fatalf("data mismatch")
	}

}

func TestInvalidPageRead(t *testing.T) {
	dm := newTestDM(t)
	defer dm.Close()

	buf := make([]byte, testPageSize)
	err := dm.ReadPage(9999, buf)
	if err == nil {
		t.Fatalf("Expected error")
	}
}

func TestPersistenceAcrossRestart(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "persist.db")

	//Scope this
	{
		// dm, err  OpenDiskManager(path, testPageSize)
		dm, err := OpenDiskManager(path, testPageSize)
		if err != nil {
			t.Fatalf("Open Disk Manaager error: %v", err)
		}

		pid, _ := dm.AllocatePage()
		data := make([]byte, testPageSize)
		copy(data, []byte("Famous last words, Fuck This!"))
		_ = dm.WritePage(pid, data)
		_ = dm.Sync()
		_ = dm.Close()

	}
	//Scope this
	{
		dm, err := OpenDiskManager(path, testPageSize)
		if err != nil {
			t.Fatalf("Open Disk Manaager error: %v", err)
		}

		buf := make([]byte, testPageSize)
		if err := dm.ReadPage(0, buf); err != nil {
			t.Fatalf("Read Page error: %v", err)
		}

		if string(buf[:29]) != "Famous last words, Fuck This!" {
			t.Fatalf("Data not durable")
		}
	}
}
