package types

import "fmt"

type HandshakeReq interface{}

type HanshakeResp struct {
	Streamid [2]byte
	Status   uint16
	Rlen     int32
	Pval     int32
	Flag     int32
}

func NewHandshakeReq() HandshakeReq {
	return [5]int32{0, 0, 0, 4, 2012}
}

type ClientReq struct {
	Streamid  [2]byte
	Requestid uint16
}

type ProtocolReq struct {
	ClientReq
	Clientpv int32
	Reserved [11]byte
	Options  int32
	zero     int32
}

func NewProtocolReq(streamId [2]byte, protocol uint16) ProtocolReq {
	return ProtocolReq{
		ClientReq: ClientReq{
			Streamid:  streamId,
			Requestid: protocol,
		},
		Clientpv: 0,
		Reserved: [11]byte{},
		Options:  0,
		zero:     0,
	}
}

type LoginReq struct {
	ClientReq

	// TODO: Use Tlen as Token length
	Pid      int32
	Username [8]byte
	Reserved byte
	Ability  byte
	Capver   [1]byte
	Role     [1]byte
	Tlen     int32
	Token    [0]byte
}

func NewLoginReq(streamId [2]byte, protocol uint16, username string) LoginReq {
	fmt.Printf("username: % x \n", username)
	loginReq := LoginReq{
		ClientReq: ClientReq{
			Streamid:  streamId,
			Requestid: protocol,
		},
		Pid:      0,
		Reserved: 0,
		Ability:  0x1,
		Capver:   [1]byte{},
		Role:     [1]byte{},
		Tlen:     0,
		Token:    [0]byte{},
	}

	for i, c := range username {
		loginReq.Username[i] = byte(c)
	}

	return loginReq
}

type PingReq struct {
	Streamid [2]byte
	KXR_ping uint16
	Reserved byte
	zero     int32
}

func NewPingReq(streamId [2]byte, KXR_ping uint16) PingReq {
	pingReq := PingReq{}
	pingReq.Streamid = streamId
	pingReq.KXR_ping = KXR_ping

	return pingReq
}
