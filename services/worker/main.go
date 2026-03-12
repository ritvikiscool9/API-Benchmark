package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"APIbenchmark/proto/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Flags to take in user input for target URL and amount of requests 
	urlPtr := flag.String("url", "https://httpbin.org/get", "URL to test")
	totalRequestsPtr := flag.Int("requests", 10, "number of requests to make")

	flag.Parse()

	if *totalRequestsPtr <= 0 {
		fmt.Printf("Request amount must be greater than 0, defaulting to 10\n")
		*totalRequestsPtr = 10
	}

	fmt.Printf("Starting benchmark on %s with %d requests...\n", *urlPtr, *totalRequestsPtr)

	// Create grpc server
	conn, err := grpc.NewClient("localhost:50051",grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Server creation failed")
	}
	defer conn.Close()

	// Create the gRPC client used to submit benchmark results
	client := pb.NewAggregatorClient(conn)

	// WaitGroup keeps track of all request goroutines
	var wg sync.WaitGroup

	for i := 0; i < *totalRequestsPtr; i++ {
		// Register one pending goroutine so Wait can block until all requests finish
		wg.Add(1)

		go func() {
			// Ensure this goroutine always signals completion, even on early return
			defer wg.Done()
			// Record when this request starts so we can calculate end-to-end latency
			startTime := time.Now()
			
			// Send the HTTP request to the target URL
			resp, err := http.Get(*urlPtr)
			
			if err != nil {
				fmt.Printf("Error: %s\n", err)
				return
			}
			defer resp.Body.Close()

			// Convert request duration to milliseconds for reporting
			duration := time.Since(startTime)
			latencyMs := duration.Milliseconds()

			// Package this request's metrics to send to the aggregator service
			res := &pb.Result{
				Time: startTime.Unix(),
				Latency: latencyMs,
				Status: int32(resp.StatusCode),
			}

			// Submit the result over gRPC.
			client.SubmitResults(context.Background(), res)

		}()
	}

	// Wait for all workers to finish
	wg.Wait()
}