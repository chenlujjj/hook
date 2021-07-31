all: clean build

.PHONY: build
build: tidy
		GOOS=$(GOOS) GOARCH=$(GOARCH) GO111MODULE=on go build -o ./bin/hook ./cmd/hook/

.PHONY: tidy
tidy:
		go mod tidy

.PHONY: clean
clean:
		@rm -rf bin/
