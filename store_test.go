package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"
)

func TestStore(t *testing.T) {
	opts := StoreOpts {
		pathTransformFunc: CASPathTransformFunc,
	}

	s := NewStore(opts)
	data := bytes.NewReader([]byte("some jpg"))
	key := "myspecailpic"
	err := s.WriteStream(key,data)
	if err != nil {
		t.Error(err)
	}

	r,err := s.Read(key)
	if err != nil {
		t.Error(err)
	}

	b,err := ioutil.ReadAll(r)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(b))

	


	
}
