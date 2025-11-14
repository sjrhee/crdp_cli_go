# Contributing to CRDP Go CLI

감사합니다! 이 프로젝트에 기여하는 것에 관심을 가져주셨습니다.

## 개발 환경 설정

### 요구사항

- Go 1.21 이상
- git
- Make (선택사항, 하지만 권장됨)

### 빌드 및 실행

```bash
# 저장소 클론
git clone https://github.com/sjrhee/crdp_cli_go.git
cd crdp_cli_go

# 빌드
make build

# 또는 직접 빌드
go build -o crdp-cli ./cmd/crdp-cli

# 실행
./crdp-cli -iterations 100
```

## 코드 스타일

- [Effective Go](https://golang.org/doc/effective_go)를 따릅니다
- `go fmt` 로 코드를 포맷팅합니다
- `go vet` 로 코드를 검사합니다

```bash
# 코드 포맷팅
make fmt

# 코드 검사
go vet ./...
```

## 커밋 메시지

커밋 메시지는 다음 형식을 따릅니다:

```
<type>: <subject>

<body>

<footer>
```

### Type

- `feat`: 새로운 기능
- `fix`: 버그 수정
- `docs`: 문서 작성 또는 수정
- `style`: 코드 스타일 수정 (포맷팅, 누락된 세미콜론 등)
- `refactor`: 코드 리팩터링
- `perf`: 성능 개선
- `test`: 테스트 추가 또는 수정
- `chore`: 빌드 프로세스, 패키지 매니저 설정 등

### 예시

```
feat: Add support for custom headers in API requests

- Implement HeaderModifier interface
- Add command-line flag --custom-headers
- Add comprehensive tests for new functionality

Closes #123
```

## Pull Request 프로세스

1. 저장소를 fork합니다
2. feature 브랜치를 만듭니다 (`git checkout -b feature/amazing-feature`)
3. 변경사항을 커밋합니다 (`git commit -m 'Add amazing feature'`)
4. 브랜치에 push합니다 (`git push origin feature/amazing-feature`)
5. Pull Request를 생성합니다

## PR 검토 체크리스트

PR을 제출하기 전에 다음을 확인하세요:

- [ ] 코드가 `make fmt`로 포맷팅되었나요?
- [ ] `make test`로 테스트를 실행했나요?
- [ ] `go vet ./...`로 코드를 검사했나요?
- [ ] 새로운 기능에 테스트가 추가되었나요?
- [ ] README가 필요한 경우 업데이트되었나요?
- [ ] 변경 사항이 문서화되었나요?

## 테스트

모든 새로운 기능 및 버그 수정에 대해 테스트를 작성하세요.

```bash
# 테스트 실행
make test

# 커버리지 포함 테스트
go test -cover ./...

# 벤치마크 실행
make bench
```

## 문제 보고

버그를 발견했거나 기능을 제안하고 싶다면, 다음 정보를 포함하여 Issue를 생성하세요:

- 명확하고 설명적인 제목
- 문제/기능의 상세한 설명
- 재현 단계 (버그의 경우)
- 예상되는 동작 vs 실제 동작
- 스크린샷 (해당하는 경우)
- 시스템 정보 (OS, Go 버전 등)

## 성능 최적화

성능 관련 PR의 경우, 벤치마크 결과를 포함해주세요:

```bash
make bench
```

## 라이센스

이 프로젝트는 MIT 라이센스로 라이센스되어 있습니다. 기여함으로써 귀하의 기여가 이 라이센스 하에 라이센스될 것에 동의합니다.

## 질문이나 도움이 필요하신가요?

- Issue를 생성하세요
- Discussion을 시작하세요
- 이메일로 연락하세요

감사합니다! 🎉
