.PHONY: build build-backend build-frontend test test-backend test-frontend \
	coverage coverage-backend coverage-frontend \
	run-backend run-frontend lint lint-backend lint-frontend \
	docker-build docker-run

PACKAGE_MANAGER ?= npm

build: build-backend build-frontend

build-backend:
	cd backend && go build ./...

build-frontend:
	cd frontend && $(PACKAGE_MANAGER) run build

test: test-backend test-frontend

test-backend:
	cd backend && go test ./...

test-frontend:
	cd frontend && $(PACKAGE_MANAGER) run test

run-backend:
	cd backend && go run ./cmd/server

run-frontend:
	cd frontend && $(PACKAGE_MANAGER) run dev

lint: lint-backend lint-frontend

lint-backend:
	cd backend && go vet ./...

lint-frontend:
	cd frontend && $(PACKAGE_MANAGER) run lint

coverage: coverage-backend coverage-frontend

coverage-backend:
	cd backend && go test ./... -cover | tee ../coverage-report/backend.txt

coverage-frontend:
	cd frontend && $(PACKAGE_MANAGER) run test:coverage | tee ../coverage-report/frontend.txt

docker-build:
	docker compose build

docker-run:
	docker compose up
