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

	doCancelSample(c, false)
	time.Sleep(2 * time.Second)
	doTimeoutSample(c, false)
	time.Sleep(2 * time.Second)

	doCancelSample(c, true)
	time.Sleep(2 * time.Second)
	doTimeoutSample(c, true)
}

func doCancelSample(c pb.GreeterClient, cancelOnServer bool) {
	log.Println("start cancel sample")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		_, err := c.Sleep(ctx, &pb.SleepRequest{TimeInSec: 2, WantCancel: cancelOnServer})
		if err != nil {
			log.Printf("could not sleep: %v", err)
		}
	}()

	// wait for 1 second and cancel
	time.Sleep(time.Second)
	cancel()
	wg.Wait()
}

func doTimeoutSample(c pb.GreeterClient, cancelOnServer bool) {
	log.Println("start timeout sample")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := c.Sleep(ctx, &pb.SleepRequest{TimeInSec: 2, WantCancel: cancelOnServer})
	if err != nil {
		log.Printf("could not sleep: %v", err)
	}
}
