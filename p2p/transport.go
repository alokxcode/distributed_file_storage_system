package p2p

import "net"

// Peer is an interface that represents a remote node
type Peer interface {
	net.Conn
	Close() error
}

// transport is anything that handles the communication between the nodes in the network.
// THis can be of form (TCP, UDP, websocets, .... )
type Transport interface {
	Dial(string) error
	ListenAndAccept() error
	Consume() <- chan RPC
	Close() error
}
