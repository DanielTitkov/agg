.PHONY: test
test:
	go test ./... -cover

.PHONY: lint
lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.29.0
	./bin/golangci-lint run -v