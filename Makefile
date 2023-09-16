.PHONY: clean
clean:
	make -C node/service clean
	make -C master/service clean
.PHONY: gen
gen:
	make -C node/api gen
	make -C node/service gen
	make -C master/api gen
	make -C master/service gen

.PHONY: test
test:
	make -C node/service test
	make -C master/service test

.PHONY: build
build:
	make -C node/service build
	make -C master/service build

.PHONY: lint
lint:
	make -C node/service lint
	make -C master/service lint


.PHONY: test-integration
test-integration:
	make -C node/service test-integration
	make -C master/service test-integration

.PHONY: test-integration-race
test-integration-race:
	make -C node/service integration-test-race
	make -C master/service integration-test-race

install-tools:
	go install github.com/golang/mock/mockgen@v1.6.0;
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.49.0;
	brew install openapi-generator;
