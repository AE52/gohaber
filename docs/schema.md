# Haber Projesi Mimari ve Şema Dokümantasyonu

## İçindekiler
1. [Genel Mimari Yapısı](#genel-mimari-yapısı)
2. [Veritabanı Şeması](#veritabanı-şeması)
3. [API Endpoints](#api-endpoints)
4. [Servisler ve Bileşenler](#servisler-ve-bileşenler)
5. [Dosya Yapısı](#dosya-yapısı)
6. [Medya Yönetimi (MinIO)](#medya-yönetimi-minio)
7. [Güvenlik ve Kimlik Doğrulama](#güvenlik-ve-kimlik-doğrulama)
8. [Dış Bağımlılıklar](#dış-bağımlılıklar)

## Genel Mimari Yapısı

Proje, Clean Architecture prensiplerine dayanmaktadır ve şu katmanlardan oluşur:

1. **Domain Katmanı**: Uygulamanın temel iş mantığını ve veri yapılarını içerir.
2. **Repository Katmanı**: Veritabanı işlemlerini soyutlar.
3. **Service Katmanı**: İş mantığı kurallarını uygular.
4. **API Katmanı**: HTTP isteklerini işler ve yanıtları döndürür.
5. **Infrastructure Katmanı**: MinIO, JWT gibi dış servislerle iletişim kurar.

Bu yapı, bağımlılıkları iç katmanlardan dış katmanlara doğru yönlendirerek, iç katmanların dış katmanlar hakkında hiçbir bilgiye sahip olmamasını sağlar.

```
+-------------------+
|    API Handlers   |
+--------+----------+
         |
         v
+--------+----------+
|     Services      |
+--------+----------+
         |
         v
+--------+----------+
|   Repositories    |
+--------+----------+
         |
         v
+--------+----------+
|      Domain       |
+-------------------+
```

## Veritabanı Şeması

Veritabanı aşağıdaki temel tablolardan oluşur:

### 1. Users
- `id`: Birincil anahtar
- `username`: Kullanıcı adı (benzersiz)
- `email`: E-posta adresi (benzersiz)
- `password_hash`: Şifre hash'i
- `full_name`: Tam ad
- `role`: Kullanıcı rolü (admin, editor, user)
- `profile_image`: Profil resmi yolu
- `created_at`: Oluşturulma tarihi
- `updated_at`: Güncellenme tarihi
- `deleted_at`: Silinme tarihi (soft delete için)

### 2. Articles
- `id`: Birincil anahtar
- `title`: Makale başlığı
- `slug`: SEO dostu URL (benzersiz)
- `content`: Makale içeriği
- `summary`: Özet
- `status`: Durum (published, draft, pending)
- `featured_image`: Öne çıkan resim
- `is_featured`: Öne çıkarılmış mı
- `view_count`: Görüntülenme sayısı
- `author_id`: Yazar ID'si (Users tablosuna referans)
- `category_id`: Kategori ID'si (Categories tablosuna referans)
- `published_at`: Yayınlanma tarihi
- `created_at`: Oluşturulma tarihi
- `updated_at`: Güncellenme tarihi
- `deleted_at`: Silinme tarihi (soft delete için)

### 3. Categories
- `id`: Birincil anahtar
- `name`: Kategori adı
- `slug`: SEO dostu URL (benzersiz)
- `description`: Açıklama
- `parent_id`: Üst kategori ID'si (self-reference)
- `created_at`: Oluşturulma tarihi
- `updated_at`: Güncellenme tarihi
- `deleted_at`: Silinme tarihi (soft delete için)

### 4. Tags
- `id`: Birincil anahtar
- `name`: Etiket adı
- `slug`: SEO dostu URL (benzersiz)
- `created_at`: Oluşturulma tarihi
- `updated_at`: Güncellenme tarihi
- `deleted_at`: Silinme tarihi (soft delete için)

### 5. ArticleTags (İlişki Tablosu)
- `article_id`: Makale ID'si
- `tag_id`: Etiket ID'si
- (Bu iki alan birlikte birincil anahtar oluşturur)

### 6. Comments
- `id`: Birincil anahtar
- `content`: Yorum içeriği
- `article_id`: Makale ID'si (Articles tablosuna referans)
- `user_id`: Kullanıcı ID'si (Users tablosuna referans)
- `parent_id`: Üst yorum ID'si (self-reference, yanıtlar için)
- `is_approved`: Onaylanmış mı
- `created_at`: Oluşturulma tarihi
- `updated_at`: Güncellenme tarihi
- `deleted_at`: Silinme tarihi (soft delete için)

### 7. Media
- `id`: Birincil anahtar
- `filename`: Dosya adı
- `path`: Dosya yolu
- `url`: Erişim URL'i
- `size`: Dosya boyutu
- `type`: Dosya türü
- `mime`: MIME tipi
- `upload_type`: Yükleme türü (image, document, vb.)
- `user_id`: Yükleyen kullanıcı ID'si
- `created_at`: Oluşturulma tarihi
- `updated_at`: Güncellenme tarihi

### 8. Settings
- `id`: Birincil anahtar
- `key`: Ayar anahtarı
- `value`: Ayar değeri
- `type`: Değer tipi
- `created_at`: Oluşturulma tarihi
- `updated_at`: Güncellenme tarihi

## API Endpoints

API, RESTful prensiplerine dayanmaktadır ve aşağıdaki ana endpoint'leri içerir:

### Kimlik Doğrulama
- `POST /auth/register`: Kullanıcı kaydı
- `POST /auth/login`: Kullanıcı girişi
- `POST /auth/refresh`: Token yenileme
- `POST /auth/reset-password`: Şifre sıfırlama isteği
- `POST /auth/reset-password/confirm`: Şifre sıfırlama onayı

### Kullanıcılar
- `GET /users/profile`: Kullanıcı profili görüntüleme
- `PUT /users/profile`: Kullanıcı profili güncelleme
- `PUT /users/password`: Şifre değiştirme
- `GET /users`: Kullanıcı listesi (Admin için)
- `GET /users/{id}`: Kullanıcı detayları (Admin için)
- `POST /users`: Yeni kullanıcı oluşturma (Admin için)
- `PUT /users/{id}`: Kullanıcı güncelleme (Admin için)
- `DELETE /users/{id}`: Kullanıcı silme (Admin için)

### Makaleler
- `GET /articles`: Makale listesi
- `GET /articles/{id}`: Makale detayları
- `GET /articles/slug/{slug}`: Makale detayları (slug ile)
- `GET /articles/category/{categoryID}`: Kategoriye göre makaleler
- `GET /articles/author/{authorID}`: Yazara göre makaleler
- `GET /articles/tag/{tagID}`: Etikete göre makaleler
- `POST /articles`: Yeni makale oluşturma (Editör ve Admin için)
- `PUT /articles/{id}`: Makale güncelleme (Editör ve Admin için)
- `DELETE /articles/{id}`: Makale silme (Editör ve Admin için)

### Kategoriler
- `GET /categories`: Kategori listesi
- `GET /categories/{id}`: Kategori detayları
- `GET /categories/slug/{slug}`: Kategori detayları (slug ile)
- `POST /categories`: Yeni kategori oluşturma (Admin için)
- `PUT /categories/{id}`: Kategori güncelleme (Admin için)
- `DELETE /categories/{id}`: Kategori silme (Admin için)

### Etiketler
- `GET /tags`: Etiket listesi
- `GET /tags/{id}`: Etiket detayları
- `GET /tags/slug/{slug}`: Etiket detayları (slug ile)
- `POST /tags`: Yeni etiket oluşturma (Editör ve Admin için)
- `PUT /tags/{id}`: Etiket güncelleme (Editör ve Admin için)
- `DELETE /tags/{id}`: Etiket silme (Editör ve Admin için)

### Yorumlar
- `GET /comments/article/{articleID}`: Makaleye göre yorumlar
- `POST /comments`: Yeni yorum oluşturma
- `PUT /comments/{id}`: Yorum güncelleme (Yorum sahibi, Editör ve Admin için)
- `DELETE /comments/{id}`: Yorum silme (Yorum sahibi, Editör ve Admin için)

### Medya Yönetimi
- `POST /uploads`: Dosya yükleme
- `GET /uploads`: Dosya listesi
- `GET /uploads/{id}`: Dosya detayları
- `DELETE /uploads/{id}`: Dosya silme

### Ayarlar
- `GET /settings`: Ayarları getir
- `PUT /settings`: Ayarları güncelle (Admin için)

## Servisler ve Bileşenler

### 1. Auth Service
JWT tabanlı kimlik doğrulama hizmetleri sunar:
- Kullanıcı kaydı
- Giriş işlemleri
- Token doğrulama
- Token yenileme
- Şifre sıfırlama

### 2. User Service
Kullanıcı yönetimi işlemlerini yönetir:
- Kullanıcı profili yönetimi
- Şifre değiştirme
- Kullanıcı CRUD işlemleri (admin için)

### 3. Article Service
Makale işlemlerini yönetir:
- Makale listeleme (filtreleme, sıralama, sayfalama)
- Makale detayları
- Makale oluşturma ve güncelleme
- Makalelerin etiketler ve kategorilerle ilişkilendirilmesi

### 4. Category Service
Kategori işlemlerini yönetir:
- Kategori ağacı oluşturma
- Kategori CRUD işlemleri

### 5. Tag Service
Etiket işlemlerini yönetir:
- Etiket CRUD işlemleri
- Etiketlerin makalelerle ilişkilendirilmesi

### 6. Comment Service
Yorum işlemlerini yönetir:
- Yorum listeleme
- Yorum oluşturma, güncelleme ve silme
- Yorum onaylama (moderasyon)

### 7. Media Service
Dosya yükleme ve yönetim işlemlerini yönetir:
- Dosya yükleme (MinIO depolama)
- Dosya listeleme
- Dosya silme

### 8. Settings Service
Sistem ayarlarını yönetir:
- Ayarları getirme
- Ayarları güncelleme
- Önbelleğe alınmış ayarlar

## Dosya Yapısı

```
.
├── cmd
│   └── server              # Ana uygulama giriş noktası
├── docs                    # Swagger dokümanları ve API dokümantasyonu
├── internal                # Uygulama-spesifik paketler
│   ├── api                 # API katmanı
│   │   ├── handler         # HTTP istek işleyicileri
│   │   └── middleware      # HTTP middleware'ler
│   ├── config              # Uygulama yapılandırması
│   ├── domain              # Domain modelleri ve iş mantığı kuralları
│   ├── repository          # Veritabanı işlemlerini soyutlar
│   ├── seed                # Seed verileri
│   └── service             # İş mantığı servisleri
├── migrations              # Veritabanı migration dosyaları
├── pkg                     # Genel amaçlı paketler
│   ├── auth                # Kimlik doğrulama paketi
│   ├── storage             # Depolama paketi (MinIO entegrasyonu)
│   └── swagger             # Swagger entegrasyonu
├── web                     # Web arayüzü dosyaları
│   ├── assets              # Statik dosyalar (CSS, JS, resimler)
│   └── templates           # HTML şablonları
├── .env                    # Ortam değişkenleri
├── Dockerfile              # Docker yapılandırması
├── docker-compose.yml      # Docker Compose yapılandırması
├── go.mod                  # Go modül dosyası
├── go.sum                  # Go modül bağımlılıkları
└── README.md               # Proje açıklaması
```

## Medya Yönetimi (MinIO)

Proje, medya dosyalarını yönetmek için MinIO nesne depolama hizmetini kullanmaktadır:

- **MinIO Service**: Dosya yükleme, silme ve listeleme işlemlerini yönetir.
- **MinIO Konfigürasyonu**: Endpoint, access key, secret key, bucket gibi MinIO yapılandırma parametrelerini içerir.
- **Dosya Yükleme İşlemi**: Dosyalar önce geçici bir konuma yüklenir, doğrulanır ve ardından MinIO'ya aktarılır.
- **URL Oluşturma**: Yüklenen dosyaların erişim URL'leri otomatik olarak oluşturulur.
- **Güvenlik**: Yalnızca kimliği doğrulanmış kullanıcılar dosya yükleyebilir.
- **Dosya Tipleri**: Resimler, belgeler ve diğer izin verilen dosya tipleri için filtreleme mevcuttur.

MinIO hizmeti, docker-compose.yml dosyasında tanımlanmıştır ve uygulamayla birlikte çalışmaktadır. Bu, lokal geliştirme veya demo amaçlı kullanım için uygundur. Üretim ortamında, daha uzun süreli ve güvenilir depolama için ayrı bir MinIO cluster'ı veya AWS S3 gibi bir servis kullanılabilir.

## Güvenlik ve Kimlik Doğrulama

Proje şu güvenlik önlemlerini içerir:

- **JWT Doğrulama**: API'ye erişim için JWT tabanlı kimlik doğrulama
- **Rol Tabanlı Erişim Kontrolü**: Farklı kullanıcı rolleri için farklı yetkilendirme seviyeleri
- **Şifre Hashleme**: Kullanıcı şifreleri güvenli bir şekilde hash'lenir
- **CORS Yapılandırması**: Cross-Origin Resource Sharing güvenlik ayarları
- **Rate Limiting**: API isteklerini sınırlama (istekte bulunulacak)
- **Recovery Middleware**: Panik durumlarının güvenli bir şekilde yakalanması
- **Validasyon**: Gelen isteklerin doğrulanması

## Dış Bağımlılıklar

Proje aşağıdaki ana dış bağımlılıkları kullanmaktadır:

- **Fiber**: Hızlı ve verimli bir HTTP framework'ü
- **GORM**: Go için ORM kütüphanesi
- **JWT-Go**: JWT oluşturma ve doğrulama için
- **Swagger**: API dokümantasyonu
- **MinIO SDK**: MinIO nesne depolama entegrasyonu
- **Golang-Migrate**: Veritabanı migrasyonları için
- **Validator**: Veri doğrulama için 