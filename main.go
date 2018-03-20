package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"

	"xrootd-test-client/types"
)

var ClientPv []byte

func main() {

	con, _ := net.Dial("tcp", "localhost:9001")
	//con, _ := net.Dial("tcp", "ccxrootdgotest.in2p3.fr:9001")

	defer con.Close()

	if err := sendHanshake(con); err != nil {
		errOut(err)
	}

	if err := sendProtocol(con, [2]byte{0xbe, 0xef}); err != nil {
		errOut(err)
	}

	if err := sendLogin(con, [2]byte{0xbe, 0xef}); err != nil {
		errOut(err)
	}

	if err := sendPing(con, [2]byte{0xbe, 0xef}); err != nil {
		errOut(err)
	}

	if err := sendInvalidRequest(con, [2]byte{0xbe, 0xef}); err != nil {
		errOut(err)
	}

	os.Exit(0)
}

func errOut(err error) {
	fmt.Println(err)
	os.Exit(1)
}

func readXRootDResp(con net.Conn) ([]byte, error) {
	header := make([]byte, 8)
	if _, err := con.Read(header); err != nil {
		return nil, err
	}

	dlen := header[7]
	data := make([]byte, dlen)
	if _, err := con.Read(data); err != nil {
		return nil, err
	}

	return append(header, data...), nil
}

func encodePayload(payload interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, payload); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func sendHanshake(con net.Conn) error {
	payload, err := encodePayload(types.NewHandshakeReq())
	if err != nil {
		return err
	}
	if _, err := con.Write(payload); err != nil {
		return err
	}

	reply, err := readXRootDResp(con)
	if err != nil {
		return err
	}

	ClientPv = reply[8:12]

	//fmt.Printf("reply: %#v\n", reply)
	fmt.Printf("Server Protocol Version: %x\n", binary.BigEndian.Uint32(reply[8:12]))
	fmt.Printf("Server Type: % x \n", binary.BigEndian.Uint32(reply[12:16]))

	return nil
}

func sendProtocol(con net.Conn, streamId [2]byte) error {
	data := types.NewProtocolReq(streamId, 3006, ClientPv)
	payload, err := encodePayload(data)
	if err != nil {
		return err
	}
	if _, err := con.Write(payload); err != nil {
		return err
	}

	reply, err := readXRootDResp(con)
	if err != nil {
		return err
	}
	fmt.Printf("Protocol Reply: % x \n", reply)

	return nil
}

func sendLogin(con net.Conn, streamID [2]byte) error {
	data := types.NewLoginReq(streamID, 3007, "gopher")
	payload, err := encodePayload(data)
	if err != nil {
		return err
	}
	if _, err := con.Write(payload); err != nil {
		return err
	}

	reply, err := readXRootDResp(con)
	if err != nil {
		return err
	}
	fmt.Printf("Login Reply: % x \n", reply)

	return nil
}

func sendPing(con net.Conn, streamID [2]byte) error {
	data := types.NewPingReq(streamID, 3011)
	payload, err := encodePayload(data)
	if err != nil {
		return err
	}
	if _, err := con.Write(payload); err != nil {
		return err
	}

	reply, err := readXRootDResp(con)
	if err != nil {
		return err
	}
	fmt.Printf("Ping Reply: % x \n", reply)

	return nil
}

func sendInvalidRequest(con net.Conn, streamID [2]byte) error {
	data := struct {
		StreamId  [2]byte
		RequestId uint16
		Params    [16]byte
		Dlen      int32
	}{
		StreamId:  streamID,
		RequestId: 0,
	}

	payload, err := encodePayload(data)
	if err != nil {
		return err
	}
	if _, err := con.Write(payload); err != nil {
		return err
	}

	reply, err := readXRootDResp(con)
	if err != nil {
		return err
	}
	fmt.Printf("Reply to invalid request: % x \n", reply)
	fmt.Printf("Error code: %d\n", binary.BigEndian.Uint16(reply[2:4]))

	return nil
}
