# iProto

Asynchronous mail.ru iproto protocol implementation on Go.
Thread safe.

## Protocol

```
<request> | <response> := <header><body>
<header> = <type:uint32><body_length:uint32><request_id:uint32>
```

## Usage

```go
package main

import (
	"fmt"
	"github.com/Cergoo/go-iproto"
	"time"
)

func main() {
	var requestType uint32 = 100
	body := []byte("iproto test message")

	conn := iproto.Connect("localhost:33013", 2*time.Minute)
	resp := conn.Request(requestType, body)

	fmt.Println("responseBody:", resp.Body)
}
```
