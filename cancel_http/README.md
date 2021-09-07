# Sample for canceling HTTP connection

## Flow

```
┌─────────────┐            ┌─────────────┐             ┌────────────┐
│cURL         │            │HTTP server  │             │gRPC server │
│             ├────────────►             ├─────────────►            │
│             │GET /sleep  │             │call Sleep   │            │
│             │    &       │             │             │            │
└─────────────┘timeout     └─────────────┘             └────────────┘
                after 2sec
```

## How to run

gRPC server
```shell
go run server_grpc/main.go
```

HTTP server
```shell
go run server_http/main.go
```

Client
```shell
curl -v -m 2 http://127.0.0.1:18081/sleep
```
`-m 2` means timeout after 2 seconds.

## Results

Client
```
❯ curl -v -m 2 http://127.0.0.1:18081/sleep
*   Trying 127.0.0.1...
* TCP_NODELAY set
* Connected to 127.0.0.1 (127.0.0.1) port 18081 (#0)
> GET /sleep HTTP/1.1
> Host: 127.0.0.1:18081
> User-Agent: curl/7.64.1
> Accept: */*
> 
* Operation timed out after 2001 milliseconds with 0 bytes received
* Closing connection 0
curl: (28) Operation timed out after 2001 milliseconds with 0 bytes received
```

HTTP server
```
❯ go run server_http/main.go
2021/09/04 23:25:19 Received: /sleep
2021/09/04 23:25:19 sleep for 5 seconds...
2021/09/04 23:25:21 could not sleep: rpc error: code = Canceled desc = context canceled
2021/09/04 23:25:21 canceled by client
```

gRPC server
```
❯ go run server_grpc/main.go
2021/09/04 23:25:19 sleep for 5s...
2021/09/04 23:25:21 sleep canceled: context canceled
```

## Code

### http server side

The http server calls cancel() in `cr.handleReadError(err)` when the connection closes.

https://github.com/golang/go/blob/bc51e930274a5d5835ac8797978afc0864c9e30c/src/net/http/server.go#L703
```go
func (cr *connReader) backgroundRead() {
	n, err := cr.conn.rwc.Read(cr.byteBuf[:])
	cr.lock()
	if n == 1 {
		cr.hasByte = true
		// We were past the end of the previous request's body already
		// (since we wouldn't be in a background read otherwise), so
		// this is a pipelined HTTP request. Prior to Go 1.11 we used to
		// send on the CloseNotify channel and cancel the context here,
		// but the behavior was documented as only "may", and we only
		// did that because that's how CloseNotify accidentally behaved
		// in very early Go releases prior to context support. Once we
		// added context support, people used a Handler's
		// Request.Context() and passed it along. Having that context
		// cancel on pipelined HTTP requests caused problems.
		// Fortunately, almost nothing uses HTTP/1.x pipelining.
		// Unfortunately, apt-get does, or sometimes does.
		// New Go 1.11 behavior: don't fire CloseNotify or cancel
		// contexts on pipelined requests. Shouldn't affect people, but
		// fixes cases like Issue 23921. This does mean that a client
		// closing their TCP connection after sending a pipelined
		// request won't cancel the context, but we'll catch that on any
		// write failure (in checkConnErrorWriter.Write).
		// If the server never writes, yes, there are still contrived
		// server & client behaviors where this fails to ever cancel the
		// context, but that's kinda why HTTP/1.x pipelining died
		// anyway.
	}
	if ne, ok := err.(net.Error); ok && cr.aborted && ne.Timeout() {
		// Ignore this error. It's the expected error from
		// another goroutine calling abortPendingRead.
	} else if err != nil {
		cr.handleReadError(err)
	}
	cr.aborted = false
	cr.inRead = false
	cr.unlock()
	cr.cond.Broadcast()
}
```

### gRPC client side

In gRPC client side, a message is received from `<-s.ctx.Done()` in `waitOnHeader()`. Then `ContextErr()` creates an error.

https://github.com/grpc/grpc-go/blob/41e044e1c82fcf6a5801d6cbd7ecf952505eecb1/internal/transport/transport.go#L326
```go
func (s *Stream) waitOnHeader() {
	if s.headerChan == nil {
		// On the server headerChan is always nil since a stream originates
		// only after having received headers.
		return
	}
	select {
	case <-s.ctx.Done():
		// Close the stream to prevent headers/trailers from changing after
		// this function returns.
		s.ct.CloseStream(s, ContextErr(s.ctx.Err()))
		// headerChan could possibly not be closed yet if closeStream raced
		// with operateHeaders; wait until it is closed explicitly here.
		<-s.headerChan
	case <-s.headerChan:
	}
}
```

### gRPC server side

Parse `grpc-timeout` header.
Set the context with cancel to `transport.Stream`.

https://github.com/grpc/grpc-go/blob/41e044e1c82fcf6a5801d6cbd7ecf952505eecb1/internal/transport/http2_server.go#L407
```go
// operateHeader takes action on the decoded headers.
func (t *http2Server) operateHeaders(frame *http2.MetaHeadersFrame, handle func(*Stream), traceCtx func(context.Context, string) context.Context) (fatal bool) {
...
	for _, hf := range frame.Fields {
		switch hf.Name {
		...
		case "grpc-timeout":
			timeoutSet = true
			var err error
			if timeout, err = decodeTimeout(hf.Value); err != nil {
				headerError = true
			}
			...
	}

	...
	if timeoutSet {
		s.ctx, s.cancel = context.WithTimeout(t.ctx, timeout)
	} else {
		s.ctx, s.cancel = context.WithCancel(t.ctx)
	}
```

https://github.com/grpc/grpc-go/blob/41e044e1c82fcf6a5801d6cbd7ecf952505eecb1/internal/transport/transport.go#L242
```go
// Stream represents an RPC in the transport layer.
type Stream struct {
	...
	cancel       context.CancelFunc // always nil for client side Stream
	...
```

When gPRC server receives RST Frame, it calls `(*http2Server).handleRSTStream` method.
The method calls `(*http2Server).closeStream`.

https://github.com/grpc/grpc-go/blob/41e044e1c82fcf6a5801d6cbd7ecf952505eecb1/internal/transport/http2_server.go#L581-L582
```go
// HandleStreams receives incoming streams using the given handler. This is
// typically run in a separate goroutine.
// traceCtx attaches trace to ctx and returns the new context.
func (t *http2Server) HandleStreams(handle func(*Stream), traceCtx func(context.Context, string) context.Context) {
		...
		switch frame := frame.(type) {
		...
		case *http2.RSTStreamFrame:
			t.handleRSTStream(frame)
```


`(*http2Server).closeStream` calls `(*http2Server).deleteStream`.

https://github.com/grpc/grpc-go/blob/41e044e1c82fcf6a5801d6cbd7ecf952505eecb1/internal/transport/http2_server.go#L1209
```go
// closeStream clears the footprint of a stream when the stream is not needed any more.
func (t *http2Server) closeStream(s *Stream, rst bool, rstCode http2.ErrCode, eosReceived bool) {
	s.swapState(streamDone)
	t.deleteStream(s, eosReceived)

	...
}
```

`(*http2Server).deleteStream` calls `(*Stream).cancel` set above.

https://github.com/grpc/grpc-go/blob/41e044e1c82fcf6a5801d6cbd7ecf952505eecb1/internal/transport/http2_server.go#L1167
```go
// deleteStream deletes the stream s from transport's active streams.
func (t *http2Server) deleteStream(s *Stream, eosReceived bool) {
	// In case stream sending and receiving are invoked in separate
	// goroutines (e.g., bi-directional streaming), cancel needs to be
	// called to interrupt the potential blocking on other goroutines.
	s.cancel()

	...
```

