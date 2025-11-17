.PHONY: build run test clean help build-all windows linux darwin bench update-deps version fmt lint jwt

# Go 버전 및 바이너리 설정
GO := go
BINARY := crdp-cli
VERSION := v1.0.0
LDFLAGS := -ldflags "-s -w -X main.Version=$(VERSION)"

# 플랫폼별 빌드 설정
PLATFORMS := linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64

help:
	@echo "CRDP Go CLI - Makefile targets:"
	@echo ""
	@echo "  make build           - 현재 플랫폼용 바이너리 빌드"
	@echo "  make build-all       - 모든 플랫폼용 바이너리 빌드"
	@echo "  make windows         - Windows 64비트용 빌드"
	@echo "  make linux           - Linux 64비트용 빌드"
	@echo "  make darwin          - macOS 64비트용 빌드"
	@echo "  make run             - 바이너리 실행 (기본 설정)"
	@echo "  make jwt             - JWT 토큰 생성 (ECDSA P-256)"
	@echo "  make test            - 테스트 실행"
	@echo "  make bench           - 벤치마크 실행"
	@echo "  make clean           - 빌드 결과물 정리"
	@echo "  make fmt             - 코드 포맷팅"
	@echo "  make lint            - 코드 스타일 검사"
	@echo "  make update-deps     - 의존성 업데이트"
	@echo "  make version         - 버전 정보 출력"
	@echo "  make help            - 이 도움말 출력"

build:
	@echo "Building $(BINARY) for current platform..."
	$(GO) build $(LDFLAGS) -o $(BINARY) ./cmd/crdp-cli
	@echo "✅ Build complete: ./$(BINARY)"

build-all:
	@echo "Building $(BINARY) for all platforms..."
	@mkdir -p dist
	@for platform in $(PLATFORMS); do \
		GOOS=$${platform%/*}; \
		GOARCH=$${platform##*/}; \
		output_name=dist/$(BINARY)-$$GOOS-$$GOARCH; \
		if [ "$$GOOS" = "windows" ]; then \
			output_name=$${output_name}.exe; \
		fi; \
		GOOS=$$GOOS GOARCH=$$GOARCH $(GO) build $(LDFLAGS) -o $$output_name ./cmd/crdp-cli; \
		echo "✅ Built for $$platform -> $$output_name"; \
	done
	@echo "✅ All builds complete in dist/"

# 개별 플랫폼 빌드
windows:
	@echo "Building for Windows amd64..."
	@mkdir -p dist
	GOOS=windows GOARCH=amd64 $(GO) build $(LDFLAGS) -o dist/$(BINARY)-windows-amd64.exe ./cmd/crdp-cli
	@echo "✅ Built: dist/$(BINARY)-windows-amd64.exe"

linux:
	@echo "Building for Linux amd64..."
	@mkdir -p dist
	GOOS=linux GOARCH=amd64 $(GO) build $(LDFLAGS) -o dist/$(BINARY)-linux-amd64 ./cmd/crdp-cli
	@echo "✅ Built: dist/$(BINARY)-linux-amd64"

darwin:
	@echo "Building for macOS amd64..."
	@mkdir -p dist
	GOOS=darwin GOARCH=amd64 $(GO) build $(LDFLAGS) -o dist/$(BINARY)-darwin-amd64 ./cmd/crdp-cli
	@echo "✅ Built: dist/$(BINARY)-darwin-amd64"

run: build
	@echo "Running $(BINARY)..."
	./$(BINARY) --iterations 10

jwt:
	@echo "Generating JWT token (ECDSA P-256)..."
	@bash scripts/generate_jwt.sh

test:
	@echo "Running tests..."
	$(GO) test -v ./...

clean:
	@echo "Cleaning build artifacts..."
	rm -f $(BINARY)
	rm -rf dist/
	$(GO) clean
	@echo "✅ Clean complete"

lint:
	@echo "Running linter..."
	golangci-lint run ./...

fmt:
	@echo "Formatting code..."
	$(GO) fmt ./...
	@echo "✅ Format complete"

# 성능 벤치마크
bench:
	@echo "Running benchmarks..."
	$(GO) test -bench=. -benchmem ./...

# 의존성 업데이트
update-deps:
	$(GO) get -u
	$(GO) mod tidy

# 버전 확인
version:
	@echo "CRDP Go CLI $(VERSION)"
