package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

var (
	Reset = "\033[0m"
	Red   = "\033[31m"
	Green = "\033[32m"
)

type result struct {
	url        string
	status     string
	statusCode int
	duration   time.Duration
	err        error
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <url1> <url2> <url3>...")
		os.Exit(1)
	}

	urls := os.Args[1:]

	results := make(chan result)

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	for _, u := range urls {

		go func(url string) {
			start := time.Now()

			resp, err := client.Get(url)
			duration := time.Since(start)

			if err != nil {
				results <- result{
					url:      url,
					duration: duration,
					err:      err,
				}
				return
			}
			defer resp.Body.Close()

			results <- result{
				url:        url,
				duration:   duration,
				status:     resp.Status,
				statusCode: resp.StatusCode,
				err:        nil,
			}
		}(u)
	}

	for range urls {
		r := <-results
		if r.err != nil {
			fmt.Printf("%-30s %sERROR%s: %v (after %s)\n", r.url, Red, Reset, r.err, r.duration)
		} else {
			var statusColor string
			if r.statusCode == 200 {
				statusColor = Green
			} else {
				statusColor = Red
			}
			fmt.Printf("%-30s %s%s%s (%s)\n", r.url, statusColor, r.status, Reset, r.duration)
		}
	}
}
