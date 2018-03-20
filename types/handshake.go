package types

import (
	"bytes"
	"encoding/binary"
)

type HandshakeReq interface{}

func NewHandshakeReq() HandshakeReq {
	return [5]int32{0, 0, 0, 4, 2012}
}

type ClientReq struct {
	StreamId  [2]byte
	RequestId uint16
}

type ProtocolReq struct {
	ClientReq
	Clientpv int32
	Reserved [11]byte
	Options  byte
	zero     int32
}

func NewProtocolReq(streamId [2]byte, requestId uint16, clientPv []byte) ProtocolReq {
	protoReq := ProtocolReq{
		ClientReq: ClientReq{
			StreamId:  streamId,
			RequestId: requestId,
		},
		Reserved: [11]byte{},
		Options:  0,
		zero:     0,
	}
	binary.Read(bytes.NewBuffer(clientPv), binary.BigEndian, &protoReq.Clientpv)
	return protoReq
}

type LoginReq struct {
	ClientReq

	// TODO: Use Dlen as Token length
	Pid      int32
	Username [8]byte
	Reserved byte
	Ability  byte
	Capver   [1]byte
	Role     [1]byte
	Dlen     int32
	//Token    [0]byte
}

func NewLoginReq(streamId [2]byte, protocol uint16, username string) LoginReq {
	loginReq := LoginReq{
		ClientReq: ClientReq{
			StreamId:  streamId,
			RequestId: protocol,
		},
		Pid:      0,
		Reserved: 0,
		Ability:  0x1,
		Capver:   [1]byte{},
		Role:     [1]byte{},
		Dlen:     0,
		//Token:    [0]byte{},
	}

	for i, c := range username {
		if i >= len(loginReq.Username) {
			break
		}
		loginReq.Username[i] = byte(c)
	}

	return loginReq
}

type PingReq struct {
	ClientReq
	Reserved [16]byte
	dlen     int32
}

func NewPingReq(streamId [2]byte, requestId uint16) PingReq {
	pingReq := PingReq{
		ClientReq: ClientReq{
			StreamId:  streamId,
			RequestId: requestId,
		},
		Reserved: [16]byte{},
		dlen:     0,
	}

	return pingReq
}
