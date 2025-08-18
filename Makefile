# Базовые команды для Go-проекта
.PHONY: help lint test build run clean

help:  ## Показать справку по командам
	@echo "Использование: make [цель]"
	@echo ""
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

lint:  ## Запустить линтер
	@golangci-lint run

test:  ## Запустить тесты
	@go test -v ./...

build:  ## Собрать бинарник
	@go build -o bin/app cmd/main.go

run: build  ## Запустить приложение
	@./bin/app

clean:  ## Очистить сгенерированные файлы
	@rm -rf bin/