package runner

import (
	"time"

	"github.com/sjrhee/crdp-cli-go/internal/client"
)

// IterationResult는 한 번의 반복 결과를 나타냅니다
type IterationResult struct {
	Data            string
	ProtectResponse *client.APIResponse
	RevealResponse  *client.APIResponse
	ProtectedToken  string
	Restored        string
	TimeS           float64
	Success         bool
	Match           bool
	RestoredCount   int // bulk용: 복원된 항목 수
	MatchedCount    int // bulk용: 일치하는 항목 수
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

// RunBulkIteration은 배치 단위로 bulk protect->reveal 반복을 실행합니다
func RunBulkIteration(c *client.Client, batch []string) (*IterationResult, error) {
	start := time.Now()

	// Bulk Protect 요청
	protectResp, err := c.ProtectBulk(batch)
	if err != nil {
		return nil, err
	}

	// protected_data_array 추출
	protectedList := extractProtectedList(protectResp)

	// Bulk Reveal 요청
	revealResp, err := c.RevealBulk(protectedList)
	if err != nil {
		return nil, err
	}

	// data_array 추출
	restoredList := extractRestoredList(revealResp)

	elapsed := time.Since(start).Seconds()

	// 일치 항목 카운트
	matchedCount := 0
	maxLen := len(batch)
	if len(restoredList) < maxLen {
		maxLen = len(restoredList)
	}
	for i := 0; i < maxLen; i++ {
		if batch[i] == restoredList[i] {
			matchedCount++
		}
	}

	result := &IterationResult{
		ProtectResponse: protectResp,
		RevealResponse:  revealResp,
		TimeS:           elapsed,
		Success:         protectResp.StatusCode >= 200 && protectResp.StatusCode < 300 &&
			revealResp.StatusCode >= 200 && revealResp.StatusCode < 300,
		Match:         matchedCount == len(batch) && len(restoredList) == len(batch),
		RestoredCount: len(restoredList),
		MatchedCount:  matchedCount,
	}

	return result, nil
}

// extractProtectedList는 protect 응답에서 protected_data_array를 추출합니다
func extractProtectedList(resp *client.APIResponse) []string {
	if resp == nil || resp.Body == nil {
		return []string{}
	}

	// protected_data_array 형태 지원
	if pdArray, ok := resp.Body["protected_data_array"].([]interface{}); ok {
		result := make([]string, 0, len(pdArray))
		for _, item := range pdArray {
			if itemMap, ok := item.(map[string]interface{}); ok {
				if pd, ok := itemMap["protected_data"].(string); ok {
					result = append(result, pd)
				}
			}
		}
		return result
	}

	return []string{}
}

// extractRestoredList는 reveal 응답에서 data_array를 추출합니다
// 최적화: 슬라이스 용량 사전 할당으로 메모리 재할당 최소화
func extractRestoredList(resp *client.APIResponse) []string {
	if resp == nil || resp.Body == nil {
		return []string{}
	}

	// data_array 형태 지원
	if dataArray, ok := resp.Body["data_array"].([]interface{}); ok {
		result := make([]string, 0, len(dataArray))
		for _, item := range dataArray {
			if itemMap, ok := item.(map[string]interface{}); ok {
				if data, ok := itemMap["data"].(string); ok {
					result = append(result, data)
				}
			}
		}
		return result
	}

	return []string{}
}
