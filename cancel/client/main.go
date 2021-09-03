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
	address       = "localhost:8080"
	wantCancel    = true
	notWantCancel = false
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	doCancelSample(c, notWantCancel)
	time.Sleep(2 * time.Second)
	doTimeoutSample(c, notWantCancel)
	time.Sleep(2 * time.Second)

	doCancelSample(c, wantCancel)
	time.Sleep(2 * time.Second)
	doTimeoutSample(c, wantCancel)
}

func doCancelSample(c pb.GreeterClient, wantCancel bool) {
	log.Println("start cancel sample")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		_, err := c.Sleep(ctx, &pb.SleepRequest{TimeInSec: 2, WantCancel: wantCancel})
		if err != nil {
			log.Printf("could not sleep: %v", err)
		}
	}()

	time.Sleep(time.Second)
	cancel()
	wg.Wait()
}

func doTimeoutSample(c pb.GreeterClient, wantCancel bool) {
	log.Println("start timeout sample")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := c.Sleep(ctx, &pb.SleepRequest{TimeInSec: 2, WantCancel: wantCancel})
	if err != nil {
		log.Printf("could not sleep: %v", err)
	}
}
