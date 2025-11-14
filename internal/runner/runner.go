package runner

import (
	"time"

	"github.com/sjrhee/crdp-cli-go/internal/client"
)

// IterationResult는 한 번의 반복 결과를 나타냅니다
type IterationResult struct {
	Data              string
	ProtectResponse   *client.APIResponse
	RevealResponse    *client.APIResponse
	ProtectedToken    string
	Restored          string
	TimeS             float64
	Success           bool
	Match             bool
}

// RunIteration은 한 번의 protect->reveal 반복을 실행합니다
func RunIteration(c *client.Client, data string) (*IterationResult, error) {
	start := time.Now()

	// Protect 요청
	protectResp, err := c.Protect(data)
	if err != nil {
		return nil, err
	}

	var protectedData string
	if protectResp.Body != nil {
		if pd, ok := protectResp.Body["protected_data"].(string); ok {
			protectedData = pd
		}
	}

	// Reveal 요청
	revealResp, err := c.Reveal(protectedData)
	if err != nil {
		return nil, err
	}

	var restoredData string
	if revealResp.Body != nil {
		if rd, ok := revealResp.Body["data"].(string); ok {
			restoredData = rd
		}
	}

	elapsed := time.Since(start).Seconds()

	result := &IterationResult{
		Data:            data,
		ProtectResponse: protectResp,
		RevealResponse:  revealResp,
		ProtectedToken:  protectedData,
		Restored:        restoredData,
		TimeS:           elapsed,
		Success:         protectResp.StatusCode >= 200 && protectResp.StatusCode < 300 &&
			revealResp.StatusCode >= 200 && revealResp.StatusCode < 300,
		Match: restoredData == data,
	}

	return result, nil
}

// IsSuccess는 응답이 2xx 상태 코드를 가지고 있는지 확인합니다
func IsSuccess(statusCode int) bool {
	return statusCode >= 200 && statusCode < 300
}
