.PHONY: build
build:
	go build -o build/bin cmd/app/main.go

.PHONY: run
run: build
	build/bin

.PHONY: docs
docs:
	swag init -g .cmd/app/main.go -o ./docs --parseDependency --parseInternal

.PHONY: gen
gen: mock docs

.PHONY: integration-test
integration-test:
	export CONFIG_PATH=../../config.yaml && \
	go test ./... -v -tags=integration

.PHONY: prod
prod:
	docker compose up --build -d

.PHONY: migrate
migrate: 
	migrate -source file://schema/migrations -database postgres://postgres:5432@127.0.0.1:5432/eff_mobile\?sslmode=disable up

.PHONY: compose-migrate
compose-migrate:
	docker compose up --build migrations