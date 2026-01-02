package storage

import (
	"bytes"
	"testing"
)

func TestPageSerializeDeserialize(t *testing.T){
	page := NewPage(67)
	copy(page.Data, []byte("FPRS"))


	raw, err := page.Serialize()
	if err != nil{
		t.Fatalf("Serialization Failed: %v", err)
	}

	loader, err := Deserialize(raw)
	if err != nil{
		t.Fatalf("Deserialzation Failed: %v", err)
	}


	if loader.Header.PageID != 67{
		t.Fatalf("Expected PageID 67, got %d", loader.Header.PageID)
	}

	if !bytes.Equal(page.Data, loader.Data){
		t.Fatalf("Data mismatch")
	}

}

func TestChecksumFailure(t *testing.T){
	page := NewPage(1)
	raw, _ := page.Serialize()

	raw[128] ^= 0xFF //Insert dirty data to corrupt 
	// raw[128] ^= 0xEA //Insert dirty data to corrupt 

	_, err := Deserialize(raw)

	if err != ErrChecksumMismatch{
		t.Fatalf("Expected checksum mismatch, got %v", err)
	}

}

func TestInvalidMagic(t *testing.T){
	page := NewPage(42)
	raw, _ := page.Serialize()
	raw[0] = 0x00 //Corrupt Magic Number

	_, err := Deserialize(raw)
	if err != ErrInvalidMagic{
		t.Fatalf("Expected invalid magic numver error")
	}
}
