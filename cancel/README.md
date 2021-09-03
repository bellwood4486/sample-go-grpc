# Cancel Sample

## How to run

Run server.
```shell
GODEBUG=http2debug=2 go run server/main.go
```

Run client.
```shell
GODEBUG=http2debug=2 go run client/main.go
```

## Results

### Summary

Disable server side cancel:

|at|client|server|
|---|---|---|
|2021/09/03 22:51:00|start cancel sample|sleep for 2s...|
|2021/09/03 22:51:01|could not sleep: rpc error: code = Canceled desc = context canceled|-|
|2021/09/03 22:51:02|-|sleep success|
|2021/09/03 22:51:03|start timeout sample|sleep for 2s...|
|2021/09/03 22:51:04|could not sleep: rpc error: code = DeadlineExceeded desc = context deadline exceeded|-|
|2021/09/03 22:51:05|-|sleep success

Enable server side cancel:

|at|client|server|
|---|---|---|
|2021/09/03 22:51:06|start cancel sample|sleep for 2s...|
|2021/09/03 22:51:07|could not sleep: rpc error: code = Canceled desc = context canceled|sleep canceled: context canceled|
|2021/09/03 22:51:09|start timeout sample|sleep for 2s...|
|2021/09/03 22:51:10|could not sleep: rpc error: code = DeadlineExceeded desc = context deadline exceeded|sleep canceled: context deadline exceeded|

### Console

client
```
❯ GODEBUG=http2debug=2 go run client/main.go
//
// client cancel & server cancel:OFF
//
2021/09/03 23:09:59 http2: Framer 0xc000224000: wrote SETTINGS len=0
2021/09/03 23:09:59 http2: Framer 0xc000224000: read SETTINGS len=6, settings: MAX_FRAME_SIZE=16384
2021/09/03 23:09:59 start cancel sample
2021/09/03 23:09:59 http2: Framer 0xc000224000: wrote SETTINGS flags=ACK len=0
2021/09/03 23:09:59 http2: Framer 0xc000224000: read SETTINGS flags=ACK len=0
2021/09/03 23:09:59 http2: Framer 0xc000224000: wrote HEADERS flags=END_HEADERS stream=1 len=70
2021/09/03 23:09:59 http2: Framer 0xc000224000: wrote DATA flags=END_STREAM stream=1 len=7 data="\x00\x00\x00\x00\x02\b\x02"
2021/09/03 23:09:59 http2: Framer 0xc000224000: read WINDOW_UPDATE len=4 (conn) incr=7
2021/09/03 23:09:59 http2: Framer 0xc000224000: read PING len=8 ping="\x02\x04\x10\x10\t\x0e\a\a"
2021/09/03 23:09:59 http2: Framer 0xc000224000: wrote PING flags=ACK len=8 ping="\x02\x04\x10\x10\t\x0e\a\a"
2021/09/03 23:10:00 could not sleep: rpc error: code = Canceled desc = context canceled
2021/09/03 23:10:00 http2: Framer 0xc000224000: wrote RST_STREAM stream=1 len=4 ErrCode=CANCEL
//
// client timeout & server cancel:OFF
//
2021/09/03 23:10:02 start timeout sample
2021/09/03 23:10:02 http2: Framer 0xc000224000: wrote HEADERS flags=END_HEADERS stream=3 len=25
2021/09/03 23:10:02 http2: Framer 0xc000224000: wrote DATA flags=END_STREAM stream=3 len=7 data="\x00\x00\x00\x00\x02\b\x02"
2021/09/03 23:10:02 http2: Framer 0xc000224000: read WINDOW_UPDATE len=4 (conn) incr=7
2021/09/03 23:10:02 http2: Framer 0xc000224000: read PING len=8 ping="\x02\x04\x10\x10\t\x0e\a\a"
2021/09/03 23:10:02 http2: Framer 0xc000224000: wrote PING flags=ACK len=8 ping="\x02\x04\x10\x10\t\x0e\a\a"
2021/09/03 23:10:03 could not sleep: rpc error: code = DeadlineExceeded desc = context deadline exceeded
2021/09/03 23:10:03 http2: Framer 0xc000224000: wrote RST_STREAM stream=3 len=4 ErrCode=CANCEL
//
// client cancel & server cancel:ON
//
2021/09/03 23:10:05 start cancel sample
2021/09/03 23:10:05 http2: Framer 0xc000224000: wrote HEADERS flags=END_HEADERS stream=5 len=7
2021/09/03 23:10:05 http2: Framer 0xc000224000: wrote DATA flags=END_STREAM stream=5 len=9 data="\x00\x00\x00\x00\x04\b\x02\x10\x01"
2021/09/03 23:10:05 http2: Framer 0xc000224000: read WINDOW_UPDATE len=4 (conn) incr=9
2021/09/03 23:10:05 http2: Framer 0xc000224000: read PING len=8 ping="\x02\x04\x10\x10\t\x0e\a\a"
2021/09/03 23:10:05 http2: Framer 0xc000224000: wrote PING flags=ACK len=8 ping="\x02\x04\x10\x10\t\x0e\a\a"
2021/09/03 23:10:06 could not sleep: rpc error: code = Canceled desc = context canceled
2021/09/03 23:10:06 http2: Framer 0xc000224000: wrote RST_STREAM stream=5 len=4 ErrCode=CANCEL
//
// client timeout & server cancel:ON
//
2021/09/03 23:10:08 start timeout sample
2021/09/03 23:10:08 http2: Framer 0xc000224000: wrote HEADERS flags=END_HEADERS stream=7 len=15
2021/09/03 23:10:08 http2: Framer 0xc000224000: wrote DATA flags=END_STREAM stream=7 len=9 data="\x00\x00\x00\x00\x04\b\x02\x10\x01"
2021/09/03 23:10:08 http2: Framer 0xc000224000: read WINDOW_UPDATE len=4 (conn) incr=9
2021/09/03 23:10:08 http2: Framer 0xc000224000: read PING len=8 ping="\x02\x04\x10\x10\t\x0e\a\a"
2021/09/03 23:10:08 http2: Framer 0xc000224000: wrote PING flags=ACK len=8 ping="\x02\x04\x10\x10\t\x0e\a\a"
2021/09/03 23:10:09 could not sleep: rpc error: code = DeadlineExceeded desc = context deadline exceeded
2021/09/03 23:10:09 http2: Framer 0xc000224000: wrote RST_STREAM stream=7 len=4 ErrCode=CANCEL
```

server
```
❯ GODEBUG=http2debug=2 go run server/main.go
//
// client cancel & server cancel:OFF
//
2021/09/03 23:09:59 http2: Framer 0xc0001e0000: wrote SETTINGS len=6, settings: MAX_FRAME_SIZE=16384
2021/09/03 23:09:59 http2: Framer 0xc0001e0000: read SETTINGS len=0
2021/09/03 23:09:59 http2: Framer 0xc0001e0000: wrote SETTINGS flags=ACK len=0
2021/09/03 23:09:59 http2: Framer 0xc0001e0000: read SETTINGS flags=ACK len=0
2021/09/03 23:09:59 http2: Framer 0xc0001e0000: read HEADERS flags=END_HEADERS stream=1 len=70
2021/09/03 23:09:59 http2: decoded hpack field header field ":method" = "POST"
2021/09/03 23:09:59 http2: decoded hpack field header field ":scheme" = "http"
2021/09/03 23:09:59 http2: decoded hpack field header field ":path" = "/helloworld.Greeter/Sleep"
2021/09/03 23:09:59 http2: decoded hpack field header field ":authority" = "localhost:8080"
2021/09/03 23:09:59 http2: decoded hpack field header field "content-type" = "application/grpc"
2021/09/03 23:09:59 http2: decoded hpack field header field "user-agent" = "grpc-go/1.40.0"
2021/09/03 23:09:59 http2: decoded hpack field header field "te" = "trailers"
2021/09/03 23:09:59 http2: Framer 0xc0001e0000: read DATA flags=END_STREAM stream=1 len=7 data="\x00\x00\x00\x00\x02\b\x02"
2021/09/03 23:09:59 http2: Framer 0xc0001e0000: wrote WINDOW_UPDATE len=4 (conn) incr=7
2021/09/03 23:09:59 http2: Framer 0xc0001e0000: wrote PING len=8 ping="\x02\x04\x10\x10\t\x0e\a\a"
2021/09/03 23:09:59 http2: Framer 0xc0001e0000: read PING flags=ACK len=8 ping="\x02\x04\x10\x10\t\x0e\a\a"
2021/09/03 23:09:59 sleep for 2s...
2021/09/03 23:10:00 http2: Framer 0xc0001e0000: read RST_STREAM stream=1 len=4 ErrCode=CANCEL
2021/09/03 23:10:01 sleep success
//
// client timeout & server cancel:OFF
//
2021/09/03 23:10:02 http2: Framer 0xc0001e0000: read HEADERS flags=END_HEADERS stream=3 len=25
2021/09/03 23:10:02 http2: decoded hpack field header field ":method" = "POST"
2021/09/03 23:10:02 http2: decoded hpack field header field ":scheme" = "http"
2021/09/03 23:10:02 http2: decoded hpack field header field ":path" = "/helloworld.Greeter/Sleep"
2021/09/03 23:10:02 http2: decoded hpack field header field ":authority" = "localhost:8080"
2021/09/03 23:10:02 http2: decoded hpack field header field "content-type" = "application/grpc"
2021/09/03 23:10:02 http2: decoded hpack field header field "user-agent" = "grpc-go/1.40.0"
2021/09/03 23:10:02 http2: decoded hpack field header field "te" = "trailers"
2021/09/03 23:10:02 http2: decoded hpack field header field "grpc-timeout" = "999949u"
2021/09/03 23:10:02 http2: Framer 0xc0001e0000: read DATA flags=END_STREAM stream=3 len=7 data="\x00\x00\x00\x00\x02\b\x02"
2021/09/03 23:10:02 http2: Framer 0xc0001e0000: wrote WINDOW_UPDATE len=4 (conn) incr=7
2021/09/03 23:10:02 http2: Framer 0xc0001e0000: wrote PING len=8 ping="\x02\x04\x10\x10\t\x0e\a\a"
2021/09/03 23:10:02 sleep for 2s...
2021/09/03 23:10:02 http2: Framer 0xc0001e0000: read PING flags=ACK len=8 ping="\x02\x04\x10\x10\t\x0e\a\a"
2021/09/03 23:10:03 http2: Framer 0xc0001e0000: read RST_STREAM stream=3 len=4 ErrCode=CANCEL
2021/09/03 23:10:04 sleep success
//
// client cancel & server cancel:ON
//
2021/09/03 23:10:05 http2: Framer 0xc0001e0000: read HEADERS flags=END_HEADERS stream=5 len=7
2021/09/03 23:10:05 http2: decoded hpack field header field ":method" = "POST"
2021/09/03 23:10:05 http2: decoded hpack field header field ":scheme" = "http"
2021/09/03 23:10:05 http2: decoded hpack field header field ":path" = "/helloworld.Greeter/Sleep"
2021/09/03 23:10:05 http2: decoded hpack field header field ":authority" = "localhost:8080"
2021/09/03 23:10:05 http2: decoded hpack field header field "content-type" = "application/grpc"
2021/09/03 23:10:05 http2: decoded hpack field header field "user-agent" = "grpc-go/1.40.0"
2021/09/03 23:10:05 http2: decoded hpack field header field "te" = "trailers"
2021/09/03 23:10:05 http2: Framer 0xc0001e0000: read DATA flags=END_STREAM stream=5 len=9 data="\x00\x00\x00\x00\x04\b\x02\x10\x01"
2021/09/03 23:10:05 http2: Framer 0xc0001e0000: wrote WINDOW_UPDATE len=4 (conn) incr=9
2021/09/03 23:10:05 http2: Framer 0xc0001e0000: wrote PING len=8 ping="\x02\x04\x10\x10\t\x0e\a\a"
2021/09/03 23:10:05 sleep for 2s...
2021/09/03 23:10:05 http2: Framer 0xc0001e0000: read PING flags=ACK len=8 ping="\x02\x04\x10\x10\t\x0e\a\a"
2021/09/03 23:10:06 http2: Framer 0xc0001e0000: read RST_STREAM stream=5 len=4 ErrCode=CANCEL
2021/09/03 23:10:06 sleep canceled: context canceled
//
// client timeout & server cancel:ON
//
2021/09/03 23:10:08 http2: Framer 0xc0001e0000: read HEADERS flags=END_HEADERS stream=7 len=15
2021/09/03 23:10:08 http2: decoded hpack field header field ":method" = "POST"
2021/09/03 23:10:08 http2: decoded hpack field header field ":scheme" = "http"
2021/09/03 23:10:08 http2: decoded hpack field header field ":path" = "/helloworld.Greeter/Sleep"
2021/09/03 23:10:08 http2: decoded hpack field header field ":authority" = "localhost:8080"
2021/09/03 23:10:08 http2: decoded hpack field header field "content-type" = "application/grpc"
2021/09/03 23:10:08 http2: decoded hpack field header field "user-agent" = "grpc-go/1.40.0"
2021/09/03 23:10:08 http2: decoded hpack field header field "te" = "trailers"
2021/09/03 23:10:08 http2: decoded hpack field header field "grpc-timeout" = "999963u"
2021/09/03 23:10:08 http2: Framer 0xc0001e0000: read DATA flags=END_STREAM stream=7 len=9 data="\x00\x00\x00\x00\x04\b\x02\x10\x01"
2021/09/03 23:10:08 sleep for 2s...
2021/09/03 23:10:08 http2: Framer 0xc0001e0000: wrote WINDOW_UPDATE len=4 (conn) incr=9
2021/09/03 23:10:08 http2: Framer 0xc0001e0000: wrote PING len=8 ping="\x02\x04\x10\x10\t\x0e\a\a"
2021/09/03 23:10:08 http2: Framer 0xc0001e0000: read PING flags=ACK len=8 ping="\x02\x04\x10\x10\t\x0e\a\a"
2021/09/03 23:10:09 sleep canceled: context deadline exceeded
2021/09/03 23:10:09 http2: Framer 0xc0001e0000: read RST_STREAM stream=7 len=4 ErrCode=CANCEL
```

## References

- [gRPC and Deadlines | gRPC](https://grpc.io/blog/deadlines/)
