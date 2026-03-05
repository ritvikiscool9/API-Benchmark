package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	// "context"
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

	client := pb.NewAggregatorClient(conn)

	var wg sync.WaitGroup

	for i := 0; i < *totalRequestsPtr; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			startTime := time.Now()
			
			resp, err := http.Get(*urlPtr)

			
			if err != nil {
				fmt.Printf("Error: %s\n", err)
				return
			}
			defer resp.Body.Close()

			duration := time.Since(startTime)

		}()
	}

	// Wait for all workers to finish
	wg.Wait()
}