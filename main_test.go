package main

import (
	"fmt"
	"net"
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestProgram(t *testing.T) {
	client, server := net.Pipe()
	defer client.Close()
	defer server.Close()

	errC := make(chan error)
	dataC := make(chan []byte)
	var wg sync.WaitGroup

	go actServer(server, errC, dataC, &wg)

	go func() {
		for {
			if err := <-errC; err != nil {
				fmt.Println("error: ", err)
				t.Error(err)
				break
			}
		}
	}()

	streamId := [2]byte{0xbe, 0xef}
	t.Run("Test Handshake", func(t *testing.T) {
		wg.Add(1)
		go func() {
			if err := sendHanshake(client); err != nil {
				t.Error(err)
			}
		}()
		time.Sleep(100 * time.Millisecond)

		req := <-dataC
		correctReq, _ := encodePayload([]int32{0, 0, 0, 4, 2012, 0})
		if !reflect.DeepEqual(req, correctReq) {
			t.Error("Invalid handshake request")
		}

		server.Write([]byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x8, 0x0, 0x0, 0x3, 0x10, 0x0, 0x0, 0x0, 0x1})
	})

	t.Run("Test Protocol", func(t *testing.T) {
		wg.Add(1)
		go func() {
			if err := sendProtocol(client, streamId); err != nil {
				t.Error(err)
			}
		}()
		time.Sleep(100 * time.Millisecond)

		req := <-dataC
		if err := testHeader(req, streamId, 3006); err != nil {
			t.Error("Invalid header: ", err.Error())
		}
		if !reflect.DeepEqual(req[4:8], []byte{0x0, 0x0, 0x3, 0x10}) {
			t.Error("Invalid clientpv: expected: ", []byte{0x0, 0x0, 0x3, 0x10}, " but found ", req[4:8])
		}

		server.Write([]byte{0xbe, 0xef, 0x0, 0x0, 0x0, 0x0, 0x0, 0x8, 0x0, 0x0, 0x3, 0x10, 0x0, 0x0, 0x0, 0x1})
	})

	t.Run("Test Login", func(t *testing.T) {
		wg.Add(1)
		go func() {
			if err := sendLogin(client, streamId); err != nil {
				t.Error(err)
			}
		}()
		time.Sleep(100 * time.Millisecond)

		req := <-dataC
		if err := testHeader(req, streamId, 3007); err != nil {
			t.Error("Invalid header: ", err.Error())
		}

		u := make([]byte, 8)
		for i, c := range "gopher" {
			u[i] = byte(c)
		}
		if !reflect.DeepEqual(req[8:16], u) {
			t.Error("Invalid username")
		}

		server.Write([]byte{0xbe, 0xef, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0})
	})

	t.Run("Test Ping", func(t *testing.T) {
		wg.Add(1)
		go func() {
			if err := sendPing(client, streamId); err != nil {
				t.Error(err)
			}
		}()
		time.Sleep(100 * time.Millisecond)

		req := <-dataC
		if err := testHeader(req, streamId, 3011); err != nil {
			t.Error("Invalid header: ", err.Error())
		}

		server.Write([]byte{0xbe, 0xef, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0})
	})

	t.Run("Test Invalid Request", func(t *testing.T) {
		wg.Add(1)
		go func() {
			if err := sendInvalidRequest(client, streamId); err != nil {
				t.Error(err)
			}
		}()
		time.Sleep(100 * time.Millisecond)

		<-dataC

		server.Write([]byte{0xbe, 0xef, 0xf, 0xa3, 0x0, 0x0, 0x0, 0x22, 0x0, 0x0, 0xb, 0xb9, 0x52, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65, 0x64, 0x20, 0x61, 0x72, 0x67, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x20, 0x6e, 0x6f, 0x74, 0x20, 0x70, 0x72, 0x65, 0x73, 0x65, 0x6e, 0x74, 0x0})
	})

	wg.Wait()
}

func testHeader(req []byte, streamId [2]byte, reqId uint16) error {
	s, _ := encodePayload(streamId)
	if !reflect.DeepEqual(req[0:2], s) {
		return fmt.Errorf("Invalid stream id")
	}
	r, _ := encodePayload([]uint16{reqId})
	if !reflect.DeepEqual(req[2:4], r) {
		return fmt.Errorf("Invalid request id")
	}
	return nil
}

func actServer(con net.Conn, errC chan error, dataC chan []byte, wg *sync.WaitGroup) {
	for {
		req, err := readXRootDReq(con)
		if err != nil {
			errC <- err
			break
		}
		dataC <- req
		wg.Done()
	}
}

func readXRootDReq(con net.Conn) ([]byte, error) {
	header := make([]byte, 24)
	if _, err := con.Read(header); err != nil {
		return nil, err
	}

	//if len(header) != 24 {
	//	return nil, fmt.Errorf("Header length is %d. Expected length is 24", len(header))
	//} // This test is useless

	dlen := header[23]
	if dlen == 0 {
		return header, nil
	}
	data := make([]byte, dlen)
	if _, err := con.Read(data); err != nil {
		return nil, err
	}

	return append(header, data...), nil
}
