.PHONY: build run test clean help

# Go 버전 및 바이너리 설정
GO := go
BINARY := crdp-cli
VERSION := v1.0.0
LDFLAGS := -ldflags "-X main.Version=$(VERSION)"

# 플랫폼별 빌드 설정
PLATFORMS := linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64

help:
	@echo "CRDP Go CLI - Makefile targets:"
	@echo ""
	@echo "  make build           - 현재 플랫폼용 바이너리 빌드"
	@echo "  make build-all       - 모든 플랫폼용 바이너리 빌드"
	@echo "  make run             - 바이너리 실행 (기본 설정)"
	@echo "  make test            - 테스트 실행"
	@echo "  make clean           - 빌드 결과물 정리"
	@echo "  make lint            - 코드 스타일 검사"
	@echo "  make fmt             - 코드 포맷팅"
	@echo "  make help            - 이 도움말 출력"

build:
	@echo "Building $(BINARY) for current platform..."
	$(GO) build $(LDFLAGS) -o $(BINARY) ./cmd/crdp-cli
	@echo "✅ Build complete: ./$(BINARY)"

build-all:
	@echo "Building $(BINARY) for all platforms..."
	@mkdir -p dist
	@for platform in $(PLATFORMS); do \
		GOOS=$${platform%/*} GOARCH=$${platform##*/} \
		$(GO) build $(LDFLAGS) -o dist/$(BINARY)-$$GOOS-$$GOARCH ./cmd/crdp-cli; \
		echo "✅ Built for $$platform"; \
	done

run: build
	@echo "Running $(BINARY)..."
	./$(BINARY) -iterations 100

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
