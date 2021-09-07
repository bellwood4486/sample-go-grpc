package main

import (
	"context"
	"log"
	"net"
	"time"

	pb "github.com/bellwood4486/sample-go-grpc/cancel/helloworld"
	"google.golang.org/grpc"
)

const (
	port = ":18080"
)

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())

	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func (s *server) Sleep(ctx context.Context, in *pb.SleepRequest) (*pb.SleepReply, error) {
	d := time.Duration(in.TimeInSec) * time.Second
	log.Printf("sleep for %v...", d)

	if in.WantCancel {
		if err := s.sleepCancelable(ctx, d); err != nil {
			log.Printf("could not sleep: %v", err)
			return nil, err
		}
	} else {
		time.Sleep(d)
	}

	log.Println("sleep success")
	return &pb.SleepReply{}, nil
}

func (s *server) sleepCancelable(ctx context.Context, duration time.Duration) error {
	done := make(chan struct{})

	go func() {
		defer func() { done <- struct{}{} }()
		time.Sleep(duration)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		if err := ctx.Err(); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
