package p2p

import (
	"errors"
	"fmt"
	"log/slog"
	"net"
)

// TCPPeer represents the remote node over a TCP established connection
type TCPPeer struct {
	// conn is the underlying connnection of  the peer
	net.Conn

	// if we dial and retrive a connection => outbound = true
	// if we accept and retrive a connection => outbound = false
	outbound bool
}

// creates a new tcp peer
func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		Conn:     conn,
		outbound: outbound,
	}
}

// close implements the Peer interface
func (p *TCPPeer) Close() error {
	return p.Conn.Close()
}

type TCPTransportOpts struct {
	ListenAddress string
	ShakeHands    HandShakeFunc
	Decoder       Decoder

	// notifies the node server when a new Peer connects 
	OnPeerConnect        func(Peer) error
}

type TCPTransport struct {
	TCPTransportOpts
	listener net.Listener
	rpcch    chan RPC
}

func NopeHandShake() error {
	return nil
}

// creates new tcp transport
func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
		rpcch:            make(chan RPC),
	}
}


// dials other peers
// Dial implements the Transport interface
func (t *TCPTransport) Dial(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	go t.handleConnection(conn)

	return nil
}


// ListenAndAccept implements the Transport interface
func (t *TCPTransport) ListenAndAccept() error {
	var err error
	t.listener, err = net.Listen("tcp", t.ListenAddress)
	if err != nil {
		return err
	}

	go t.startAcceptLoop()
	slog.Info("listening on port","port",t.ListenAddress)

	return nil
}

func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if errors.Is(err, net.ErrClosed) {
			return
		}

		slog.Info("incoming new connection","from", conn.RemoteAddr(), "to",conn.LocalAddr())

		go t.handleConnection(conn)
	}
}

func (t *TCPTransport) handleConnection(conn net.Conn) {
	var err error

	defer func() {
		conn.Close()
		slog.Info("dropping peer connection")
	}()

	peer := NewTCPPeer(conn, false)

	if err = t.ShakeHands(peer); err != nil {
		peer.Close()
		fmt.Printf("tcp handshake error : %+v", err)
		return
	}

	if t.OnPeerConnect != nil {
		err = t.OnPeerConnect(peer)
		if err != nil {
			return
		}
	}


	// read loop

	for {

		rpc := RPC{}
		var msgType [1]byte
		peer.Conn.Read(msgType[:])

		if msgType[0] == MessageTypeGOB {
			err = t.Decoder.GOBDecoder(peer.Conn, &rpc)
			if  err != nil {
				slog.Error("gob decoder failed to decode metadata","error", err)
				continue
			}
		}

		var msgType2 [1]byte
		peer.Conn.Read(msgType2[:])
		if msgType2[0] == MessageTypeRaw{
			err = t.Decoder.RawDecoder(peer.Conn,&rpc)
			if err != nil {
				slog.Error("default decoder failed to decode","error",err)
				continue
			}
	
		}
	
		rpc.From = peer.Conn.RemoteAddr()
		t.rpcch <- rpc
	}
}


// consume implements the Transport interface, which will return read-only channnel
// for reading the incoming messages recieved from another peer in the network
func (t *TCPTransport) Consume() <-chan RPC {
	return t.rpcch
}
// Close implements the Transport interface.
func (t *TCPTransport) Close() error {
	return t.listener.Close()
}


