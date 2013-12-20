/*
	Asynchronous mail.ru iproto protocol implementation on Go.
	Thread safe

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
	"github.com/Cergoo/cache"
	"net"
	"runtime"
	"sync/atomic"
	"time"
)

type IProto struct {
	addr        string        //
	connection  *net.TCPConn  //
	requestID   int32         // counter
	chan_writer chan *Request // chanel to wtite
	chan_stop   chan bool     // chanel to stop all gorutines
	requests    cache.Cache   // requests storage
}

type Request struct {
	RequestType int32
	Body        []byte
	Chan        chan<- *Response
}

type Response struct {
	RequestType int32
	Body        []byte
}

// callback function on timeout response, return nil
func callback(key *string, val *interface{}) {
	var (
		ch     chan *Response
		isChan bool
	)
	ch, isChan = (*val).(chan *Response)
	if isChan {
		ch <- nil
	}
}

// constructor
func Connect(addr string, timeout time.Duration) (connection *IProto, err error) {
	raddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return
	}
	conn, err := net.DialTCP("tcp", nil, raddr)
	if err != nil {
		return
	}
	connection = &IProto{
		addr:        addr,
		connection:  conn,
		chan_writer: make(chan *Request, 100),
		chan_stop:   make(chan bool),
		requests:    cache.New(100, false, timeout, callback),
	}

	go connection.read()
	go connection.write()

	// destroy action
	stopAllGorutines := func(t *IProto) {
		close(t.chan_stop)
	}
	runtime.SetFinalizer(connection, stopAllGorutines)

	return
}

// async request
func (conn *IProto) RequestGo(requestType int32, body []byte, ch chan<- *Response) {
	conn.chan_writer <- &Request{
		RequestType: requestType,
		Body:        body,
		Chan:        ch,
	}
}

// sync request
func (conn *IProto) Request(requestType int32, body []byte) *Response {
	ch := make(chan *Response)
	conn.chan_writer <- &Request{
		RequestType: requestType,
		Body:        body,
		Chan:        ch,
	}
	// waiting response
	return <-ch
}

func (conn *IProto) read() {
	var (
		err    error
		ch     chan<- *Response
		isChan bool
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
		select {
		case <-conn.chan_stop:
			return
		default:
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

			ch, isChan = conn.requests.Del(string(header[2])).(chan<- *Response)
			if isChan {
				ch <- &Response{header[0], bodyBuf}
			}
		}
	}
}

func (conn *IProto) write() {
	var (
		err       error
		r         *Request
		requestID int32
		buf       *bytes.Buffer
	)

	for {
		select {
		case <-conn.chan_stop:
			return
		case r = <-conn.chan_writer:
			requestID = atomic.AddInt32(&conn.requestID, 1)
			// write header in a packet
			binary.Write(buf, binary.LittleEndian, []int32{r.RequestType, int32(len(r.Body)), requestID})
			// write body in a packet
			buf.Write(r.Body)
			conn.requests.Set(string(requestID), r.Chan)
			_, err = conn.connection.Write(buf.Bytes())
			if err != nil {
				panic(err)
			}
			buf.Reset()
		}
	}
}
