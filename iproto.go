package iproto

import (
	"bytes"
	"encoding/binary"
	"net"
)

type IProto struct {
	addr        string
	connection  *net.TCPConn
	requestID   int32
}

type Response struct {
	requestType  int32
	bodyLength   int32
	requestID    int32
	returnCode   int32
	responseBody *bytes.Buffer
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
	connection = &IProto{ addr, conn, 0 }
	return
}

func (conn *IProto) Request(requestType int32, body *bytes.Buffer) (response *Response, err error) {
	// Prepare packet
	packet, err := conn.pack(requestType, body)
	if err != nil {
		return
	}

	// Send it
	err = conn.send(packet)
	if err != nil {
		return
	}

	//Wait for answer
	response, err = conn.recv()
	return
}

func (conn *IProto) pack(requestType int32, body *bytes.Buffer) (packet *bytes.Buffer, err error) {	
	// Each request should have uniq RequestID
	conn.requestID++
	// And t should not get out of 32 bits
	if conn.requestID >= 1<<31 - 1 {
		conn.requestID = 1
	}
	bodyLength := int32(body.Len())

	// <header> ::= <type><body_length><request_id>
	header := [] int32 { requestType, bodyLength, conn.requestID }

	// Put integers into packet
	packet = new(bytes.Buffer)
	err = binary.Write(packet, binary.LittleEndian, header)
	if err != nil {
		return
	}

	// Put body into packet
	_, err = packet.Write(body.Bytes())
	return
}

func (conn *IProto) send(packet *bytes.Buffer) (err error) {
	_, err = conn.connection.Write(packet.Bytes())
	return
}

func (conn *IProto) recv() (response *Response, err error) {
	headerBuf := make([]byte, 16)
	// Read header 16 bytes (12 + responseCode)
	_, err = conn.connection.Read(headerBuf)
	if err != nil {
		return
	}

	// Unpack data
	res := make([]int32, 4)
	err = binary.Read(bytes.NewBuffer(headerBuf), binary.LittleEndian, &res)

	requestType  := res[0]
	bodyLength   := res[1]
	requestID    := res[2]
	responseCode := res[3]

	// Read body
	bodyRest := bodyLength - 4
	bodyBuf  := make([]byte, bodyRest)
	if bodyRest > 0 {
		_, err = conn.connection.Read(bodyBuf)
		if err != nil {
			return
		}
	}

	response = &Response{ requestType, bodyLength, requestID, responseCode, bytes.NewBuffer(bodyBuf) }
	return
}
