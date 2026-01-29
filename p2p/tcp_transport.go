package p2p

import (
	"fmt"
	"net"
	"sync"
)

// TCPPeer represents the remote node over a TCP established connection
type TCPPeer struct {

	// conn is the underlying connnection of  the peer
	conn net.Conn

	// if we dial and retrive a connection => outbound = true
	// if we accept and retrive a connection => outbound = false
	outbound bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

type TCPTransportOpts struct {
	ListenAddress string
	ShakeHands    HandShakeFunc
	Decoder       Decoder
}

type TCPTransport struct {
	TCPTransportOpts
	listener net.Listener

	mu    sync.RWMutex
	peers map[net.Addr]Peer
}

func NopeHandShake() error {
	return nil
}
func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
	}

}

func (t *TCPTransport) ListenAndAccept() error {
	var err error
	t.listener, err = net.Listen("tcp", t.ListenAddress)
	if err != nil {
		return err
	}

	go t.startAcceptLoop()

	return nil
}

func (t *TCPTransport) startAcceptLoop() error {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			return err
		}

		fmt.Printf("incoming new connection %+v\n", conn)
		go t.handleConnnection(conn)
	}
}

func (t *TCPTransport) handleConnnection(conn net.Conn) {
	peer := NewTCPPeer(conn, true)

	// do handshake
	if err := t.ShakeHands(peer); err != nil {
		conn.Close()
		fmt.Printf("TCP handshake error : %+v", err)
		return
	}

	// read loop
	rpc := &RPC{}
	// buff := make([]byte, 2000)
	for {
		// n, err := conn.Read(buff)
		// if err != nil {
		// 	return
		// }
		if err := t.Decoder.Decode(conn, rpc); err != nil {
			fmt.Printf("TCP error : %+v\n", err)
			continue
		}
		rpc.From = conn.RemoteAddr()
		fmt.Printf("message : %v\n", rpc)
	}

}
