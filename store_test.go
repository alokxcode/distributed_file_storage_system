package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
"testing"
)

// tests newStore
func TestStore(t *testing.T) {
	fileStoreOpts := FileStoreOpts {
		pathTransformFunc: CASPathTransformFunc,
	}

	s := NewFileStore(fileStoreOpts)
	fmt.Println("fileStore created",s)
}

// tests WriteSteam
func TestWriteStream(t *testing.T) {
	opts := FileStoreOpts {
		pathTransformFunc: CASPathTransformFunc,
	}

	s := NewFileStore(opts)
	data := []byte("first jpg")
	key := "file" 

	err := s.WriteStream(key,bytes.NewReader(data))
	if err != nil {
		t.Error(err)
	}
	
	
}

// tests ReadStream
func TestReadStream(t *testing.T) {
	opts := FileStoreOpts {
		pathTransformFunc: CASPathTransformFunc,
	}

	s := NewFileStore(opts)
	data := []byte("first jpg")
	key := "second_file"

	// Write the file to the store first before reading it
	err := s.WriteStream(key, bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}

	r, err := s.Read(key)
	if err != nil {
		t.Fatal(err)
	}

	b, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(b))

	if string(b) != string(data) {
		t.Errorf("want %s have %s", data, b)
	}
}


// tests delete
func TestDelete(t *testing.T) {
	opts := FileStoreOpts {
		pathTransformFunc: CASPathTransformFunc,
	}

	s := NewFileStore(opts)
	key := "second_file"

	err := s.Delete(key)
	if err != nil {
		t.Error(err)
	}
}

// tests DeleteAll
func TestDeleteAll(t *testing.T) {
	opts := FileStoreOpts {
		pathTransformFunc: CASPathTransformFunc,
	}

	s := NewFileStore(opts)

	err := s.DeleteAll()
	if err != nil {
		t.Error(err)
	}
}





