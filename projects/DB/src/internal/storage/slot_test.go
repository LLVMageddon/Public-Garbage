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

func TestSlottedPageReclaim(t *testing.T) {

	p := NewSlottedPage(1)
	id, err := p.Insert([]byte("hello1"))
	if err != nil {
		t.Fatal(err)
	}
	data, err := p.Get(0)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "hello1" {
		t.Fatalf("expected hello1, got %s", data)
	}

	id, err = p.Insert([]byte("hello2"))
	if err != nil {
		t.Fatal(err)
	}
	data, err = p.Get(1)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "hello2" {
		t.Fatalf("expected hello2, got %s", data)
	}

	id, err = p.Insert([]byte("hello3"))
	if err != nil {
		t.Fatal(err)
	}
	data, err = p.Get(2)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "hello3" {
		t.Fatalf("expected hello3, got %s", data)
	}

	err = p.Delete(id - 1)
	if err != nil {
		t.Fatal(err)
	}

	id, err = p.Insert([]byte("hllo"))
	if err != nil {
		t.Fatal(err)
	}
	data, err = p.Get(1)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "hllo" {
		t.Fatalf("expected hllo, got %s", data)
	}

	err = p.Delete(id)
	if err != nil {
		t.Fatal(err)
	}

	id, err = p.Insert([]byte("Chewable Code"))
	if err != nil {
		t.Fatal(err)
	}
	data, err = p.Get(3)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "Chewable Code" {
		t.Fatalf("expected hllo, got %s", data)
	}

}
