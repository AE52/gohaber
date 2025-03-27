#!/bin/bash

echo "Haber Sitesi Başlatılıyor..."

# Gerekli bağımlılıkları kontrol et
echo "Bağımlılıklar kontrol ediliyor..."
go mod tidy

# Ana klasörler oluşturuluyor
mkdir -p web/uploads

# Uygulamayı çalıştır
echo "Uygulama başlatılıyor..."
go run cmd/server/main.go 