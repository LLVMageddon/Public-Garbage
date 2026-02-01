package storage

import (
	"testing"
)

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
	// p.compact2()
}

func TestTest(t *testing.T) {

	p := NewSlottedPage(1)
	p.Insert([]byte("1, Alice, Band"))
	p.Insert([]byte("2, Robert, Marley"))
	p.Insert([]byte("3, Beate, Becker"))
	p.Delete(1)
	p.Insert([]byte("4, Daikia, Sato"))
	p.Delete(1)
	p.Insert([]byte("Chewable Code"))
	p.Delete(3)
	p.Insert([]byte("More Chewable Code"))
	const pageSize = 4096
	page := make([]byte, pageSize)

	id, err := p.Insert(page)
	id += 0
	if err != nil {
		t.Fatal(err)
	}
	p.Insert([]byte("4, Daikia, Sato"))

	data, err := p.Get(0)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "1, Alice, Band" {
		t.Fatalf("expected {1, Alice, Band}, got {%s}", data)
	}
	data, err = p.Get(1)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "Chewable Code" {
		t.Fatalf("expected {Chewable Code}, got {%s}", data)
	}
	data, err = p.Get(2)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "3, Beate, Becker" {
		t.Fatalf("expected {3, Beate, Becker}, got {%s}", data)
	}
	data, err = p.Get(3)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "More Chewable Code" {
		t.Fatalf("expected {More Chewable Code}, got {%s}", data)
	}
	err = p.Delete(4)
	if err != nil {
		t.Fatal(err)
	}
	p.compact()
	data, err = p.Get(4)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "4, Daikia, Sato" {
		t.Fatalf("expected {4, Daikia, Sato}, got {%s}", data)
	}

}
