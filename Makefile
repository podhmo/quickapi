test:
	@mkdir -p /tmp/xxx
	GOBIN=/tmp/xxx go install -v ./...
	go test ./...