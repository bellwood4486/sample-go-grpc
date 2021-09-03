# Cancel Sample

## How to run

```shell
go run server/main.go
```

```shell
go run client/main.go
```

## Results

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
