# CRDP Go CLI

CRDP (Cryptographic Redaction Protocol) 커맨드라인 클라이언트의 Go 구현입니다.

## 특징

- **의존성 없음**: Go 표준 라이브러리만 사용
- **정적 바이너리**: 단일 실행 파일로 배포 (약 8.3MB)
- **프로덕션 준비**: 견고한 에러 처리 및 최적화

## 설치

### 소스에서 빌드

```bash
git clone https://github.com/sjrhee/crdp_cli_go.git
cd crdp_cli_go
go build -o crdp-cli ./cmd/crdp-cli
```

### 사전 컴파일된 바이너리

```bash
# macOS
curl -LO https://github.com/sjrhee/crdp_cli_go/releases/latest/download/crdp-cli-darwin-amd64
chmod +x crdp-cli-darwin-amd64

# Linux
curl -LO https://github.com/sjrhee/crdp_cli_go/releases/latest/download/crdp-cli-linux-amd64
chmod +x crdp-cli-linux-amd64
```

## 사용법

### 기본 실행

```bash
# 기본 설정으로 100회 실행
./crdp-cli

# 반복 횟수 지정
./crdp-cli -iterations 1000
```

### 옵션

| 옵션 | 설명 | 기본값 |
|------|------|--------|
| `--host` | API 호스트 주소 | 192.168.0.231 |
| `--port` | API 포트 번호 | 32082 |
| `--policy` | 보호 정책 이름 | P03 |
| `--iterations` | 반복 횟수 | 100 |
| `--timeout` | 요청 타임아웃 (초) | 10 |
| `--tls` | HTTPS 사용 | false |
| `--verbose` | 상세 로그 출력 | false |
| `--show-progress` | 반복별 진행 상황 출력 | false |
| `--show-body` | HTTP 요청/응답 본문 출력 (자동으로 show-progress 활성화) | false |
| `--bulk` | 대량 처리 API 사용 (protectbulk/revealbulk) | false |
| `--batch-size` | 대량 처리 시 배치 크기 | 50 |

### 사용 예시

```bash
# 기본 실행 (100회)
./crdp-cli

# HTTPS로 1000회 실행
./crdp-cli --tls --port 32182 --iterations 1000

# 다른 호스트에 연결
./crdp-cli --host 192.168.0.233 --port 32082

# 진행 상황 확인
./crdp-cli --iterations 10 --show-progress

# HTTP 요청/응답 본문 확인
./crdp-cli --iterations 3 --show-body

# 대량 처리 모드 (배치 크기 100)
./crdp-cli --bulk --batch-size 100 --iterations 1000
```

## 프로젝트 구조

```
crdp_cli_go/
├── cmd/
│   └── crdp-cli/
│       └── main.go           # 진입점 및 CLI 인터페이스
├── internal/
│   ├── client/
│   │   └── client.go         # CRDP API 클라이언트
│   └── runner/
│       └── runner.go         # 실행 로직 및 검증
├── go.mod
└── README.md
```

### 주요 컴포넌트

- **Client**: CRDP API와의 통신 담당 (protect/reveal 호출)
- **Runner**: 반복 실행 및 데이터 검증 로직

## 빌드 옵션

```bash
# 개발 빌드
go build -o crdp-cli ./cmd/crdp-cli

# 프로덕션 빌드 (최적화)
go build -ldflags="-s -w" -o crdp-cli ./cmd/crdp-cli
```

## 라이센스

MIT
