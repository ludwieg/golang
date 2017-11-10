all:
	@go list github.com/stretchr/testify 2>&1 1>/dev/null || go get github.com/stretchr/testify
	@go test