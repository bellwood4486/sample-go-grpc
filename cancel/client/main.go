package main

import (
	"context"
	"log"
	"sync"
	"time"

	pb "github.com/bellwood4486/sample-go-grpc/cancel/helloworld"
	"google.golang.org/grpc"
)

const (
	address = "localhost:8080"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	doCancelSample(c)
	doTimeoutSample(c)
}

func doCancelSample(c pb.GreeterClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		_, err := c.Sleep(ctx, &pb.SleepRequest{TimeInSec: 2})
		if err != nil {
			log.Printf("could not sleep: %v", err)
		}
	}()

	time.Sleep(time.Second)
	cancel()
	wg.Wait()
}

func doTimeoutSample(c pb.GreeterClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := c.Sleep(ctx, &pb.SleepRequest{TimeInSec: 2})
	if err != nil {
		log.Printf("could not sleep: %v", err)
	}
}
