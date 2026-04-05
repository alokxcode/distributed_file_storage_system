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

const DefaultRootDir = "DJAlok_Network"

type PathKey struct {
	FirstPath string
	PathName string
	FileName string
}

// returns full path
func (p *PathKey) FullPath() string {
	return fmt.Sprintf("%s/%s",p.PathName,p.FileName)
}

// Path Transformer interface func
type pathTransformFunc func(string) PathKey 

// Default Path Transformer
var DefaultPathTransformFunc = func(key string) PathKey {
	return PathKey{
		PathName: key,
		FileName: key,
	} 
}

// content addressable storage path transformer 
func CASPathTransformFunc(key string) PathKey {
	hash := sha1.Sum([]byte(key))  
	hashStr := hex.EncodeToString(hash[:]) 

	blockSize := 8
	sliceLen := len(hashStr) / blockSize
	paths := make([]string,sliceLen)

	for i:= range sliceLen {
		from := i*blockSize
		to := (i*blockSize)+blockSize
		paths[i] = hashStr[from:to]
	}
	pathName := strings.Join(paths,"/")

	return PathKey{
		FirstPath: paths[0],
		PathName: pathName,
		FileName: hashStr,
	}
}

// store config
type StoreOpts struct {
	// root is the folder name of root containing all the folder/files of the system.
	root string
	pathTransformFunc pathTransformFunc

}

type Store struct {
	StoreOpts
}

func NewStore(opts StoreOpts) *Store {
	if len(opts.root) == 0 {
		opts.root = DefaultRootDir
	} 
	return &Store {
		StoreOpts: opts,
	}
}

func (s *Store) WriteStream(key string, r io.Reader) error {
	pathKey := s.pathTransformFunc(key)

	// makes the nested dirs
	err := os.MkdirAll(s.root + "/" + pathKey.PathName, os.ModePerm)
	if err != nil {
		return err
	}

	fullPath := pathKey.FullPath()
	f,err := os.Create(s.root + "/" + fullPath)
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

func (s *Store) ReadStream(key string) (io.ReadCloser,error) {
	pathKey := s.pathTransformFunc(key)
	return os.Open(s.root + "/" + pathKey.FullPath())
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

func (s *Store) Delete(key string) error {
	pathKey := s.pathTransformFunc(key)
	defer func() {
		fmt.Printf("delete [%s] from the disk",pathKey.FileName)
	}()

	err := os.RemoveAll(s.root + "/" + pathKey.FirstPath)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) DeleteAll() error {
	defer func() {
		fmt.Printf("root dir deleted : %s",s.root)
	}()

	err := os.RemoveAll(s.root)
	if err != nil {
		return err
	}
	return nil
}


