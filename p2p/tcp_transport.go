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

type TCPTransport struct {
	listenAddress string
	listener      net.Listener
	shakeHands    HandShakeFunc
	decoder       Decoder

	mu    sync.RWMutex
	peers map[net.Addr]Peer
}

func NopeHandShake() error {
	return nil
}
func NewTCPTransport(listenAddr string) *TCPTransport {
	return &TCPTransport{
		listenAddress: listenAddr,
		shakeHands:    NOPHandShakeFunc,
	}

}

func (t *TCPTransport) ListenAndAccept() error {
	var err error
	t.listener, err = net.Listen("tcp", t.listenAddress)
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

		fmt.Printf("incoming new connection %+v", conn)
		go t.handleConnnection(conn)
	}
}

func (t *TCPTransport) handleConnnection(conn net.Conn) {
	peer := NewTCPPeer(conn, true)

	// do handshake
	if err := t.shakeHands(conn); err != nil {
		conn.Close()
		fmt.Printf("TCP handshake error : %+v", err)
		return
	}

	// read loop
	for {
		if err := t.decoder.Decode(peer); err != nil {
			fmt.Printf("TCP error : %+v", err)
			continue
		}
	}

}
