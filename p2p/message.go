package p2p

import "net"


const (
	GETReq byte = 0x1
	POSTReq byte = 0x2

	MessageTypeRaw byte = 0x3
	MessageTypeGOB byte = 0x4
)

type File_MetaData struct {
	Name string
	Size int64
}

// RPC holds any arbitory data that is being sent over the each
// transport between the two nodes in the network
type RPC struct {
	From    net.Addr
	Peer Peer
	ReqType byte
	MetaData File_MetaData
	Payload []byte
}
