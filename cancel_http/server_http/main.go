package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	pb "github.com/bellwood4486/sample-go-grpc/cancel/helloworld"
	"google.golang.org/grpc"
)

const (
	grpcAddress = "localhost:18080"
	httpAddress = ":18081"
)

func main() {
	http.HandleFunc("/sleep", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Received: /sleep")
		ctx := r.Context()
		fmt.Println("sleep for 5 seconds...")
		if err := invokeSleep(ctx, 5); err != nil {
			log.Printf("could not sleep: %v", err)
			status := http.StatusInternalServerError
			if errors.Is(ctx.Err(), context.Canceled) {
				log.Printf("canceled:  %v", err)
				status = http.StatusBadRequest
			}
			http.Error(w, fmt.Sprint(err), status)
			return
		}
		fmt.Println("sleep success")

		_, _ = fmt.Fprintf(w, "slept")
	})

	log.Fatal(http.ListenAndServe(httpAddress, nil))
}

func invokeSleep(ctx context.Context, timeInSec int32) error {
	conn, err := grpc.Dial(grpcAddress, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return err
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	_, err = c.Sleep(ctx, &pb.SleepRequest{TimeInSec: timeInSec, WantCancel: false})
	if err != nil {
		return err
	}

	return nil
}

func doCancelSample(c pb.GreeterClient, cancelOnServer bool) {
	log.Println("start cancel sample")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// wait for 1 second and cancel
	go func() {
		time.Sleep(time.Second)
		cancel()
	}()

	_, err := c.Sleep(ctx, &pb.SleepRequest{TimeInSec: 2, WantCancel: cancelOnServer})
	if err != nil {
		log.Printf("could not sleep: %v", err)
	}
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
