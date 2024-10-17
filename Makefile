.PHONY: build
build:
	@echo "Building..."
	@goreleaser build --snapshot --clean

.PHONY: fmt
fmt:
	golines . --write-output --max-len=80 --base-formatter="gofmt" --tab-len=2
	golangci-lint run --fix

.PHONY: test
test:
	go test -v -cover ./...

.PHONY: all
all: build

.PHONY: build-dev
build-dev:
	@echo "Building..."
	@go build -o dist/ ./...
	@echo "Done!"

.PHONY: dev
dev: build-dev
	@echo "Deploying..."
	@cp ./dist/cli ~/.local/bin/golem
	@chmod +x ~/.local/bin/golem
	@echo "Done!"
