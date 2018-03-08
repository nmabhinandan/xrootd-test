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

	//if err := sendPing(con, [2]byte{0xbe, 0xef}); err != nil {
	//	errOut(err)
	//}

	os.Exit(0)
}

func errOut(err error) {
	fmt.Println(err)
	os.Exit(1)
}

func sendHanshake(con net.Conn) error {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, types.NewHandshakeReq()); err != nil {
		return err
	}
	if _, err := con.Write(buf.Bytes()); err != nil {
		return err
	}

	reply := make([]byte, 4096)
	if _, err := con.Read(reply); err != nil {
		return err
	}
	ClientPv = reply[8:12]

	fmt.Printf("Server Protocol Version: %x", binary.BigEndian.Uint32(reply[8:12]))
	fmt.Println()
	fmt.Printf("Server Type: % x ", binary.BigEndian.Uint32(reply[12:16]))
	fmt.Println()

	return nil
}

func sendProtocol(con net.Conn, streamId [2]byte) error {
	buf := new(bytes.Buffer)
	data := types.NewProtocolReq(streamId, 3006, ClientPv)
	if err := binary.Write(buf, binary.BigEndian, data); err != nil {
		return err
	}
	fmt.Printf("Protocol Request: % x \n", data)
	if _, err := con.Write(buf.Bytes()); err != nil {
		return err
	}

	reply := make([]byte, 24)
	if _, err := con.Read(reply); err != nil {
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
	fmt.Printf("Login Request: % x \n", data)
	if _, err := con.Write(buf.Bytes()); err != nil {
		return err
	}

	reply := make([]byte, 42)
	if _, err := con.Read(reply); err != nil {
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
	fmt.Printf("% x \n", data)
	if _, err := con.Write(buf.Bytes()); err != nil {
		return err
	}

	reply := make([]byte, 4)
	if _, err := con.Read(reply); err != nil {
		return err
	}
	fmt.Printf("% x \n", reply)

	return nil
}
