# iProto

Asynchronous mail.ru iproto protocol implementation on Go.

## Protocol

```
<request> | <response> := <header><body>
<header> = <type:int32><body_length:int32><request_id:int32>
```

## Usage

```go
package main

import (
	"fmt"
	"github.com/Cergoo/go-iproto"
)

func main() {
	var requestType int32 = 100
	body := []byte("iproto test message")

	conn := iproto.Connect("localhost:33013")
	resp, err = conn.Request(requestType, body)

	fmt.Println("responseBody:", resp.Body)
}
```
