package storage

import "testing"

func TestSlottedPageInsertRead(t *testing.T) {
	p := NewSlottedPage(1)
	id, err := p.Insert([]byte("hello"))
	if err != nil {
		t.Fatal(err)
	}

	data, err := p.Get(id)
	if err != nil {
		t.Fatal(err)
	}

	if string(data) != "hello" {
		t.Fatalf("expected hello, got %s", data)
	}
}

func TestSlottedPageDelete(t *testing.T) {
	p := NewSlottedPage(1)
	id, _ := p.Insert([]byte("FRPS"))

	err := p.Delete(id)
	if err != nil {
		t.Fatal(err)
	}

	_, err = p.Get(id)
	if err == nil {
		t.Fatal("expected error for deleted record")
	}

}
