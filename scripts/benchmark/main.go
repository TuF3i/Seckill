package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math"
	"net/http"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

type Response struct {
	Status int32       `json:"status"`
	Info   string      `json:"info"`
	Data   interface{} `json:"data"`
}

func main() {
	gateway := flag.String("gateway", "http://localhost:8888", "Gateway address")
	itemId := flag.String("item", "", "Item ID to order (required)")
	token := flag.String("token", "", "JWT access token (required)")
	concurrency := flag.Int("c", 50, "Concurrency")
	total := flag.Int("n", 2000, "Total requests")
	flag.Parse()

	if *itemId == "" || *token == "" {
		fmt.Println("Usage: go run main.go --item <itemId> --token <accessToken>")
		fmt.Println("  -gateway  string  Gateway address (default: http://localhost:8888)")
		fmt.Println("  -item     string  Item ID to order (required)")
		fmt.Println("  -token    string  JWT access token (required)")
		fmt.Println("  -c        int     Concurrency (default: 50)")
		fmt.Println("  -n        int     Total requests (default: 2000)")
		return
	}

	fmt.Printf("Seckill Benchmark\n")
	fmt.Printf("  Gateway:     %s\n", *gateway)
	fmt.Printf("  Item ID:     %s\n", *itemId)
	fmt.Printf("  Token:       %s...\n", (*token)[:min(16, len(*token))])
	fmt.Printf("  Concurrency: %d\n", *concurrency)
	fmt.Printf("  Requests:    %d\n\n", *total)

	runBenchmark(*gateway, *token, *itemId, *concurrency, *total)
}

func runBenchmark(gateway, token, itemId string, concurrency, total int) {
	body := fmt.Sprintf(`{"itemId":"%s"}`, itemId)

	var wg sync.WaitGroup
	var success, failed int64
	durations := make([]float64, total)
	errMap := sync.Map{}
	idxMu := sync.Mutex{}
	idx := 0

	start := time.Now()

	sem := make(chan struct{}, concurrency)

	fmt.Println("=== Benchmark Running ===")

	ticker := time.NewTicker(2 * time.Second)
	done := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				elapsed := time.Since(start).Seconds()
				doneCount := int(atomic.LoadInt64(&success) + atomic.LoadInt64(&failed))
				if doneCount > 0 {
					currentRps := float64(doneCount) / elapsed
					fmt.Printf("  Progress: %d/%d (%.1f%%) | RPS: %.0f\n",
						doneCount, total, float64(doneCount)/float64(total)*100, currentRps)
				}
			case <-done:
				ticker.Stop()
				return
			}
		}
	}()

	for i := 0; i < total; i++ {
		wg.Add(1)
		sem <- struct{}{}

		go func() {
			defer wg.Done()
			defer func() { <-sem }()
			reqStart := time.Now()

			req, _ := http.NewRequest("POST", gateway+"/order/create",
				bytes.NewReader([]byte(body)))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)

			resp, err := http.DefaultClient.Do(req)
			cost := time.Since(reqStart).Seconds() * 1000

			idxMu.Lock()
			pos := idx
			idx++
			durations[pos] = cost
			idxMu.Unlock()

			if err != nil {
				atomic.AddInt64(&failed, 1)
				errMap.Store(pos, err.Error())
				return
			}
			defer resp.Body.Close()

			respBody, _ := io.ReadAll(resp.Body)
			var result Response
			json.Unmarshal(respBody, &result)

			if result.Status == 20000 {
				atomic.AddInt64(&success, 1)
			} else {
				atomic.AddInt64(&failed, 1)
				errMap.Store(pos, fmt.Sprintf("status=%d info=%s", result.Status, result.Info))
			}
		}()
	}

	wg.Wait()
	close(done)
	elapsed := time.Since(start)

	durMs := durations[:atomic.LoadInt64(&success)+atomic.LoadInt64(&failed)]
	sort.Float64s(durMs)
	n := len(durMs)

	var sum float64
	for _, d := range durMs {
		sum += d
	}
	avg := sum / float64(n)

	var p50, p90, p99, p999, minVal, maxVal float64
	if n > 0 {
		minVal = durMs[0]
		maxVal = durMs[n-1]
		p50 = percentile(durMs, 50)
		p90 = percentile(durMs, 90)
		p99 = percentile(durMs, 99)
		p999 = percentile(durMs, 99.9)
	}

	errCounts := map[string]int64{}
	errMap.Range(func(_, value interface{}) bool {
		errCounts[value.(string)]++
		return true
	})

	totalDone := success + failed

	fmt.Printf("\n=== Summary ===\n")
	fmt.Printf("Total:          %d\n", total)
	fmt.Printf("Completed:      %d (%.1f%%)\n", totalDone, float64(totalDone)/float64(total)*100)
	fmt.Printf("Success:        %d (%.1f%%)\n", success, float64(success)/float64(total)*100)
	fmt.Printf("Failed:         %d (%.1f%%)\n", failed, float64(failed)/float64(total)*100)

	fmt.Printf("\n=== Throughput ===\n")
	fmt.Printf("QPS (Queries Per Second):       %.2f\n", float64(totalDone)/elapsed.Seconds())
	fmt.Printf("TPS (Transactions Per Second):  %.2f\n", float64(success)/elapsed.Seconds())
	fmt.Printf("Elapsed:                        %v\n", elapsed)

	fmt.Printf("\n=== Latency (ms) ===\n")
	fmt.Printf("Min:    %.2f\n", minVal)
	fmt.Printf("Avg:    %.2f\n", avg)
	fmt.Printf("Max:    %.2f\n", maxVal)
	fmt.Printf("P50:    %.2f\n", p50)
	fmt.Printf("P90:    %.2f\n", p90)
	fmt.Printf("P99:    %.2f\n", p99)
	fmt.Printf("P99.9:  %.2f\n", p999)

	if len(errCounts) > 0 {
		fmt.Printf("\n=== Error Breakdown ===\n")
		for code, count := range errCounts {
			fmt.Printf("  %s: %d\n", code, count)
		}
	}
	fmt.Println()
}

func percentile(data []float64, p float64) float64 {
	if len(data) == 0 {
		return 0
	}
	rank := p / 100.0 * float64(len(data)-1)
	lower := int(math.Floor(rank))
	upper := int(math.Ceil(rank))
	if lower == upper {
		return data[lower]
	}
	return data[lower] + (rank-float64(lower))*(data[upper]-data[lower])
}
