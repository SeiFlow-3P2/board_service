generate_proto:
	buf generate

tidy:
	go mod tidy

lint: tidy
	go vet ./...

format: tidy
	go fmt ./...

env:
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo "Создан файл .env из .env.example"; \
	else \
		echo "Файл .env уже существует"; \
	fi

board-build:
	docker compose --env-file .env build
	@echo "Сборка завершена"

board-up:
	docker compose --env-file .env up -d
	@echo "Сервис запущен"

board-down:
	docker compose --env-file .env down
	@echo "Сервис остановлен"