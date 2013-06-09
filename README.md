# iProto

Mail.RU iProto protocol implementation on Go.

## Protocol

```
<request> | <response> := <header><body>
<header> = <type:int32><body_length:int32><request_id:int32>
```