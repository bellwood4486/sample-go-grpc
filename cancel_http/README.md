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
