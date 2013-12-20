package iproto

import (
	"bytes"
	"testing"
	"time"
)

func TestConnect(t *testing.T) {
	var (
		rid      int32
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
	response, err = conn.Request(rtype, body)
	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}
	if response.RequestType != rtype {
		t.Errorf("Error: requestType should be %d, not %d", rid, response.requestType)
	}

	rtype = 20
	response, err = conn.Request(rtype, body)
	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}
	if response.requestType != rtype {
		t.Errorf("Error: requestType should be %d, not %d", rid, response.requestType)
	}

}
