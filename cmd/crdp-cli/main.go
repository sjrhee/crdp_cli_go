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
	// 모든 플래그 정의 (기본값은 공백/0/false로 설정)
	configPath := flag.String("config", "", "path to config.yaml file (default: auto-search)")
	
	// CLI 플래그 정의 (기본값을 빈 값이나 0으로 설정하여 명시적 제공 여부 감지)
	host := flag.String("host", "", "API host")
	port := flag.Int("port", 0, "API port")
	policy := flag.String("policy", "", "protection_policy_name")
	startData := flag.String("start-data", "", "numeric data to start from")
	iterations := flag.Int("iterations", 0, "number of iterations")
	timeout := flag.Int("timeout", 0, "per-request timeout seconds")
	verbose := flag.Bool("verbose", false, "enable debug logging")
	showProgress := flag.Bool("show-progress", false, "show per-iteration progress output")
	showBody := flag.Bool("show-body", false, "show request/response URLs and JSON bodies")
	useBulk := flag.Bool("bulk", false, "use bulk protect/reveal endpoints")
	batchSize := flag.Int("batch-size", 0, "batch size for bulk operations")
	useTLSFlag := flag.String("tls", "", "use HTTPS (true/false, default: config value)")
	jwtEnabledFlag := flag.String("jwt-enabled", "", "enable JWT authentication (true/false)")
	jwtTokenFlag := flag.String("jwt-token", "", "JWT token for authentication")

	// 커스텀 Usage 함수
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  --config string          path to config.yaml file (default: auto-search)\n")
		fmt.Fprintf(os.Stderr, "  --host string            API host (default \"192.168.0.231\")\n")
		fmt.Fprintf(os.Stderr, "  --port int               API port (default 32082)\n")
		fmt.Fprintf(os.Stderr, "  --policy string          protection_policy_name (default \"P03\")\n")
		fmt.Fprintf(os.Stderr, "  --start-data string      numeric data to start from (default \"1234567890123\")\n")
		fmt.Fprintf(os.Stderr, "  --iterations int         number of iterations (default 100)\n")
		fmt.Fprintf(os.Stderr, "  --timeout int            per-request timeout seconds (default 10)\n")
		fmt.Fprintf(os.Stderr, "  --verbose                enable debug logging\n")
		fmt.Fprintf(os.Stderr, "  --show-progress          show per-iteration progress output\n")
		fmt.Fprintf(os.Stderr, "  --show-body              show request/response URLs and JSON bodies\n")
		fmt.Fprintf(os.Stderr, "  --bulk                   use bulk protect/reveal endpoints\n")
		fmt.Fprintf(os.Stderr, "  --batch-size int         batch size for bulk operations (default 50)\n")
		fmt.Fprintf(os.Stderr, "  --tls string             use HTTPS (true/false, default: config value)\n")
		fmt.Fprintf(os.Stderr, "  --jwt-enabled string     enable JWT authentication (true/false)\n")
		fmt.Fprintf(os.Stderr, "  --jwt-token string       JWT token for authentication\n")
	}

	// 플래그 파싱
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

	// CLI 플래그로 설정된 값이 있으면 config 값을 오버라이드
	// flag.Visit()를 사용하여 실제로 명시적으로 제공된 플래그만 처리
	flag.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "host":
			if *host != "" {
				cfg.API.Host = *host
			}
		case "port":
			if *port != 0 {
				cfg.API.Port = *port
			}
		case "policy":
			if *policy != "" {
				cfg.Protection.Policy = *policy
			}
		case "start-data":
			if *startData != "" {
				cfg.Execution.StartData = *startData
			}
		case "iterations":
			if *iterations != 0 {
				cfg.Execution.Iterations = *iterations
			}
		case "timeout":
			if *timeout != 0 {
				cfg.API.Timeout = *timeout
			}
		case "verbose":
			cfg.Output.Verbose = *verbose
		case "show-progress":
			cfg.Output.ShowProgress = *showProgress
		case "show-body":
			cfg.Output.ShowBody = *showBody
		case "bulk":
			cfg.Batch.Enabled = *useBulk
		case "batch-size":
			if *batchSize != 0 {
				cfg.Batch.Size = *batchSize
			}
		case "tls":
			if *useTLSFlag != "" {
				cfg.API.TLS = *useTLSFlag == "true"
			}
		case "jwt-enabled":
			if *jwtEnabledFlag != "" {
				// jwt-enabled 플래그는 별도 처리
			}
		case "jwt-token":
			if *jwtTokenFlag != "" {
				// jwt-token 플래그는 별도 처리
			}
		}
	})

	// JWT 설정 처리
	jwtEnabled := cfg.Auth.JWTEnabled
	jwtToken := cfg.Auth.JWTToken
	if *jwtEnabledFlag != "" {
		jwtEnabled = *jwtEnabledFlag == "true"
	}
	if *jwtTokenFlag != "" {
		jwtToken = *jwtTokenFlag
	}

	// show-body가 활성화되면 show-progress도 자동 활성화
	if cfg.Output.ShowBody {
		cfg.Output.ShowProgress = true
	}

	// 클라이언트 생성
	c := client.NewClient(cfg.API.Host, cfg.API.Port, cfg.Protection.Policy, cfg.API.Timeout, cfg.API.TLS)
	c.SetShowBody(cfg.Output.ShowBody)
	
	// verbose 로그 출력
	if cfg.Output.Verbose {
		log.Printf("Config loaded: JWT enabled=%v", jwtEnabled)
	}
	
	c.SetJWT(jwtEnabled, jwtToken)

	// 반복 실행
	startTime := time.Now()
	
	// 배치 수 미리 계산하여 슬라이스 용량 사전 할당
	expectedBatches := (cfg.Execution.Iterations + cfg.Batch.Size - 1) / cfg.Batch.Size
	results := make([]*runner.IterationResult, 0, expectedBatches)
	successfulItems := 0
	matchedItems := 0
	totalItems := 0
	sumBatchTimes := 0.0

	if cfg.Batch.Enabled {
		// Bulk 모드: 입력 데이터 생성
		inputs := generateDataSequence(cfg.Execution.StartData, cfg.Execution.Iterations, cfg.Output.Verbose)

		// 배치 단위로 처리
		for i := 0; i < len(inputs); i += cfg.Batch.Size {
			end := i + cfg.Batch.Size
			if end > len(inputs) {
				end = len(inputs)
			}
			batch := inputs[i:end]

			result, err := runner.RunBulkIteration(c, batch)
			if err != nil {
			if cfg.Output.Verbose {
				log.Printf("Error at batch %d: %v", i/cfg.Batch.Size+1, err)
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

			if cfg.Output.ShowProgress {
				printBulkProgress(i/cfg.Batch.Size+1, len(batch), result.TimeS,
					result.ProtectResponse.StatusCode, result.RevealResponse.StatusCode, result.MatchedCount)
			}
		}
	} else {
		// 일반 모드 (단일 처리)
		dataSequence := generateDataSequence(cfg.Execution.StartData, cfg.Execution.Iterations, cfg.Output.Verbose)
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

			if cfg.Output.ShowProgress {
				printIterationProgress(iterNum, data, result.TimeS,
					result.ProtectResponse.StatusCode, result.RevealResponse.StatusCode,
					result.Match, cfg.Output.ShowBody)
			}
		}
	}

	total := time.Since(startTime)

	// 결과 출력
	if cfg.Batch.Enabled {
		printSummary(totalItems, successfulItems, matchedItems, total)
	} else {
		printSummary(len(results), successfulItems, matchedItems, total)
	}
}
