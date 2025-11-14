package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/sjrhee/crdp-cli-go/internal/client"
	"github.com/sjrhee/crdp-cli-go/internal/runner"
)

func main() {
	// CLI 플래그 정의
	host := flag.String("host", "192.168.0.231", "API host")
	port := flag.Int("port", 32082, "API port")
	policy := flag.String("policy", "P03", "protection_policy_name")
	iterations := flag.Int("iterations", 100, "number of iterations")
	timeout := flag.Int("timeout", 10, "per-request timeout seconds")
	verbose := flag.Bool("verbose", false, "enable debug logging")
	showProgress := flag.Bool("show-progress", false, "show per-iteration progress output")
	useTLS := flag.Bool("tls", false, "use HTTPS instead of HTTP")

	flag.Parse()

	// 클라이언트 생성
	c := client.NewClient(*host, *port, *policy, *timeout, *useTLS)

	// 반복 실행
	startTime := time.Now()
	var results []*runner.IterationResult
	successful := 0
	matched := 0

	for i := 1; i <= *iterations; i++ {
		data := fmt.Sprintf("1234567890123%d", i-1)
		result, err := runner.RunIteration(c, data)
		if err != nil {
			if *verbose {
				log.Printf("Error at iteration %d: %v", i, err)
			}
			continue
		}

		results = append(results, result)

		if result.Success {
			successful++
		}
		if result.Match {
			matched++
		}

		if *showProgress {
			fmt.Fprintf(os.Stderr, "#%03d data=%s time=%.4fs protect_status=%d reveal_status=%d match=%v\n",
				i, data, result.TimeS, result.ProtectResponse.StatusCode, result.RevealResponse.StatusCode, result.Match)
		}
	}

	total := time.Since(startTime)

	// 결과 출력
	fmt.Printf("\nSummary:\n")
	fmt.Printf("Iterations attempted: %d\n", len(results))
	fmt.Printf("Successful (both 2xx): %d\n", successful)
	fmt.Printf("Revealed matched original data: %d\n", matched)
	fmt.Printf("Total time: %.4fs\n", total.Seconds())
	if len(results) > 0 {
		avgTime := total.Seconds() / float64(len(results))
		fmt.Printf("Average per-iteration time: %.4fs\n", avgTime)
	}
}
