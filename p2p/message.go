package p2p

import "net"


const (
	MessageTypeRaw byte = 0x1
	MessageTypeGOB byte = 0x2
)

type File_MetaData struct {
	Name string
	Size int64
}

// RPC holds any arbitory data that is being sent over the each
// transport between the two nodes in the network
type RPC struct {
	From    net.Addr
	MetaData File_MetaData
	Payload []byte
}
