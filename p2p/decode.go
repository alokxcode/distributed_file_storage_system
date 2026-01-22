package p2p

type Decoder interface {
	Decode(Peer) error
}
