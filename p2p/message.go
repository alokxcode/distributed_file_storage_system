package p2p

import "net"

// RPC holds any arbitory data that is being sent over the each
// transport between the two nodes in the network
type RPC struct {
	From    net.Addr
	Payload []byte
}
