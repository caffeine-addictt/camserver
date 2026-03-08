# Syntax
#
# <target>: [dep targets] [## docstring]
# <target>: [dep targets] \
# 	... \
# 	[## multilined]

all: help

.PHONY: help
help: ## prints this message
	@echo "Usage: make [target] [args...]"
	@printf "\nAvailable targets:\n"
	@awk ' \
		BEGIN { FS = ":.*?## " } \
		/^[a-zA-Z0-9_\/\.-]+:/ { \
			target = $$1; \
			gsub(/[ :\\].*/, "", target); \
			desc = ""; \
			if ($$0 ~ /## /) { \
				desc = $$0; \
				sub(/.*## /, "", desc); \
			} else { \
				while ($$0 ~ /\\$$/) { \
					if (getline <= 0) break; \
					if ($$0 ~ /[ \t]*## /) { \
						desc = $$0; \
						sub(/^[ \t]*## /, "", desc); \
						break; \
					} \
				} \
			} \
			if (desc) printf "  \033[36m%-15s\033[0m %s\n", target, desc; \
		} \
	' $(MAKEFILE_LIST)

.PHONY: build
build: ## builds our binaries
	ls -1 ./cmd | xargs -I{} \
		go build -ldflags="-s -w" -o "./bin/{}" "./cmd/{}"

.PHONY: test
test: ## run tests
	go mod tidy
	go mod verify
	go vet ./...
	go test -v -race ./...

.PHONY: bench
bench: ## benchmarking
	go test -v -bench=. -benchmem ./...

.PHONY: security
security: ## run security checks
	go run github.com/securego/gosec/v2/cmd/gosec@latest -quiet ./...
	go run github.com/go-critic/go-critic/cmd/gocritic@latest check -enableAll ./...
	go run github.com/google/osv-scanner/cmd/osv-scanner@latest -r .

.PHONY: lint
lint: ## static code analysis
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@latest run

.PHONY: format
format: ## enforcing code style
	go fmt ./...
	go mod tidy -v
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@latest run --fix

.PHONY: clean
clean: ## clean up aritifacts
	go clean ./...
	rm -rf ./bin
