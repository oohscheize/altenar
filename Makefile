DOCKER_COMPOSE := docker-compose
PKGS_UNIT      := $(shell go list ./... | grep -v '/cmd/')  # all except cmd/*
COVER_UNIT     := coverage.out
COVER_INT      := integration.out

.PHONY: run down test integration-test coverage icov lint tidy
default: test

run:
	$(DOCKER_COMPOSE) up -d --build
	@echo "Stack is up  →  http://localhost:8080"

down:
	$(DOCKER_COMPOSE) down -v

test:
	go test $(PKGS_UNIT) \
		-race -coverpkg=./... -covermode=atomic -coverprofile=$(COVER_UNIT)

integration-test:
	go test -tags=integration $(PKGS_UNIT) \
		-coverpkg=./... -covermode=atomic -coverprofile=$(COVER_INT)

coverage: test
	go tool cover -html=$(COVER_UNIT) -o coverage.html
	@echo "Unit-coverage report → coverage.html"
	@if command -v xdg-open >/dev/null 2>&1; then xdg-open coverage.html; \
	elif command -v open >/dev/null 2>&1;     then open coverage.html; fi

icov: integration-test
	go tool cover -html=$(COVER_INT) -o integration_coverage.html
	@echo "Integration-coverage report → integration_coverage.html"
	@if command -v xdg-open >/dev/null 2>&1; then xdg-open integration_coverage.html; \
	elif command -v open >/dev/null 2>&1;     then open integration_coverage.html; fi


tidy:
	go mod tidy

