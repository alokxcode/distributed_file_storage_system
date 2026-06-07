package main

import (
	"bytes"
	"log"
	"log/slog"
	"time"

	"github.com/alokxcode/distributed_file_storage_system/p2p"
)

func makeNodeServer(ListenAddress string, nodes ...string) *NodeServer {
	tcpTransport := p2p.NewTCPTransport(p2p.TCPTransportOpts{
		ListenAddress: ListenAddress,
		ShakeHands:    p2p.NOPHandShakeFunc,
		Decoder:       p2p.DefaultDecoder{},
	})

	fileStore := NewFileStore(FileStoreOpts{
		StorageRoot:       ListenAddress,
		pathTransformFunc: CASPathTransformFunc,
	})

	nodeServer := NewNodeServer(NodeServerOpts{
		FileStore:      fileStore,
		Transport:      tcpTransport,
		BootStrapNodes: nodes,
	})

	tcpTransport.OnPeerConnect = nodeServer.AddPeer

	return nodeServer
}

func main() {

	s1 := makeNodeServer(":3000")
	go func() {
		slog.Info("starting node server","server","s1")
		if err := s1.Start(); err != nil {
			slog.Error("failed to start the node server","error",err)
		}
	}()


	s2 := makeNodeServer(":4000", ":3000")
	go func() {
		time.Sleep(50*time.Millisecond)
		slog.Info("starting node server ", "server","s2")
		if err := s2.Start(); err != nil {
			slog.Error("failed to start the node server","error",err)
		}
	}()


	time.Sleep(100*time.Millisecond)
	testData := []byte("hello from server2 : PORT - 4000")

	file_metaData := p2p.File_MetaData {
		Name: "test.jpeg",
		Size: int64(len(testData)),
	} 
	if err := s2.Store(file_metaData, bytes.NewReader(testData)); err != nil {
		log.Fatal(err)
	}

	select {}
}
