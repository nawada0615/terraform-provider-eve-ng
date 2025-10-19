.PHONY: build test test-race clean install fmt vet lint

build:
	go build -o terraform-provider-eve-ng

test:
	go test ./... -v

test-race:
	go test ./... -v -race

clean:
	rm -f terraform-provider-eve-ng

install: build
	mkdir -p ~/.terraform.d/plugins/local/nawada0615/eve-ng/0.1.0/linux_amd64
	cp terraform-provider-eve-ng ~/.terraform.d/plugins/local/nawada0615/eve-ng/0.1.0/linux_amd64/

fmt:
	go fmt ./...

vet:
	go vet ./...

lint:
	golangci-lint run

# Development targets
dev-setup:
	go mod tidy
	go mod download

# Test specific targets
test-unit:
	go test ./eve/... -v

test-client:
	go test ./internal/client/... -v

test-mock:
	go test ./tests/... -v
