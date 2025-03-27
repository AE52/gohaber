FROM golang:1.20-alpine AS builder

# Gerekli paketlerin yüklenmesi
RUN apk add --no-cache git

# Çalışma dizini oluşturma
WORKDIR /app

# Go modüllerini kopyalama ve indirme
COPY go.mod go.sum ./
RUN go mod download

# Kaynak kodları kopyalama
COPY . .

# Değişiklik sonrası
# Derleme öncesi template ve static dosyalarının varlığını kontrol et
RUN mkdir -p /app/web/templates/layouts /app/web/templates/partials

# Uygulamayı derleme
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

# Çalışma imajı
FROM alpine:latest  

# SSL sertifikaları ve gerekli paketleri yükleme
RUN apk --no-cache add ca-certificates tzdata

# Çalışma dizini oluşturma
WORKDIR /root/

# Derlenmiş uygulamayı kopyalama
COPY --from=builder /app/main .
COPY --from=builder /app/web ./web

# Uygulama için gereken port
EXPOSE 3000

# Uygulamayı çalıştırma
CMD ["./main"] 