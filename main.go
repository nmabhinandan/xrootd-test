package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"

	"xrootd-test-client/types"
)

func main() {

	con, _ := net.Dial("tcp", "localhost:9001")

	defer con.Close()

	if err := sendHanshake(con); err != nil {
		os.Exit(1)
	}

	if err := sendProtocol(con, [2]byte{0xbe, 0xef}); err != nil {
		os.Exit(1)
	}

	if err := sendLogin(con, [2]byte{0xbe, 0xef}); err != nil {
		os.Exit(1)
	}

	//if err := sendPing(con, [2]byte{0xbe, 0xef}); err != nil {
	//	os.Exit(1)
	//}

	os.Exit(0)
}

func sendHanshake(con net.Conn) error {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, types.NewHandshakeReq())
	con.Write(buf.Bytes())
	reply := make([]byte, 4096)
	con.Read(reply)

	fmt.Printf("Server Protocol Version: %d", binary.BigEndian.Uint32(reply[8:12]))
	fmt.Println()
	fmt.Printf("Server Type: %d ", binary.BigEndian.Uint32(reply[12:16]))
	fmt.Println()

	return nil
}

func sendProtocol(con net.Conn, streamId [2]byte) error {
	buf := new(bytes.Buffer)
	data := types.NewProtocolReq(streamId, 3006)
	binary.Write(buf, binary.BigEndian, data)
	fmt.Printf("Protocol Request: % x \n", data)
	con.Write(buf.Bytes())

	reply := make([]byte, 24)
	con.Read(reply)
	fmt.Printf("Protocol Reply: % x \n", reply)

	return nil
}

func sendLogin(con net.Conn, streamID [2]byte) error {
	buf := new(bytes.Buffer)
	data := types.NewLoginReq(streamID, 3007, "gopher")
	binary.Write(buf, binary.BigEndian, data)
	fmt.Printf("Login Request: % x \n", data)
	con.Write(buf.Bytes())

	reply := make([]byte, 42)
	con.Read(reply)
	fmt.Printf("Login Reply: % x \n", reply)

	return nil
}

func sendPing(con net.Conn, streamID [2]byte) error {
	buf := new(bytes.Buffer)
	data := types.NewPingReq(streamID, 3011)
	binary.Write(buf, binary.BigEndian, data)
	fmt.Printf("% x \n", data)
	con.Write(buf.Bytes())

	reply := make([]byte, 4)
	con.Read(reply)
	fmt.Printf("% x \n", reply)

	return nil
}
