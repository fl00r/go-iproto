/*
	Asynchronous mail.ru iproto protocol implementation on Go.
    NOT thread safe

	Protocol description
	<request> | <response> := <header><body>
	<header> = <type:int32><body_length:int32><request_id:int32>

	(c) 2013 Cergoo (forked from fl00r/go-iproto)
	under terms of ISC license

*/
package iproto

import (
	"bytes"
	"encoding/binary"
	"net"
)

type IProto struct {
	addr       string
	connection *net.TCPConn
	requestID  int32
	requests   map[int32]chan *Response
	writeChan  chan []byte
}

type Response struct {
	RequestType int32
	Body        []byte
}

func Connect(addr string) (connection *IProto, err error) {
	raddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return
	}
	conn, err := net.DialTCP("tcp", nil, raddr)
	if err != nil {
		return
	}
	connection = &IProto{
		addr:       addr,
		connection: conn,
		requests:   make(map[int32]chan *Response),
		writeChan:  make(chan []byte),
	}

	go connection.read()
	go connection.write()

	return
}

func (conn *IProto) Request(requestType int32, body []byte) (response *Response, err error) {
	// create packet
	packet := new(bytes.Buffer)
	conn.requestID++
	// write header in a packet
	err = binary.Write(packet, binary.LittleEndian, []int32{requestType, int32(len(body)), conn.requestID})
	if err != nil {
		return
	}
	// write body in a packet
	_, err = packet.Write(body)
	if err != nil {
		return
	}

	conn.requests[conn.requestID] = make(chan *Response)
	// send request
	conn.writeChan <- packet.Bytes()
	// waiting response
	response = <-conn.requests[conn.requestID]
	// delete chanel
	delete(conn.requests, conn.requestID)
	return
}

func (conn *IProto) read() {
	var (
		err error
	)
	/*
		requestType = header[0]
		bodyLength  = header[1]
		requestID   = header[2]
	*/
	header := make([]int32, 3)
	headerBuf := make([]byte, 12)
	headerReader := bytes.NewReader(headerBuf)
	for {
		// read header
		_, err = conn.connection.Read(headerBuf)
		if err != nil {
			panic(err)
		}
		err = binary.Read(headerReader, binary.LittleEndian, &header)
		if err != nil {
			panic(err)
		}
		// read body
		bodyBuf := make([]byte, header[1])
		_, err = conn.connection.Read(bodyBuf)
		if err != nil {
			panic(err)
		}

		conn.requests[header[2]] <- &Response{header[0], bodyBuf}
	}
}

func (conn *IProto) write() {
	var (
		err error
	)
	for {
		_, err = conn.connection.Write(<-conn.writeChan)
		if err != nil {
			panic(err)
		}
	}
}
