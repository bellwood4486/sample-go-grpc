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
	port = ":8080"
)

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())

	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func (s *server) Sleep(ctx context.Context, in *pb.SleepRequest) (*pb.SleepReply, error) {
	log.Printf("Sleep %d sec...", in.TimeInSec)
	time.Sleep(time.Duration(in.TimeInSec) * time.Second)

	return &pb.SleepReply{}, nil
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
