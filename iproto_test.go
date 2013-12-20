package iproto

import (
	"testing"
	"time"
)

func TestConnect(t *testing.T) {
	var (
		rtype    uint32
		response *Response
		err      error
		conn     *IProto
	)

	conn, err = Connect("localhost:33013", 5*time.Minute)
	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}

	body := []byte("iproto test message")
	rtype = 17
	response = conn.Request(rtype, body)
	if response.RequestType != rtype {
		t.Errorf("Error: requestType should be %d, not %d", rtype, response.RequestType)
	}

	rtype = 20
	response = conn.Request(rtype, body)
	if response.RequestType != rtype {
		t.Errorf("Error: requestType should be %d, not %d", rtype, response.RequestType)
	}

}
