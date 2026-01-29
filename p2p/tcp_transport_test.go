package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTCPTransport(t *testing.T) {
	opts := TCPTransportOpts{
		ListenAddress: ":3000",
		ShakeHands:    NOPHandShakeFunc,
		Decoder:       DefaultDecoder{},
	}
	tr := NewTCPTransport(opts)
	assert.Equal(t, tr.ListenAddress, opts.ListenAddress)

	//server

	assert.Nil(t, tr.ListenAndAccept())
}
