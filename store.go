package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
)

const DefaultRootDir = "alokxcode_Network"

type PathKey struct {
	FirstDir string
	DirPath  string
	FileName string
	// DirPath/FileName
	FullPath string
}

// Path Transformer interface func
type pathTransformFunc func(string) PathKey

// Default Path Transformer
var DefaultPathTransformFunc = func(key string) PathKey {
	return PathKey{
		DirPath:  key,
		FileName: key,
		FullPath: fmt.Sprintf("%s/%s", key, key),
	}
}

// content addressable storage path transformer
func CASPathTransformFunc(key string) PathKey {
	hash := sha1.Sum([]byte(key))
	hashStr := hex.EncodeToString(hash[:])

	blockSize := 8
	sliceLen := len(hashStr) / blockSize
	paths := make([]string, sliceLen)

	for i := range sliceLen {
		from := i * blockSize
		to := (i * blockSize) + blockSize
		paths[i] = hashStr[from:to]
	}
	pathName := strings.Join(paths, "/")

	return PathKey{
		FirstDir: paths[0],
		DirPath:  pathName,
		FileName: hashStr,
		FullPath: fmt.Sprintf("%s/%s", pathName, hashStr),
	}
}

// store config
type FileStoreOpts struct {
	// StorageRoot is the folder name of root containing all the folder/files of the system.
	StorageRoot       string
	pathTransformFunc pathTransformFunc
}

type FileStore struct {
	FileStoreOpts
}

func NewFileStore(opts FileStoreOpts) *FileStore {
	if len(opts.StorageRoot) == 0 {
		opts.StorageRoot = DefaultRootDir
	}
	return &FileStore{
		FileStoreOpts: opts,
	}
}

func (s *FileStore) WriteStream(key string, r io.Reader) error {
	pathKey := s.pathTransformFunc(key)
	nestedHashDirPath := s.StorageRoot + "/" + pathKey.DirPath

	// makes the nested dirs
	err := os.MkdirAll(nestedHashDirPath, os.ModePerm)
	if err != nil {
		return err
	}

	fullPathWithRoot := s.StorageRoot + "/" + pathKey.FullPath

	// creates the file
	f, err := os.Create(fullPathWithRoot)
	if err != nil {
		return err
	}

	// writes bytes to file
	n, err := io.Copy(f, r)
	if err != nil {
		return err
	}

	slog.Info("written to disk","size", n,"path", pathKey.FullPath)

	return nil
}

func (s *FileStore) ReadStream(key string) (io.ReadCloser, error) {
	pathKey := s.pathTransformFunc(key)
	return os.Open(s.StorageRoot + "/" + pathKey.FullPath)
}

func (s *FileStore) Read(key string) (io.Reader, error) {
	f, err := s.ReadStream(key)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buff := new(bytes.Buffer)
	_, err = io.Copy(buff, f)
	return buff, err
}

func (s *FileStore) Delete(key string) error {
	pathKey := s.pathTransformFunc(key)
	defer func() {
		fmt.Printf("delete [%s] from the disk", pathKey.FileName)
	}()

	err := os.RemoveAll(s.StorageRoot + "/" + pathKey.FirstDir)
	if err != nil {
		return err
	}
	return nil
}

func (s *FileStore) DeleteAll() error {
	defer func() {
		fmt.Printf("root dir deleted : %s", s.StorageRoot)
	}()

	err := os.RemoveAll(s.StorageRoot)
	if err != nil {
		return err
	}
	return nil
}

func (s *FileStore) Has(key string) (io.Reader,error) {
	pathKey := s.pathTransformFunc(key)

	_,err := os.Stat(pathKey.FullPath)
	if err != nil {
		return nil,err
	}

	f,err := os.Open(pathKey.FullPath)
	if err != nil {
		return nil,err
	}

	return f,nil
 
}
