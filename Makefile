.PHONY: build run clean dev test docker-build docker-up docker-down lint migrate-up migrate-down deps minio-up

# Ana komutlar
build:
	go build -o bin/server ./cmd/server

run: build
	./bin/server

clean:
	rm -rf bin/
	rm -rf tmp/

dev:
	go run cmd/server/main.go

# Docker komutları
docker-build:
	docker-compose build

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

# MinIO servisini başlat
minio-up:
	docker-compose up -d minio createbuckets

# Test komutları
test:
	go test -v ./...

test-cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Veritabanı komutları
migrate-up:
	psql -U postgres -d haberdb -f ./migrations/01_initial_schema.sql

migrate-down:
	@echo "Migration geri alma işlemi yapılıyor..."
	@echo "Bu işlem henüz desteklenmiyor"

# Kalite kontrol
lint:
	golangci-lint run ./...

# Bağımlılıklar
deps:
	go mod tidy

# Yardım
help:
	@echo "Kullanılabilir komutlar:"
	@echo "  make build            - Uygulamayı derler"
	@echo "  make run              - Uygulamayı derler ve çalıştırır"
	@echo "  make clean            - Derleme sonuçlarını temizler"
	@echo "  make dev              - Geliştirme modunda uygulamayı çalıştırır"
	@echo "  make docker-build     - Docker imajını oluşturur"
	@echo "  make docker-up        - Docker containerları başlatır"
	@echo "  make docker-down      - Docker containerları durdurur"
	@echo "  make minio-up         - Sadece MinIO ve bucket oluşturucuyu başlatır"
	@echo "  make test             - Testleri çalıştırır"
	@echo "  make test-cover       - Test kapsama raporunu oluşturur"
	@echo "  make migrate-up       - Veritabanı migrationlarını uygular"
	@echo "  make migrate-down     - Veritabanı migrationlarını geri alır"
	@echo "  make lint             - Kod kalite kontrolü yapar"
	@echo "  make deps             - Bağımlılıkları yönetir"
	@echo "  make help             - Bu yardım mesajını gösterir"

# Varsayılan komut
default: help 