package p2p

import (
	"encoding/gob"
	"io"
)

type Decoder interface {
	RawDecoder(io.Reader, *RPC) error
	GOBDecoder(io.Reader, *RPC) error
}

type DefaultDecoder struct {}

func (dec DefaultDecoder) GOBDecoder(r io.Reader, rpc *RPC) error {
	err := gob.NewDecoder(r).Decode(&rpc.MetaData)
	 return err 
}


func (dec DefaultDecoder) RawDecoder(r io.Reader, rpc *RPC) error {

	buff := make([]byte, rpc.MetaData.Size)
	n, err := io.ReadFull(r,buff)
	if err != nil {
		return err
	}
	rpc.Payload = buff[:n]
	return nil
}
