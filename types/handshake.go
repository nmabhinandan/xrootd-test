package types

import "fmt"

type HandshakeReq interface{}

type HanshakeResp struct {
	Streamid [2]KXR_char
	Status   uint16
	Rlen     KXR_int32
	Pval     KXR_int32
	Flag     KXR_int32
}

func NewHandshakeReq() HandshakeReq {
	return [5]KXR_int32{0, 0, 0, 4, 2012}
}

type ClientReq struct {
	Streamid  [2]KXR_char
	Requestid KXR_unt16
}

type ProtocolReq struct {
	ClientReq
	Clientpv KXR_int32
	Reserved [11]KXR_char
	Options  KXR_int32
	zero     KXR_int32
}

func NewProtocolReq(streamId [2]KXR_char, protocol KXR_unt16) ProtocolReq {
	return ProtocolReq{
		ClientReq: ClientReq{
			Streamid:  streamId,
			Requestid: protocol,
		},
		Clientpv: 0,
		Reserved: [11]KXR_char{},
		Options:  0,
		zero:     0,
	}
}

type LoginReq struct {
	ClientReq

	// TODO: Use Tlen as Token length
	Pid      KXR_int32
	Username [8]KXR_char
	Reserved KXR_char
	Ability  KXR_char
	Capver   [1]KXR_char
	Role     [1]KXR_char
	Tlen     KXR_int32
	Token    [0]KXR_char
}

func NewLoginReq(streamId [2]KXR_char, protocol KXR_unt16, username string) LoginReq {
	fmt.Printf("username: % x \n", username)
	loginReq := LoginReq{
		ClientReq: ClientReq{
			Streamid:  streamId,
			Requestid: protocol,
		},
		Pid:      0,
		Reserved: 0,
		Ability:  0x1,
		Capver:   [1]KXR_char{},
		Role:     [1]KXR_char{},
		Tlen:     0,
		Token:    [0]KXR_char{},
	}

	for i, c := range username {
		loginReq.Username[i] = KXR_char(c)
	}

	return loginReq
}

type PingReq struct {
	Streamid [2]KXR_char
	KXR_ping KXR_unt16
	Reserved KXR_char
	zero     KXR_int32
}

func NewPingReq(streamId [2]KXR_char, KXR_ping KXR_unt16) PingReq {
	pingReq := PingReq{}
	pingReq.Streamid = streamId
	pingReq.KXR_ping = KXR_ping

	return pingReq
}
