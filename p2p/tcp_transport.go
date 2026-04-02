package p2p

import (
	"fmt"
	"net"
)

// TCPPeer represents the remote node over a TCP established connection
type TCPPeer struct {

	// conn is the underlying connnection of  the peer
	conn net.Conn

	// if we dial and retrive a connection => outbound = true
	// if we accept and retrive a connection => outbound = false
	outbound bool
}

// close implements the Peer interface
func (p *TCPPeer) Close() error {
	return p.Close()
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
	OnPeer func(Peer) error
}

type TCPTransport struct {
	TCPTransportOpts
	listener net.Listener
	rpcch chan RPC
}

func NopeHandShake() error {
	
	return nil
}
func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
		rpcch: make(chan RPC),
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

// consume implements the Transport interface, which will return read-only channnel
// for reading the incoming messages recieved from another peer in the network
func (t *TCPTransport) Consume() <- chan RPC {
	return t.rpcch
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
	var err error
	defer func() {
		conn.Close()
		fmt.Println("dropping peer connection", err)
	}()

	peer := NewTCPPeer(conn, true)

	// do handshake
	if err = t.ShakeHands(peer); err != nil {
		conn.Close()
		fmt.Printf("TCP handshake error : %+v", err)
		return
	}

	if t.OnPeer != nil {
		err = t.OnPeer(peer)
		if err != nil {
			return
		}
	}

	// read loop
	rpc := RPC{}
	for {
		if err = t.Decoder.Decode(conn, &rpc); err != nil {
			fmt.Printf("TCP error : %+v\n", err)
			continue
		}
		rpc.From = conn.RemoteAddr()
		t.rpcch <- rpc
	}

}
