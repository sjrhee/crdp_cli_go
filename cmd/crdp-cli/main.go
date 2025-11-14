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
	showBody := flag.Bool("show-body", false, "show request/response URLs and JSON bodies")
	useBulk := flag.Bool("bulk", false, "use bulk protect/reveal endpoints")
	batchSize := flag.Int("batch-size", 50, "batch size for bulk operations")
	useTLS := flag.Bool("tls", false, "use HTTPS instead of HTTP")

	// 커스텀 Usage 함수로 -- 형식 표시
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  --host string\n        API host (default \"192.168.0.231\")\n")
		fmt.Fprintf(os.Stderr, "  --port int\n        API port (default 32082)\n")
		fmt.Fprintf(os.Stderr, "  --policy string\n        protection_policy_name (default \"P03\")\n")
		fmt.Fprintf(os.Stderr, "  --iterations int\n        number of iterations (default 100)\n")
		fmt.Fprintf(os.Stderr, "  --timeout int\n        per-request timeout seconds (default 10)\n")
		fmt.Fprintf(os.Stderr, "  --verbose\n        enable debug logging\n")
		fmt.Fprintf(os.Stderr, "  --show-progress\n        show per-iteration progress output\n")
		fmt.Fprintf(os.Stderr, "  --show-body\n        show request/response URLs and JSON bodies\n")
		fmt.Fprintf(os.Stderr, "  --bulk\n        use bulk protect/reveal endpoints\n")
		fmt.Fprintf(os.Stderr, "  --batch-size int\n        batch size for bulk operations (default 50)\n")
		fmt.Fprintf(os.Stderr, "  --tls\n        use HTTPS instead of HTTP\n")
	}

	flag.Parse()

	// 클라이언트 생성
	c := client.NewClient(*host, *port, *policy, *timeout, *useTLS)
	c.SetShowBody(*showBody)

	// 반복 실행
	startTime := time.Now()
	var results []*runner.IterationResult
	successfulItems := 0
	matchedItems := 0
	totalItems := 0
	sumBatchTimes := 0.0

	if *useBulk {
		// Bulk 모드: 입력 데이터 생성
		inputs := make([]string, *iterations)
		for i := 0; i < *iterations; i++ {
			inputs[i] = fmt.Sprintf("1234567890123%d", i)
		}

		// 배치 단위로 처리
		for i := 0; i < len(inputs); i += *batchSize {
			end := i + *batchSize
			if end > len(inputs) {
				end = len(inputs)
			}
			batch := inputs[i:end]

			result, err := runner.RunBulkIteration(c, batch)
			if err != nil {
				if *verbose {
					log.Printf("Error at batch %d: %v", i/ *batchSize+1, err)
				}
				continue
			}

			results = append(results, result)
			totalItems += len(batch)
			sumBatchTimes += result.TimeS

			// 성공 및 매칭된 항목 카운트
			if result.ProtectResponse.StatusCode >= 200 && result.ProtectResponse.StatusCode < 300 &&
				result.RevealResponse.StatusCode >= 200 && result.RevealResponse.StatusCode < 300 {
				if result.RestoredCount == len(batch) {
					successfulItems += len(batch)
				} else {
					successfulItems += result.RestoredCount
				}
				}
			matchedItems += result.MatchedCount

			if *showProgress {
				fmt.Fprintf(os.Stderr, "Batch #%03d size=%d time=%.4fs protect_status=%d reveal_status=%d matched=%d\n",
					i/ *batchSize+1, len(batch), result.TimeS,
					result.ProtectResponse.StatusCode, result.RevealResponse.StatusCode, result.MatchedCount)
			}
		}
	} else {
		// 일반 모드 (단일 처리)
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
				successfulItems++
			}
			if result.Match {
				matchedItems++
			}

			if *showProgress {
				fmt.Fprintf(os.Stderr, "#%03d data=%s time=%.4fs protect_status=%d reveal_status=%d match=%v\n",
					i, data, result.TimeS, result.ProtectResponse.StatusCode, result.RevealResponse.StatusCode, result.Match)
			}
		}
	}

	total := time.Since(startTime)

	// 결과 출력
	fmt.Printf("\nSummary:\n")
	if *useBulk {
		fmt.Printf("Iterations attempted: %d\n", totalItems)
		fmt.Printf("Successful (both 2xx): %d\n", successfulItems)
		fmt.Printf("Revealed matched original data: %d\n", matchedItems)
		fmt.Printf("Total time: %.4fs\n", total.Seconds())
		if totalItems > 0 {
			avgTime := sumBatchTimes / float64(totalItems)
			fmt.Printf("Average per-iteration time: %.4fs\n", avgTime)
		}
	} else {
		fmt.Printf("Iterations attempted: %d\n", len(results))
		fmt.Printf("Successful (both 2xx): %d\n", successfulItems)
		fmt.Printf("Revealed matched original data: %d\n", matchedItems)
		fmt.Printf("Total time: %.4fs\n", total.Seconds())
		if len(results) > 0 {
			avgTime := total.Seconds() / float64(len(results))
			fmt.Printf("Average per-iteration time: %.4fs\n", avgTime)
		}
	}
}
