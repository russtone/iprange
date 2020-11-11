#
# Build
#

.PHONY: build
build:
	@go build

#
# Test
#

.PHONY: test
test:
	@go test ./... -v -coverprofile coverage.out

.PHONY: coverage
coverage:
	@go tool cover -func=coverage.out

.PHONY: coverage-html
coverage-html:
	@go tool cover -html=coverage.out

