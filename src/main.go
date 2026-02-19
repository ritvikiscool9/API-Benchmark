package main

import (
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
	targetURL := "https://httpbin.org/get"
	totalRequests := 50

	fmt.Printf("Starting benchmark on %s with %d requests...\n", targetURL, totalRequests)

	// Define a channel for goroutines to send their results to
	resultsChan := make(chan RequestResult, totalRequests)

	var wg sync.WaitGroup

	for i := 0; i < totalRequests; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			startTime := time.Now()
			
			resp, err := http.Get(targetURL)
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

	average := totalDuration / time.Duration(totalRequests)
	fmt.Printf("\nTotal Time: %v | Average Request Time: %v\n", totalDuration, average)
}