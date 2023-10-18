clean:
	go clean -testcache
	go mod tidy
	rm -rf bin/
gen: clean
	go generate -x ./...
lint:
	golangci-lint run -v --timeout=3m  ./...
build: clean gen
	go build -o bin/grading cmd/grading/*
run: clean gen
	go run ./cmd/grading/
test: clean gen
	go test -count 1 -parallel 8 ./...
test-integration: clean gen
	$(eval current_dir=$(shell pwd))
	(CONFIG_FILE_PATH=$(current_dir)/tests/integration/ go test -count 1 -parallel 8 --tags=integration ./tests/...)
test-integrations: gen test-integration
test-integration-race:
	$(eval current_dir=$(shell pwd))
	(CONFIG_FILE_PATH=$(current_dir)/tests/integration/ go test -count 1 -race --tags=integration ./...)
test-integrations-race: gen test-integration-race
test-parallel: gen
	go test -count 1 -parallel 8 ./...
test-race-cond: gen
	go test -count 1 -race -parallel 8 ./...

test-coveralls:
	# runs all the tests with coverage and sends to coveralls
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	goveralls -coverprofile=coverage.out -service=travis-ci -repotoken $(COVERALLS_TOKEN)

# These step are needed for CodeFresh
# do not change the target name
integration-test: test-integration
integration-test-race: test-integration-race
