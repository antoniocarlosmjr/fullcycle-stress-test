package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

var (
	url         string
	requests    int
	concurrency int
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "stress_test",
		Short: "Stress testing tool",
		Run: func(cmd *cobra.Command, args []string) {
			if concurrency > requests {
				fmt.Println("Error: Concurrency cannot be greater than the number of requests.")
				os.Exit(1)
			}
			runLoadTest(url, requests, concurrency)
		},
	}

	rootCmd.Flags().StringVar(&url, "url", "", "URL of the service to test")
	rootCmd.Flags().IntVar(&requests, "requests", 1, "Total number of requests")
	rootCmd.Flags().IntVar(&concurrency, "concurrency", 1, "Number of concurrent requests")

	if err := rootCmd.MarkFlagRequired("url"); err != nil {
		log.Fatalf("url flag is required: %v", err)
	}

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("cmd.Execute() failed: %v", err)
	}
}

func runLoadTest(url string, requests int, concurrency int) {
	var wg sync.WaitGroup
	var mu sync.Mutex

	startTime := time.Now()

	statusCodes := make(map[int]int)
	totalRequests := 0

	sem := make(chan struct{}, concurrency)
	done := make(chan struct{})

	go func() {
		fmt.Println("Making requests, please wait....")
	}()

	for i := 0; i < requests; i++ {
		wg.Add(1)
		sem <- struct{}{}

		go func() {
			defer wg.Done()
			defer func() { <-sem }()

			resp, err := http.Get(url)
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				if os.IsTimeout(err) {
					statusCodes[408]++ // Request Timeout HTTP status code
				} else {
					statusCodes[500]++ // Assume 500 for any other error
				}
				totalRequests++
				return
			}
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					log.Printf("Error closing response body: %v", err)
				}
			}(resp.Body)

			statusCodes[resp.StatusCode]++
			totalRequests++
		}()
	}

	// Wait for all requests to complete
	go func() {
		wg.Wait()
		close(done)
	}()

	<-done // Wait until all requests are done

	totalTime := time.Since(startTime)
	generateReport(totalRequests, statusCodes, totalTime)
}

func generateReport(totalRequests int, statusCodes map[int]int, totalTime time.Duration) {
	fmt.Println("Report: ")
	fmt.Printf("Total requests: %d\n", totalRequests)
	fmt.Printf("Time taken: %s\n", totalTime)
	fmt.Printf("Status code distribution:\n")

	for code, count := range statusCodes {
		fmt.Printf("[%d] %d requests \n", code, count)
	}
}
