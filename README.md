# Haber Sitesi

Go ve Fiber framework ile geliştirilmiş, SEO dostu ve yüksek performanslı bir haber sitesi projesi.

## Özellikler

- Hızlı ve optimize edilmiş Go Fiber backend
- Responsive tasarım (Bootstrap 5)
- SEO optimizasyonu
- Admin paneli ile içerik yönetimi
- PostgreSQL veritabanı
- Redis önbellek
- JWT tabanlı kimlik doğrulama
- Kullanıcı yorumları
- Kategori ve etiket sistemi
- Reklam alanları yönetimi
- TinyMCE WYSIWYG editör
- Medya yönetim sistemi

## Gereksinimler

- Go 1.18 veya üzeri
- PostgreSQL 13 veya üzeri
- Redis 6 veya üzeri

## Kurulum

1. Proje klonlanır:
   ```bash
   git clone https://github.com/username/haber.git
   cd haber
   ```

2. Bağımlılıklar indirilir:
   ```bash
   go mod tidy
   ```

3. `.env` dosyası oluşturulur:
   ```
   # Sunucu
   SERVER_PORT=3000
   ENVIRONMENT=development
   
   # Veritabanı
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=postgres
   DB_NAME=haberdb
   DB_SSL_MODE=disable
   
   # Redis
   REDIS_HOST=localhost
   REDIS_PORT=6379
   REDIS_PASSWORD=
   REDIS_DB=0
   
   # JWT
   JWT_SECRET=change-this-in-production
   JWT_ACCESS_TOKEN_EXP=60
   JWT_REFRESH_TOKEN_EXP=168
   ```

4. Veritabanı oluşturulur:
   ```bash
   psql -U postgres -c "CREATE DATABASE haberdb;"
   ```

5. Uygulamayı çalıştırın:
   ```bash
   go run cmd/server/main.go
   ```

## Proje Yapısı

```
haber/
│
├── cmd/
│   └── server/             # Ana uygulama başlangıç noktası
│
├── internal/
│   ├── config/             # Konfigürasyon yapıları
│   ├── middleware/         # Middleware fonksiyonları
│   ├── models/             # Veritabanı modelleri
│   ├── handlers/           # HTTP işleyicileri
│   ├── repository/         # Veritabanı erişim katmanı
│   ├── services/           # İş mantığı
│   └── utils/              # Yardımcı fonksiyonlar
│
├── pkg/
│   ├── auth/               # Kimlik doğrulama
│   ├── cache/              # Önbellek
│   ├── validator/          # Veri doğrulama
│   └── logger/             # Loglama
│
├── migrations/             # Veritabanı migrasyon dosyaları
│
├── web/
│   ├── templates/          # HTML şablonları
│   ├── static/             # Statik dosyalar (CSS, JS, resimler)
│   └── uploads/            # Yüklenen dosyalar
│
├── docker/                 # Docker yapılandırması
└── scripts/                # Yardımcı scriptler
```

## Admin Paneli

Admin paneline `http://localhost:3000/admin` adresinden erişilebilir.

Varsayılan kullanıcı:
- Kullanıcı Adı: `admin`
- Şifre: `admin123`

İlk girişten sonra şifreyi değiştirmeniz önerilir.

## Docker ile Çalıştırma

```bash
docker-compose up -d
```

## Katkıda Bulunma

1. Projeyi fork edin
2. Yeni bir özellik branchi oluşturun (`git checkout -b yeni-ozellik`)
3. Değişikliklerinizi commit edin (`git commit -am 'Yeni özellik: xyz'`)
4. Branch'i push edin (`git push origin yeni-ozellik`)
5. Pull Request oluşturun

## Lisans

Bu proje MIT lisansı altında lisanslanmıştır. Daha fazla bilgi için `LICENSE` dosyasına bakın. 