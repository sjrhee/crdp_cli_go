package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

// APIResponse는 API 응답을 나타냅니다
type APIResponse struct {
	StatusCode int
	Body       map[string]interface{}
}

// Client는 CRDP API 클라이언트입니다
type Client struct {
	baseURL  string
	policy   string
	timeout  time.Duration
	client   *http.Client
	showBody bool
}

// NewClient는 새로운 CRDP 클라이언트를 생성합니다
func NewClient(host string, port int, policy string, timeoutSec int, useTLS bool) *Client {
	protocol := "http"
	if useTLS {
		protocol = "https"
	}

	baseURL := fmt.Sprintf("%s://%s:%d", protocol, host, port)

	// HTTP 클라이언트 설정 (TCP_NODELAY 효과를 내기 위해 커스텀 Transport 사용)
	transport := &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   time.Duration(timeoutSec) * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		MaxIdleConns:        10,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
		DisableKeepAlives:   false,
	}

	// TLS 사용 시 인증서 검증 비활성화
	if useTLS {
		transport.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	httpClient := &http.Client{
		Transport: transport,
		Timeout:   time.Duration(timeoutSec) * time.Second,
	}

	return &Client{
		baseURL:  baseURL,
		policy:   policy,
		timeout:  time.Duration(timeoutSec) * time.Second,
		client:   httpClient,
		showBody: false,
	}
}

// SetShowBody는 요청/응답 본문 출력 여부를 설정합니다
func (c *Client) SetShowBody(show bool) {
	c.showBody = show
}

// PostJSON은 JSON 페이로드로 POST 요청을 보냅니다
func (c *Client) PostJSON(endpoint string, payload map[string]interface{}) (*APIResponse, error) {
	url := c.baseURL + endpoint

	// JSON 인코딩
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	// show-body 옵션이 활성화된 경우 요청 정보 출력
	if c.showBody {
		method := "POST"
		fmt.Printf("%s %s\n%s\n", method, url, string(body))
	}

	// POST 요청 생성
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	// 요청 전송
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 응답 본문 읽기
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// JSON 파싱
	var data map[string]interface{}
	if len(respBody) > 0 {
		if err := json.Unmarshal(respBody, &data); err != nil {
			// JSON 파싱 실패 시 문자열로 저장
			data = map[string]interface{}{"raw": string(respBody)}
		}
	}

	// show-body 옵션이 활성화된 경우 응답 정보 출력
	if c.showBody {
		fmt.Print(string(respBody))
	}

	return &APIResponse{
		StatusCode: resp.StatusCode,
		Body:       data,
	}, nil
}

// Protect는 데이터를 보호합니다
func (c *Client) Protect(data string) (*APIResponse, error) {
	payload := map[string]interface{}{
		"data":                      data,
		"protection_policy_name": c.policy,
	}
	return c.PostJSON("/v1/protect", payload)
}

// Reveal은 보호된 데이터를 복원합니다
func (c *Client) Reveal(protectedData string) (*APIResponse, error) {
	payload := map[string]interface{}{
		"protected_data":              protectedData,
		"protection_policy_name": c.policy,
	}
	return c.PostJSON("/v1/reveal", payload)
}

// ProtectBulk는 여러 데이터를 한 번에 보호합니다 (Thales API 형식)
func (c *Client) ProtectBulk(dataList []string) (*APIResponse, error) {
	payload := map[string]interface{}{
		"protection_policy_name": c.policy,
		"data_array":             dataList,
	}
	return c.PostJSON("/v1/protectbulk", payload)
}

// RevealBulk은 여러 보호된 데이터를 한 번에 복원합니다 (Thales API 형식)
func (c *Client) RevealBulk(protectedDataList []string) (*APIResponse, error) {
	// protected_data_array 형태로 구성
	pdArray := make([]map[string]interface{}, len(protectedDataList))
	for i, pd := range protectedDataList {
		pdArray[i] = map[string]interface{}{
			"protected_data": pd,
		}
	}

	payload := map[string]interface{}{
		"protection_policy_name": c.policy,
		"protected_data_array":   pdArray,
	}
	return c.PostJSON("/v1/revealbulk", payload)
}
