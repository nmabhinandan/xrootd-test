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

func readXRootResp(con net.Conn) ([]byte, error) {
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

func sendHanshake(con net.Conn) error {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, types.NewHandshakeReq()); err != nil {
		return err
	}
	if _, err := con.Write(buf.Bytes()); err != nil {
		return err
	}
	//con.SetReadDeadline(time.Now().Add(time.Second * 30))

	reply, err := readXRootResp(con)
	if err != nil {
		return err
	}

	//ClientPv = reply[8:12]
	fmt.Printf("reply: % x\n", reply)
	//fmt.Printf("Server Protocol Version: %x", binary.BigEndian.Uint32(reply[8:12]))
	//fmt.Printf("Server Type: % x ", binary.BigEndian.Uint32(reply[12:16]))

	return nil
}

func sendProtocol(con net.Conn, streamId [2]byte) error {
	buf := new(bytes.Buffer)
	data := types.NewProtocolReq(streamId, 3006, ClientPv)
	if err := binary.Write(buf, binary.BigEndian, data); err != nil {
		return err
	}
	//fmt.Printf("Protocol Request: % x \n", data)
	//fmt.Println("Length of protocol request: ", len(buf.Bytes()))

	if _, err := con.Write(buf.Bytes()); err != nil {
		return err
	}
	reply, err := readXRootResp(con)
	if err != nil {
		return err
	}
	fmt.Printf("Protocol Reply: % x \n", reply)

	return nil
}

func sendLogin(con net.Conn, streamID [2]byte) error {
	buf := new(bytes.Buffer)
	data := types.NewLoginReq(streamID, 3007, "gopher")
	if err := binary.Write(buf, binary.BigEndian, data); err != nil {
		return err
	}
	//fmt.Printf("Login Request: % x \n", data)
	//fmt.Println("Length of login request: ", len(buf.Bytes()))
	if _, err := con.Write(buf.Bytes()); err != nil {
		return err
	}

	reply, err := readXRootResp(con)
	if err != nil {
		return err
	}
	fmt.Printf("Login Reply: % x \n", reply)

	return nil
}

func sendPing(con net.Conn, streamID [2]byte) error {
	buf := new(bytes.Buffer)
	data := types.NewPingReq(streamID, 3011)
	if err := binary.Write(buf, binary.BigEndian, data); err != nil {
		return err
	}
	//fmt.Printf("Ping Request: % x \n", data)
	//fmt.Println("Length of ping request: ", len(buf.Bytes()))

	if _, err := con.Write(buf.Bytes()); err != nil {
		return err
	}

	reply, err := readXRootResp(con)
	if err != nil {
		return err
	}
	fmt.Printf("Ping Reply: % x \n", reply)

	return nil
}

func sendInvalidRequest(con net.Conn, streamID [2]byte) error {
	buf := new(bytes.Buffer)
	data := struct {
		StreamId  [2]byte
		RequestId uint16
		Params    [16]byte
		Dlen      int32
	}{
		StreamId:  streamID,
		RequestId: 0,
	}

	if err := binary.Write(buf, binary.BigEndian, data); err != nil {
		return err
	}
	//fmt.Printf("Invalid Request: % x \n", data)
	//fmt.Println("Length of invalid request: ", len(buf.Bytes()))

	if _, err := con.Write(buf.Bytes()); err != nil {
		return err
	}

	reply, err := readXRootResp(con)
	if err != nil {
		return err
	}
	fmt.Printf("Reply to invalid request: % x \n", reply)
	fmt.Printf("Error code: %d\n", binary.BigEndian.Uint16(reply[2:4]))

	return nil
}
