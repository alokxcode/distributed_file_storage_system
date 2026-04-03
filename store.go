package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func CASPathTransformFunc(key string) PathKey {
	hash := sha1.Sum([]byte(key))
	hashStr := hex.EncodeToString(hash[:]) 

	blockSize := 8
	sliceLen := len(hashStr) / blockSize
	paths := make([]string,sliceLen)

	for i:=0; i<sliceLen; i++ {
		from := i*blockSize
		to := (i*blockSize)+blockSize
		paths[i] = hashStr[from:to]
	}
	pathName := strings.Join(paths,"/")

	return PathKey{
		PathName: pathName,
		FileName: hashStr,
	}
}

type pathTransformFunc func(string) PathKey 
var DefaultPathTransformFunc = func(key string) string {
	return key
}

type PathKey struct {
	PathName string
	FileName string
}

func (p *PathKey) FullPath() string {
	return fmt.Sprintf("%s/%s",p.PathName,p.FileName)
}

type StoreOpts struct {
	pathTransformFunc pathTransformFunc

}

type Store struct {
	StoreOpts

}

func NewStore(opts StoreOpts) *Store {
	return &Store {
		StoreOpts: opts,
	}
}

func (s *Store) Read(key string) (io.Reader,error){
	f,err := s.ReadStream(key)
	if err != nil {
		return nil,err
	}
	defer f.Close()

	buff := new(bytes.Buffer)
	_,err = io.Copy(buff,f)
	return buff,err
}

func (s *Store) ReadStream(key string) (io.ReadCloser,error) {
	pathKey := s.pathTransformFunc(key)
	return os.Open(pathKey.FullPath())
}

func (s *Store) WriteStream(key string, r io.Reader) error {
	pathKey := s.pathTransformFunc(key)
	err := os.MkdirAll(pathKey.PathName, os.ModePerm)
	if err != nil {
		return err
	}

	fullPath := pathKey.FullPath()
	f,err := os.Create(fullPath)
	if err != nil {
		return nil
	}
	
	n, err := io.Copy(f,r)
	if err != nil {
		return err
	}

	log.Printf("written %v bytes to disk : %s",n,fullPath)


	return nil
}
