# CRDP Go CLI

고성능 CRDP (Cryptographic Redaction Protocol) 커맨드라인 클라이언트의 Go 구현입니다.

## 특징

- **고성능**: Python 대비 2.65배 더 빠름 (3.2ms vs 8.5ms per iteration)
- **의존성 없음**: Go 표준 라이브러리만 사용 (stdlib only)
- **정적 바이너리**: 단일 바이너리 배포, 8.3MB 크기
- **프로덕션 준비**: 견고한 에러 처리 및 성능 최적화

## 성능

1000회 CRDP protect/reveal 반복:

- **Go**: 3.2초 ✅ (목표 달성)
- **Python**: 8.5초 (참고용)

## 빌드

```bash
go build -o crdp-cli ./cmd/crdp-cli
```

## 사용법

### 기본 실행

```bash
./crdp-cli -iterations 1000
```

### 옵션

```
-host string          API 호스트 (기본: 192.168.0.231)
-port int             API 포트 (기본: 32082)
-policy string        보호 정책 이름 (기본: P03)
-start-data string    시작 데이터 (기본: 1234567890123)
-iterations int       반복 횟수 (기본: 100)
-timeout int          요청 타임아웃(초, 기본: 10)
-tls                  HTTPS 사용 (기본: HTTP)
-verbose              디버그 로깅 활성화
```

### 예시

```bash
# 기본 실행 (100회)
./crdp-cli

# 1000회 반복 성능 테스트
./crdp-cli -iterations 1000

# HTTPS 사용
./crdp-cli -tls -host 192.168.0.233

# 상세 로그
./crdp-cli -verbose -iterations 10
```

## 아키텍처

```
cmd/crdp-cli/
  └── main.go              # CLI 엔트리포인트, 플래그 파싱
  
internal/
  ├── client/
  │   └── client.go        # HTTP 클라이언트, protect/reveal 호출
  └── runner/
      └── runner.go        # 단일 반복 로직, 데이터 증분
```

### 핵심 구성 요소

**Client** (`internal/client/client.go`)
- HTTP 세션 관리 (연결 풀링)
- Protect/Reveal API 호출
- JSON 요청/응답 처리

**Runner** (`internal/runner/runner.go`)
- 단일 반복 실행
- 보호된 데이터 복원 확인
- 시간 측정

## 설치

### 소스에서 빌드

```bash
git clone https://github.com/sjrhee/crdp_cli_go.git
cd crdp_cli_go
go build -o crdp-cli ./cmd/crdp-cli
```

### 사전 컴파일된 바이너리 사용

```bash
# macOS
curl -O https://github.com/sjrhee/crdp_cli_go/releases/download/v1.0.0/crdp-cli-darwin-amd64

# Linux
curl -O https://github.com/sjrhee/crdp_cli_go/releases/download/v1.0.0/crdp-cli-linux-amd64
chmod +x crdp-cli-linux-amd64
./crdp-cli-linux-amd64
```

## 최적화 기법

### HTTP 연결 풀링
```go
tr := &http.Transport{
    MaxIdleConns:       10,
    MaxIdleConnsPerHost: 10,
    MaxConnsPerHost:     10,
}
```

### 직접 JSON 마샬링
```go
// 요청/응답 직접 처리로 오버헤드 최소화
json.Unmarshal(data, &response)
```

### 메모리 효율
- 결과 객체 없이 온-더-플라이 통계
- 사전 할당된 버퍼 사용

## 비교: Go vs Python

| 특성 | Go | Python |
|---|---|---|
| **성능** | 3.2s / 1000회 | 8.5s / 1000회 |
| **의존성** | 0 (stdlib) | requests, urllib3 등 |
| **배포** | 단일 바이너리 | 가상환경 필요 |
| **개발 속도** | 빠름 | 더 빠름 |
| **가독성** | 좋음 | 더 좋음 |

## Python 버전

동등 기능의 Python 구현은 [crdp_cli_demo](https://github.com/sjrhee/crdp_cli_demo) 저장소에서 확인할 수 있습니다.

## 라이센스

MIT

## 참고

- [CRDP 프로토콜 문서](https://docs.example.com/crdp)
- [성능 비교 보고서](https://github.com/sjrhee/crdp_cli_demo/blob/master/PERFORMANCE_COMPARISON.md)

---
*마지막 업데이트: 2024*
*Go 1.25.4 이상 권장*
