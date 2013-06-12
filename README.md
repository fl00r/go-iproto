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
	"github.com/fl00r/go-iproto"
	"bytes"
)

func main() {
	var requestID int32 = 100
	body := new(bytes.Buffer)

	conn := iproto.Connect("localhost:33013")
	resp, err = conn.Request(requestID, body)

	fmt.Println("responseBody:", resp.Body)
}
```