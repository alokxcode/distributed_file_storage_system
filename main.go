package main

import (
	"fmt"
	"log"

	"github.com/alokxcode/distributed_file_storage_system/p2p"
)

func main() {
	opts := p2p.TCPTransportOpts{
		ListenAddress: ":3000",
		ShakeHands:    p2p.NOPHandShakeFunc,
		Decoder:       p2p.DefaultDecoder{},
		OnPeer: func(p2p.Peer) error {return fmt.Errorf("failed the OnPeer")},
	}

	tr := p2p.NewTCPTransport(opts)
	
	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}
	for {
		fmt.Printf("%v\n", <- tr.Consume())
	}
}
