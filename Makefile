.PHONY: build
build:
	go build -o build/bin cmd/app/main.go

.PHONY: run
run: build
	build/bin

.PHONY: docs
docs:
	swag init -g ./app/main.go -o ./docs --parseDependency --parseInternal

.PHONY: gen
gen: mock docs

.PHONY: dev
dev:
	docker compose up --build postgres service

.PHONY: full-compose-up
full-compose-up:
	docker-compose up --build -d 

.PHONY: migrate
migrate: 
	migrate -source file://schema/migrations -database postgres://postgres:5432@127.0.0.1:5432/eff_mobile\?sslmode=disable up

.PHONY: compose-migrate
compose-migrate:
	docker compose up --build migrations

.PHONY: init-platform-submodule
init-platform-submodule:
	git submodule update --init --remote platform

.PHONY: update-platform-submodule
update-platform-submodule:
	git submodule update --remote platform
