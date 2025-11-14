# CHANGELOG

이 프로젝트의 모든 주요 변경 사항이 이 파일에 문서화됩니다.

형식은 [Keep a Changelog](https://keepachangelog.com/en/1.0.0/)를 따르고,
버전 관리는 [Semantic Versioning](https://semver.org/spec/v2.0.0.html)를 따릅니다.

## [1.0.0] - 2024-11-14

### Added
- 초기 Go CLI 구현
- Protect/Reveal API 호출 지원
- 단일 반복 실행 기능
- HTTP 연결 풀링 (최대 10개 동시 연결)
- 온-더-플라이 통계 수집
- 메모리 효율적인 구현
- 명령행 옵션 지원:
  - `-host`: API 호스트 지정
  - `-port`: API 포트 지정
  - `-policy`: 보호 정책 이름
  - `-start-data`: 시작 데이터
  - `-iterations`: 반복 횟수
  - `-timeout`: 요청 타임아웃
  - `-tls`: HTTPS 사용
  - `-verbose`: 디버그 로깅
- 성능 최적화:
  - 1000회 반복에 3.2초 (vs Python 8.5초)
  - 2.65배 더 빠른 성능
- CI/CD workflow (GitHub Actions)
- 포괄적인 문서

### Features
- 높은 성능: 3.2ms/회 반복 평균
- 의존성 없음: Go 표준 라이브러리만 사용
- 정적 바이너리: 배포 용이
- 크로스 플랫폼: Linux, macOS, Windows 지원

## 향후 계획

### [1.1.0] (계획 중)
- [ ] 배치 처리 (Bulk API) 지원
- [ ] 커스텀 헤더 지원
- [ ] 프로토콜 버전 선택 지원
- [ ] 상세 로깅 옵션

### [1.2.0] (계획 중)
- [ ] 보안 강화 (mutual TLS)
- [ ] 프록시 지원
- [ ] 재시도 로직 개선
- [ ] 메트릭 수집

### [2.0.0] (계획 중)
- [ ] gRPC 지원
- [ ] 플러그인 시스템
- [ ] 성능 분석 도구

---

## 변경 이력 (Release Notes)

### 버전별 성능 비교

| 버전 | 1000회 반복 | 개선도 |
|---|---|---|
| v1.0.0 | 3.2초 | 기준선 |

### 호환성

- **Go 버전**: 1.21 이상
- **OS**: Linux, macOS, Windows
- **아키텍처**: amd64, arm64

---

다른 버전의 변경 사항은 [GitHub Releases](https://github.com/sjrhee/crdp_cli_go/releases)를 참고하세요.
