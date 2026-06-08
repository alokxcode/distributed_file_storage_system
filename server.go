package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"log/slog"
	"os"
	"sync"

	"github.com/alokxcode/distributed_file_storage_system/p2p"
)

type NodeServerOpts struct {
	FileStore *FileStore
	Transport         p2p.Transport
	BootStrapNodes    []string
}

type NodeServer struct {
	NodeServerOpts
	quitch    chan struct{}
	peersLock sync.Mutex
	Peers     map[string]p2p.Peer
}

func NewNodeServer(nodeServerOpts NodeServerOpts) *NodeServer {

	return &NodeServer{
		NodeServerOpts: nodeServerOpts,
		quitch:         make(chan struct{}),
		Peers:          make(map[string]p2p.Peer),
	}
}

// AddPeer registers a peer in the node's peers map,
// maintaining a record of all peers currently connected to the network
func (server *NodeServer) AddPeer(peer p2p.Peer) error {
	server.peersLock.Lock()
	defer server.peersLock.Unlock()
	server.Peers[peer.RemoteAddr().String()] = peer
	return nil
}

func (server *NodeServer) Start() error {
	err := server.Transport.ListenAndAccept()
	if err != nil {
		slog.Error("failed to listen on port ","error",err)
		return err
	}

	if err := server.BootStrapNetwork(); err != nil {
		slog.Error("failed to bootstrap nodes","error",err)
	}
	server.readloop()

	return nil
}

func (s *NodeServer) readloop() {
	defer func() {
		slog.Info("node server stooped due to user quit action")
		s.Transport.Close()
	}()
	for {
		select {
		case msg := <-s.Transport.Consume():
			if msg.ReqType == p2p.POSTReq {
				slog.Info("writing to disk","size",len(msg.Payload), "payload", string(msg.Payload))
				err := s.FileStore.WriteStream(msg.MetaData.Name,bytes.NewReader(msg.Payload))
				if err != nil {
					slog.Error("Failed to write payload on disk","payload",string(msg.Payload),"error",err)
					return
				}
			}
			if msg.ReqType == p2p.GETReq {

				f,err := s.FileStore.Has(msg.MetaData.Name)
				if err != nil {
					return 
				}

				fileInfo, _ := f.(*os.File).Stat()

				metadata := p2p.File_MetaData{
					Name: msg.MetaData.Name,
					Size: fileInfo.Size(),
				}

				msg.Peer.Write([]byte{p2p.POSTReq})

				msg.Peer.Write([]byte{p2p.MessageTypeGOB})
				gob.NewEncoder(msg.Peer).Encode(metadata)

				msg.Peer.Write([]byte{p2p.MessageTypeRaw})
				io.Copy(msg.Peer,f)
			}
		case <-s.quitch:
			return
		}
	}
}

func (s *NodeServer) Store(metaData p2p.File_MetaData, r io.Reader) error {
	var buff bytes.Buffer
	tee  := io.TeeReader(r,&buff)
	// writes in its own disk
	err := s.FileStore.WriteStream(metaData.Name, tee)
	if err != nil {
		return err
	}

	// dials all BootStrapNetwork to save the file on thier disk
	err = s.BootStrapNetwork()
	if err != nil {
		return err
	}

	// broadcasts to all other peers
	s.Broadcast(metaData, &buff)

	return nil
}

//
// const (
// 	MessageTypeRaw byte = 0x1
// 	MessageTypeGOB byte = 0x2
// )
//

func (s *NodeServer) Broadcast(metaData p2p.File_MetaData, r io.Reader) error {
	writers := make([]io.Writer,0,len(s.Peers))


	for _, p := range s.Peers {
		p.Write([]byte{p2p.POSTReq})
		p.Write([]byte{p2p.MessageTypeGOB})
		enc := gob.NewEncoder(p)
		enc.Encode(metaData)
	}

	for _, p := range s.Peers {
		p.Write([]byte{p2p.MessageTypeRaw})
		writers = append(writers, p)
	}

	mw := io.MultiWriter(writers...)
	_,err := io.Copy(mw,r)
	if err != nil {
		return fmt.Errorf("Failed to stream to Peers : %w",err)
	}
	
	return nil
}

func (s *NodeServer) BootStrapNetwork() error {
	for _, addr := range s.BootStrapNodes {
		if len(addr) == 0 {
			continue
		}
		go func(addr string) {
			err := s.Transport.Dial(addr)
			if err != nil {
				slog.Error("failed to dial","port",addr, "error",err)
			}
		}(addr)
	}
	return nil
}

func (s *NodeServer) Stop() {
	close(s.quitch)
}


func (s *NodeServer) Get(key string) (io.Reader, error) {
	// looks locally
	f,err := s.FileStore.Has(key)
	if err != nil {
		if os.IsNotExist(err) {
			slog.Error("File does not exist", "error", err)
		}

		s.BootStrapNetwork()
		s.Fetch(key)

	} 

	return f,nil
}

func (s *NodeServer) Fetch(key string) {

	metadata := p2p.File_MetaData {
		Name: key,
	}

	for _,p := range s.Peers {
		p.Write([]byte{p2p.GETReq})
		err := gob.NewEncoder(p).Encode(metadata)
		if err != nil {
			return 
		}
	}
}
