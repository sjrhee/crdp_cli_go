package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/sjrhee/crdp-cli-go/internal/client"
	"github.com/sjrhee/crdp-cli-go/internal/config"
	"github.com/sjrhee/crdp-cli-go/internal/runner"
)

// incrementNumericString increments a numeric string by 1
func incrementNumericString(s string) (string, error) {
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return "", fmt.Errorf("not a valid numeric string: %w", err)
	}
	return fmt.Sprintf("%d", n+1), nil
}

// generateDataSequence generates a sequence of data strings starting from startData
func generateDataSequence(startData string, count int, verbose bool) []string {
	inputs := make([]string, 0, count)
	currentData := startData
	for i := 0; i < count; i++ {
		inputs = append(inputs, currentData)
		if i < count-1 {
			nextData, err := incrementNumericString(currentData)
			if err != nil {
				if verbose {
					log.Printf("Cannot increment data '%s': %v", currentData, err)
				}
				break
			}
			currentData = nextData
		}
	}
	return inputs
}

// printSummary prints execution summary
// printSummary prints execution summary
func printSummary(attempted, successful, matched int, totalTime time.Duration) {
	fmt.Printf("\nSummary\n")
	fmt.Printf("- Iterations attempted: %d\n", attempted)
	fmt.Printf("- Successful (both 2xx): %d\n", successful)
	fmt.Printf("- Revealed matched original data: %d\n", matched)
	fmt.Printf("- Total time: %.4fs\n", totalTime.Seconds())
	if attempted > 0 {
		avgTime := totalTime.Seconds() / float64(attempted)
		fmt.Printf("- Average per-iteration time: %.4fs\n", avgTime)
	}
}

// printBulkProgress prints progress for bulk mode
func printBulkProgress(batchNum, batchSize int, timeS float64, protectStatus, revealStatus, matched int) {
	fmt.Fprintf(os.Stderr, "Batch #%03d size=%d time=%.4fs protect_status=%d reveal_status=%d matched=%d\n",
		batchNum, batchSize, timeS, protectStatus, revealStatus, matched)
}

// printIterationProgress prints progress for single iteration mode
func printIterationProgress(iterNum int, data string, timeS float64, protectStatus, revealStatus int, match bool, withBlankLine bool) {
	if withBlankLine {
		fmt.Fprintf(os.Stderr, "#%03d data=%s time=%.4fs protect_status=%d reveal_status=%d match=%v\n\n",
			iterNum, data, timeS, protectStatus, revealStatus, match)
	} else {
		fmt.Fprintf(os.Stderr, "#%03d data=%s time=%.4fs protect_status=%d reveal_status=%d match=%v\n",
			iterNum, data, timeS, protectStatus, revealStatus, match)
	}
}

func main() {
	// 설정 파일 경로 플래그 (먼저 정의)
	configPath := flag.String("config", "", "path to config.yaml file (default: auto-search)")

	// 플래그 파싱 전에 config 플래그만 처리하기 위해 특수 처리
	flag.Parse()

	// 설정 파일 로드
	var cfg *config.Config
	var err error
	if *configPath != "" {
		// 명시적으로 지정된 config 파일 사용
		cfg, err = config.LoadConfig(*configPath)
	} else {
		// 자동 검색
		cfg, err = config.LoadConfig(config.GetConfigPath())
	}
	if err != nil {
		log.Printf("Warning: failed to load config file: %v. Using defaults.\n", err)
		cfg = config.DefaultConfig()
	}

	// CLI 플래그 다시 정의 (config.yaml에서 읽은 기본값 사용)
	host := flag.String("host", cfg.API.Host, "API host")
	port := flag.Int("port", cfg.API.Port, "API port")
	policy := flag.String("policy", cfg.Protection.Policy, "protection_policy_name")
	startData := flag.String("start-data", cfg.Execution.StartData, "numeric data to start from")
	iterations := flag.Int("iterations", cfg.Execution.Iterations, "number of iterations")
	timeout := flag.Int("timeout", cfg.API.Timeout, "per-request timeout seconds")
	verbose := flag.Bool("verbose", cfg.Output.Verbose, "enable debug logging")
	showProgress := flag.Bool("show-progress", cfg.Output.ShowProgress, "show per-iteration progress output")
	showBody := flag.Bool("show-body", cfg.Output.ShowBody, "show request/response URLs and JSON bodies")
	useBulk := flag.Bool("bulk", cfg.Batch.Enabled, "use bulk protect/reveal endpoints")
	batchSize := flag.Int("batch-size", cfg.Batch.Size, "batch size for bulk operations")
	useTLS := flag.Bool("tls", cfg.API.TLS, "use HTTPS instead of HTTP")
	// JWT 플래그: boolean 플래그의 특수성 때문에 문자열로 받아서 처리합니다
	jwtEnabledFlag := flag.String("jwt-enabled", "", "enable JWT authentication (true/false)")
	jwtTokenFlag := flag.String("jwt-token", "", "JWT token for authentication")

	// 커스텀 Usage 함수로 -- 형식 표시
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  --config string\n        path to config.yaml file (default: auto-search)\n")
		fmt.Fprintf(os.Stderr, "  --host string\n        API host (default \"192.168.0.231\")\n")
		fmt.Fprintf(os.Stderr, "  --port int\n        API port (default 32082)\n")
		fmt.Fprintf(os.Stderr, "  --policy string\n        protection_policy_name (default \"P03\")\n")
		fmt.Fprintf(os.Stderr, "  --start-data string\n        numeric data to start from (default \"1234567890123\")\n")
		fmt.Fprintf(os.Stderr, "  --iterations int\n        number of iterations (default 100)\n")
		fmt.Fprintf(os.Stderr, "  --timeout int\n        per-request timeout seconds (default 10)\n")
		fmt.Fprintf(os.Stderr, "  --verbose\n        enable debug logging\n")
		fmt.Fprintf(os.Stderr, "  --show-progress\n        show per-iteration progress output\n")
		fmt.Fprintf(os.Stderr, "  --show-body\n        show request/response URLs and JSON bodies (auto-enables show-progress)\n")
		fmt.Fprintf(os.Stderr, "  --bulk\n        use bulk protect/reveal endpoints\n")
		fmt.Fprintf(os.Stderr, "  --batch-size int\n        batch size for bulk operations (default 50)\n")
		fmt.Fprintf(os.Stderr, "  --tls\n        use HTTPS instead of HTTP\n")
		fmt.Fprintf(os.Stderr, "  --jwt-enabled string\n        enable JWT authentication (true/false)\n")
		fmt.Fprintf(os.Stderr, "  --jwt-token string\n        JWT token for authentication\n")
	}

	flag.Parse()

	// JWT 설정 처리: config.yaml 기본값 + 플래그 오버라이드
	jwtEnabled := cfg.Auth.JWTEnabled
	jwtToken := cfg.Auth.JWTToken
	if *jwtEnabledFlag != "" {
		jwtEnabled = *jwtEnabledFlag == "true"
	}
	if *jwtTokenFlag != "" {
		jwtToken = *jwtTokenFlag
	}

	// show-body가 활성화되면 show-progress도 자동 활성화
	if *showBody {
		*showProgress = true
	}

	// 클라이언트 생성
	c := client.NewClient(*host, *port, *policy, *timeout, *useTLS)
	c.SetShowBody(*showBody)
	
	// verbose 로그 출력
	if *verbose {
		log.Printf("Config loaded: JWT enabled=%v, token=%s", jwtEnabled, jwtToken)
	}
	
	c.SetJWT(jwtEnabled, jwtToken)

	// 반복 실행
	startTime := time.Now()
	
	// 배치 수 미리 계산하여 슬라이스 용량 사전 할당
	expectedBatches := (*iterations + *batchSize - 1) / *batchSize
	results := make([]*runner.IterationResult, 0, expectedBatches)
	successfulItems := 0
	matchedItems := 0
	totalItems := 0
	sumBatchTimes := 0.0

	if *useBulk {
		// Bulk 모드: 입력 데이터 생성
		inputs := generateDataSequence(*startData, *iterations, *verbose)

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
				printBulkProgress(i/ *batchSize+1, len(batch), result.TimeS,
					result.ProtectResponse.StatusCode, result.RevealResponse.StatusCode, result.MatchedCount)
			}
		}
	} else {
		// 일반 모드 (단일 처리)
		dataSequence := generateDataSequence(*startData, *iterations, *verbose)
		for i, data := range dataSequence {
			iterNum := i + 1
			result, err := runner.RunIteration(c, data)
			if err != nil {
				if *verbose {
					log.Printf("Error at iteration %d: %v", iterNum, err)
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
				printIterationProgress(iterNum, data, result.TimeS,
					result.ProtectResponse.StatusCode, result.RevealResponse.StatusCode,
					result.Match, *showBody)
			}
		}
	}

	total := time.Since(startTime)

	// 결과 출력
	if *useBulk {
		printSummary(totalItems, successfulItems, matchedItems, total)
	} else {
		printSummary(len(results), successfulItems, matchedItems, total)
	}
}
