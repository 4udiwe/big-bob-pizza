.PHONY: test unit-test integration-test e2e-test cover generate-mocks

# Run all unit tests
unit-test:
	@echo "Running unit tests..."
	go test -v -count=1 ./internal/service/order/...

# Run integration tests
integration-test:
	@echo "Starting integration tests..."
	docker compose -f test/integration_test/docker-compose.test.yaml up \
		--abort-on-container-exit \
		--exit-code-from integration \
		--build \
		--remove-orphans

# Run e2e tests
e2e-test:
	@echo "Starting e2e tests..."
	docker compose -f test/e2e_test/docker-compose.test.yaml up \
		--abort-on-container-exit \
		--exit-code-from e2e \
		--build \
		--remove-orphans


# Generate code coverage
cover:
	@echo "Generating coverage report..."
	go test -count=1 -coverprofile=coverage.out $(shell go list ./internal/service/order/...)
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out

# Generate mocks
generate-mocks:
	@echo "Generating mocks..."
	go generate ./internal/service/order/contracts.go
	go generate ./pkg/transactor/transactor.go


