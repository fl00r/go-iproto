package iproto

import (
	"testing"
	"bytes"
)

func TestConnect(t *testing.T) {
	conn, err := Connect("localhost:33013")
	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}

	var rid int32 = 17
	response, err := conn.Request(rid, new(bytes.Buffer))
	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}
	if response.requestType != rid {
		t.Errorf("Error: requestType should be %d, not %d", rid, response.requestType)
	}
	if response.requestID != 1 {
		t.Errorf("Error: requestID should be %d, not %d", 1, response.requestID)
	}
}