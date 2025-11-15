# CRDP Go CLI

CRDP (Cryptographic Redaction Protocol) 커맨드라인 클라이언트의 Go 구현입니다.

## 특징

- **의존성 없음**: Go 표준 라이브러리만 사용
- **정적 바이너리**: 단일 실행 파일로 배포
- **프로덕션 준비**: 견고한 에러 처리 및 최적화

## 설치

### 소스에서 빌드

```bash
git clone https://github.com/sjrhee/crdp_cli_go.git
cd crdp_cli_go
go build -o crdp-cli ./cmd/crdp-cli
```

## 설정

### config.yaml

애플리케이션은 `config.yaml` 파일에서 기본값을 읽습니다. 설정 파일은 다음 순서로 검색됩니다:

1. 현재 디렉토리 (`./config.yaml`)
2. 홈 디렉토리 (`~/.crdp/config.yaml`)
3. 실행 파일 디렉토리 (`<executable_dir>/config.yaml`)

**crdp_file_converter 호환성**: `config.yaml` 구조는 [crdp_file_converter](https://github.com/sjrhee/crdp_file_converter)와 호환됩니다. 동일한 설정 파일을 양쪽 애플리케이션에서 사용할 수 있습니다.

**config.yaml 예시:**

```yaml
# API 연결 설정
api:
  host: "192.168.0.231"
  port: 32082
  timeout: 10
  tls: false

# 보호 정책 설정
protection:
  policy: "P03"

# 반복 실행 설정
execution:
  iterations: 100
  start_data: "1234567890123"

# 배치 처리 설정
batch:
  enabled: false
  size: 50

# 출력 설정
output:
  show_progress: false
  show_body: false
  verbose: false

# 파일 처리 설정 (crdp_file_converter 호환)
file:
  delimiter: ","
  column: 0
  skip_header: false

# 병렬 처리 설정 (crdp_file_converter 호환)
parallel:
  workers: 1
```

CLI 플래그는 `config.yaml`의 기본값을 **덮어씁니다**.

## 사용법

### 기본 실행

```bash
# 기본 설정으로 100회 실행 (config.yaml 사용)
./crdp-cli

# 반복 횟수 지정
./crdp-cli --iterations 1000
```

### 옵션

| 옵션 | 설명 | 기본값 |
|------|------|--------|
| `--host` | API 호스트 주소 | 192.168.0.231 |
| `--port` | API 포트 번호 | 32082 |
| `--policy` | 보호 정책 이름 | P03 |
| `--start-data` | 시작 데이터 (숫자 문자열) | 1234567890123 |
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

# 특정 데이터부터 시작
./crdp-cli --start-data 1234567890123 --iterations 10

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
│   ├── config/
│   │   └── config.go         # 설정 파일 로더
│   └── runner/
│       └── runner.go         # 실행 로직 및 검증
├── config.yaml               # 설정 파일
├── go.mod
└── README.md
```

### 주요 컴포넌트

- **Client**: CRDP API와의 통신 담당 (protect/reveal 호출)
- **Config**: `config.yaml` 파일 읽기 및 기본값 관리
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
