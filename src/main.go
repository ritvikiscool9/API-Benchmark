package main

import (
	"flag"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// Struct to hod the result of each request
type RequestResult struct{
	Status string
	Duration time.Duration
}

func main() {
	urlPtr := flag.String("url", "https://httpbin.org/get", "URL to test")
	totalRequestsPtr := flag.Int("requests", 10, "number of requests to make")

	flag.Parse()

	if *totalRequestsPtr <= 0 {
		fmt.Printf("Request amount must be greater than 0, defaulting to 10\n")
		*totalRequestsPtr = 10
	}

	fmt.Printf("Starting benchmark on %s with %d requests...\n", *urlPtr, *totalRequestsPtr)

	// Define a channel for goroutines to send their results to
	resultsChan := make(chan RequestResult, *totalRequestsPtr)

	var wg sync.WaitGroup

	for i := 0; i < *totalRequestsPtr; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			startTime := time.Now()
			
			resp, err := http.Get(*urlPtr)
			if err != nil {
				// For now, we'll just print errors, but ideally these go to an error channel too
				fmt.Printf("Error: %s\n", err)
				return
			}
			defer resp.Body.Close()

			duration := time.Since(startTime)

			// Send results to the channel
			resultsChan <- RequestResult{
				Status:   resp.Status,
				Duration: duration,
			}
		}()
	}

	// Wait for all workers to finish
	wg.Wait()

	close(resultsChan)

	var totalDuration time.Duration
	
	// Loop over the channel 
	for result := range resultsChan {
		totalDuration += result.Duration
	}

	average := totalDuration / time.Duration(*totalRequestsPtr)
	fmt.Printf("\nTotal Time: %v | Average Request Time: %v\n", totalDuration, average)
}