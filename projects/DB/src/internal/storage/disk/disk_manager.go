package disk

import (
	"errors"
	"fmt"
	"os"
)

var(
	ErrInvalidPageID = errors.New("invalid page id")
	ErrCorruptPage = errors.New("corrupt page detected")
)

type PageID uint64

type DiskManager struct{
	file *os.File
	pageSize uint64
	path string 
}

func OpenDiskManager(path string, pageSize uint64) (*DiskManager, error){
	flag := os.O_RDWR | os.O_CREATE
	file, err := os.OpenFile(path, flag, 0644)
	if err != nil{
		return nil, err
	}

	info, err := file.Stat()
	if err != nil{
		return nil, err
	}

	if info.Size()%int64(pageSize) != 0{
		return nil, fmt.Errorf("invalid db file size")
	}

	return &DiskManager{
		file: file,
		pageSize: pageSize,
		path: path,
	}, nil
}

func (dm *DiskManager) pageOffset(pid PageID) int64{
	return int64(pid) * int64(dm.pageSize)
}

func (dm *DiskManager) ReadPage(pid PageID, buf []byte) error{
	if uint64(len(buf))!= dm.pageSize {
		return fmt.Errorf("buffer size mismatch")
	}

	offset := dm.pageOffset(pid)
	n, err := dm.file.ReadAt(buf, offset)
	if err != nil{
		return err
	}

	if uint64(n) != dm.pageSize{
		return ErrCorruptPage
	}

	return nil

}

func (dm *DiskManager) WritePage(pid PageID, data []byte) error{
	if uint64(len(data))!= dm.pageSize{
		fmt.Errorf("page size mismatch")
	}

	tmpPath := fmt.Sprintf("%s.page.%d.tmp", dm.path, pid)
	tmpFile, err := os.OpenFile(tmpPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil{
		return err
	}
	defer os.Remove(tmpPath)

	if _, err := tmpFile.Write(data); err != nil{
		return err
	}

	if err := tmpFile.Sync(); err != nil{
		return err
	}

	if err := tmpFile.Close(); err != nil{
		return err
	}

	targetOffset := dm.pageOffset(pid)

	if _ , err := dm.file.WriteAt(data, targetOffset); err != nil{
		return err 
	}

	return nil
}

func (dm *DiskManager) AllocatePage() (PageID, error){
	info, err := dm.file.Stat()
	if err != nil{
		return 0, err
	}

	pid := PageID(info.Size()/ int64(dm.pageSize))
	zero := make([]byte, dm. pageSize)

	if _, err := dm.file.Write(zero); err != nil{
		return 0, err
	}

	return pid, nil
}

func (dm *DiskManager) Sync() error{
	return dm.file.Sync()
}

func (dm *DiskManager) Close() error{
	return dm.file.Close()
}
