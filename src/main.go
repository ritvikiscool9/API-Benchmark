package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

func main(){
	targetURL := "https://httpbin.org/get"
	totalRequests := 50

	var wg sync.WaitGroup

	for i := 0; i < totalRequests; i++ {
		wg.Add(1)
		go func(){
			defer wg.Done()
			startTime := time.Now()
			resp, err := http.Get(targetURL)
			if err != nil {
				fmt.Printf("Request failed: %s\n", err)
				return // exit this specific goroutine if it failed
			}
			defer resp.Body.Close()
			timeSince :=time.Since(startTime)
			fmt.Printf("Status: %s | Time: %v\n", resp.Status, timeSince)
		}()
	}
	
	wg.Wait()
	fmt.Println("Benchmark complete!")

}