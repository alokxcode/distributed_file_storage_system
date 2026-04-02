package p2p

// Peer is an interface that represents a remote node
type Peer interface {
	Close() error
}

// transport is anything that handles the communication between the nodes in the network.
// THis can be of form (TCP, UDP, websocets, .... )
type Transport interface {
	ListenAndAccept() error
	Consume() <- chan RPC
}
