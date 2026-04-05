package main

type FileServerOpts struct {
	Listend_addr string
	StorageRoot string
} 

type FileServer struct {
	FileServerOpts
	Store *Store
}

func NewFileServer(opts FileServerOpts, store Store) *FileServer {
	return &FileServer{
		FileServerOpts: opts,
		Store: &store,
	}
}
