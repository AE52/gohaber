package seed

import (
	"fmt"
	"log"
	"time"

	"github.com/username/haber/internal/models"
	"github.com/username/haber/pkg/auth"
	"gorm.io/gorm"
)

// SeedDB veritabanını örnek verilerle doldurur
func SeedDB(db *gorm.DB) error {
	// Veritabanını kontrol et, boş değilse seed işlemini atla
	var userCount int64
	db.Model(&models.User{}).Count(&userCount)
	if userCount > 0 {
		fmt.Println("Veritabanı zaten dolu, seed işlemi atlanıyor.")
		return nil
	}

	// Seed işlemini transaction içinde yap
	return db.Transaction(func(tx *gorm.DB) error {
		// 1. Kullanıcılar
		users, err := seedUsers(tx)
		if err != nil {
			return err
		}

		// 2. Kategoriler
		categories, err := seedCategories(tx)
		if err != nil {
			return err
		}

		// 3. Etiketler
		tags, err := seedTags(tx)
		if err != nil {
			return err
		}

		// 4. Makaleler
		articles, err := seedArticles(tx, users, categories, tags)
		if err != nil {
			return err
		}

		// 5. Yorumlar
		if err := seedComments(tx, users, articles); err != nil {
			return err
		}

		// 6. Reklam Alanları
		if err := seedAdSpaces(tx); err != nil {
			return err
		}

		fmt.Println("Tüm seed verileri başarıyla oluşturuldu.")
		return nil
	})
}

// seedUsers kullanıcı örnek verilerini oluşturur
func seedUsers(db *gorm.DB) ([]models.User, error) {
	// Şifreleri hashleme
	adminPass, err := auth.HashPassword("admin123")
	if err != nil {
		return nil, err
	}

	editorPass, err := auth.HashPassword("editor123")
	if err != nil {
		return nil, err
	}

	userPass, err := auth.HashPassword("user123")
	if err != nil {
		return nil, err
	}

	users := []models.User{
		{
			Username:     "admin",
			Email:        "admin@haber.com",
			PasswordHash: adminPass,
			FullName:     "Admin Kullanıcı",
			Role:         "admin",
			ProfileImage: "/uploads/users/admin.jpg",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			Username:     "editor",
			Email:        "editor@haber.com",
			PasswordHash: editorPass,
			FullName:     "Editör Kullanıcı",
			Role:         "editor",
			ProfileImage: "/uploads/users/editor.jpg",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			Username:     "user",
			Email:        "user@haber.com",
			PasswordHash: userPass,
			FullName:     "Normal Kullanıcı",
			Role:         "user",
			ProfileImage: "/uploads/users/user.jpg",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
	}

	// Kullanıcıları veritabanına ekle
	if err := db.Create(&users).Error; err != nil {
		return nil, err
	}

	log.Println("Kullanıcı seed verileri oluşturuldu")
	return users, nil
}

// seedCategories kategori örnek verilerini oluşturur
func seedCategories(db *gorm.DB) ([]models.Category, error) {
	categories := []models.Category{
		{
			Name:        "Gündem",
			Slug:        "gundem",
			Description: "Güncel haberler ve gelişmeler",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Name:        "Ekonomi",
			Slug:        "ekonomi",
			Description: "Ekonomi haberleri ve piyasa gelişmeleri",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Name:        "Spor",
			Slug:        "spor",
			Description: "Spor haberleri ve gelişmeleri",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Name:        "Dünya",
			Slug:        "dunya",
			Description: "Dünya haberleri ve gelişmeleri",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Name:        "Teknoloji",
			Slug:        "teknoloji",
			Description: "Teknoloji haberleri ve gelişmeler",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Name:        "Sağlık",
			Slug:        "saglik",
			Description: "Sağlık haberleri ve bilgileri",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Name:        "Kültür-Sanat",
			Slug:        "kultur-sanat",
			Description: "Kültür ve sanat haberleri",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	// Kategorileri veritabanına ekle
	if err := db.Create(&categories).Error; err != nil {
		return nil, err
	}

	log.Println("Kategori seed verileri oluşturuldu")
	return categories, nil
}

// seedTags etiket örnek verilerini oluşturur
func seedTags(db *gorm.DB) ([]models.Tag, error) {
	tags := []models.Tag{
		{
			Name:      "Türkiye",
			Slug:      "turkiye",
			CreatedAt: time.Now(),
		},
		{
			Name:      "Dünya",
			Slug:      "dunya",
			CreatedAt: time.Now(),
		},
		{
			Name:      "Ekonomi",
			Slug:      "ekonomi",
			CreatedAt: time.Now(),
		},
		{
			Name:      "Siyaset",
			Slug:      "siyaset",
			CreatedAt: time.Now(),
		},
		{
			Name:      "Futbol",
			Slug:      "futbol",
			CreatedAt: time.Now(),
		},
		{
			Name:      "Basketbol",
			Slug:      "basketbol",
			CreatedAt: time.Now(),
		},
		{
			Name:      "Sağlık",
			Slug:      "saglik",
			CreatedAt: time.Now(),
		},
		{
			Name:      "Teknoloji",
			Slug:      "teknoloji",
			CreatedAt: time.Now(),
		},
		{
			Name:      "Bilim",
			Slug:      "bilim",
			CreatedAt: time.Now(),
		},
		{
			Name:      "Eğitim",
			Slug:      "egitim",
			CreatedAt: time.Now(),
		},
		{
			Name:      "Kültür",
			Slug:      "kultur",
			CreatedAt: time.Now(),
		},
		{
			Name:      "Sanat",
			Slug:      "sanat",
			CreatedAt: time.Now(),
		},
	}

	// Etiketleri veritabanına ekle
	if err := db.Create(&tags).Error; err != nil {
		return nil, err
	}

	log.Println("Etiket seed verileri oluşturuldu")
	return tags, nil
}

// seedArticles makale örnek verilerini oluşturur
func seedArticles(db *gorm.DB, users []models.User, categories []models.Category, tags []models.Tag) ([]models.Article, error) {
	// Gündem haberleri
	gundemArticles := []models.Article{
		{
			Title:         "İstanbul'da Su Kesintisi! İşte Etkilenecek İlçeler",
			Slug:          "istanbulda-su-kesintisi-iste-etkilenecek-ilceler",
			Content:       "İstanbul'da planlı bakım çalışmaları nedeniyle yarın 6 ilçede su kesintisi yapılacak. İstanbul Su ve Kanalizasyon İdaresi (İSKİ) tarafından yapılan açıklamada...",
			Summary:       "İstanbul'da planlı bakım çalışmaları nedeniyle yarın 6 ilçede su kesintisi yapılacak.",
			FeaturedImage: "https://images.unsplash.com/photo-1563543054715-4f60d55bccb3?ixlib=rb-4.0.3&auto=format&fit=crop&w=800&q=80",
			AuthorID:      users[0].ID,
			CategoryID:    categories[0].ID,
			Status:        "published",
			ViewCount:     120,
			IsFeatured:    true,
			PublishedAt:   ptrTime(time.Now().Add(-48 * time.Hour)),
			CreatedAt:     time.Now().Add(-72 * time.Hour),
			UpdatedAt:     time.Now(),
		},
		{
			Title:         "Yeni Vergi Düzenlemesi Meclisten Geçti",
			Slug:          "yeni-vergi-duzenlemesi-meclisten-gecti",
			Content:       "Uzun süredir tartışılan vergi düzenlemesi meclisten geçti. Yeni düzenlemeyle birlikte, kurumlar vergisi oranları yeniden belirlendi ve...",
			Summary:       "Uzun süredir tartışılan vergi düzenlemesi meclisten geçti.",
			FeaturedImage: "https://images.unsplash.com/photo-1554224155-8d04cb21cd6c?ixlib=rb-4.0.3&auto=format&fit=crop&w=800&q=80",
			AuthorID:      users[1].ID,
			CategoryID:    categories[0].ID,
			Status:        "published",
			ViewCount:     95,
			IsFeatured:    false,
			PublishedAt:   ptrTime(time.Now().Add(-24 * time.Hour)),
			CreatedAt:     time.Now().Add(-36 * time.Hour),
			UpdatedAt:     time.Now(),
		},
		{
			Title:         "Meteoroloji'den Kuvvetli Yağış Uyarısı",
			Slug:          "meteorolojiden-kuvvetli-yagis-uyarisi",
			Content:       "Meteoroloji Genel Müdürlüğü tarafından yapılan son tahminlere göre, önümüzdeki hafta ülke genelinde kuvvetli yağış beklenirken, bazı bölgelerde sel riski bulunuyor...",
			Summary:       "Meteoroloji Genel Müdürlüğü tarafından yapılan son tahminlere göre, önümüzdeki hafta ülke genelinde kuvvetli yağış bekleniyor.",
			FeaturedImage: "https://images.unsplash.com/photo-1514632595-4944383f2737?ixlib=rb-4.0.3&auto=format&fit=crop&w=800&q=80",
			AuthorID:      users[1].ID,
			CategoryID:    categories[0].ID,
			Status:        "published",
			ViewCount:     85,
			IsFeatured:    false,
			PublishedAt:   ptrTime(time.Now().Add(-12 * time.Hour)),
			CreatedAt:     time.Now().Add(-24 * time.Hour),
			UpdatedAt:     time.Now(),
		},
	}

	// Ekonomi haberleri
	ekonomiArticles := []models.Article{
		{
			Title:         "Dolar Rekor Kırdı! İşte Son Durum",
			Slug:          "dolar-rekor-kirdi-iste-son-durum",
			Content:       "Dolar kuru bugün uluslararası piyasalardaki gelişmeler ve yurt içindeki ekonomik belirsizliklerin etkisiyle rekor seviyeyi gördü. Dolar kuru sabah saatlerinde...",
			Summary:       "Dolar kuru bugün rekor seviyeyi görerek tarihi zirveyi gördü.",
			FeaturedImage: "https://images.unsplash.com/photo-1591696205602-2f950c417cb9?ixlib=rb-4.0.3&auto=format&fit=crop&w=800&q=80",
			AuthorID:      users[0].ID,
			CategoryID:    categories[1].ID,
			Status:        "published",
			ViewCount:     230,
			IsFeatured:    true,
			PublishedAt:   ptrTime(time.Now().Add(-5 * time.Hour)),
			CreatedAt:     time.Now().Add(-8 * time.Hour),
			UpdatedAt:     time.Now(),
		},
		{
			Title:         "Merkez Bankası Faiz Kararını Açıkladı",
			Slug:          "merkez-bankasi-faiz-kararini-acikladi",
			Content:       "Merkez Bankası, bugün gerçekleştirdiği Para Politikası Kurulu toplantısında faiz oranlarıyla ilgili kararını açıkladı. Yapılan açıklamaya göre faiz oranları...",
			Summary:       "Merkez Bankası, bugün gerçekleştirdiği Para Politikası Kurulu toplantısında faiz oranlarıyla ilgili kararını açıkladı.",
			FeaturedImage: "https://images.unsplash.com/photo-1553729459-efe14ef6055d?ixlib=rb-4.0.3&auto=format&fit=crop&w=800&q=80",
			AuthorID:      users[1].ID,
			CategoryID:    categories[1].ID,
			Status:        "published",
			ViewCount:     180,
			IsFeatured:    false,
			PublishedAt:   ptrTime(time.Now().Add(-32 * time.Hour)),
			CreatedAt:     time.Now().Add(-40 * time.Hour),
			UpdatedAt:     time.Now(),
		},
		{
			Title:         "Asgari Ücret Zammı İçin Görüşmeler Başladı",
			Slug:          "asgari-ucret-zammi-icin-gorusmeler-basladi",
			Content:       "Asgari ücret zammı için Çalışma Bakanlığı, işçi ve işveren sendikaları arasındaki görüşmeler bugün başladı. Toplantıda, enflasyon oranları ve ekonomik göstergeler eşliğinde...",
			Summary:       "Asgari ücret zammı için Çalışma Bakanlığı, işçi ve işveren sendikaları arasındaki görüşmeler bugün başladı.",
			FeaturedImage: "https://images.unsplash.com/photo-1607863680198-23d4b2565df0?ixlib=rb-4.0.3&auto=format&fit=crop&w=800&q=80",
			AuthorID:      users[0].ID,
			CategoryID:    categories[1].ID,
			Status:        "published",
			ViewCount:     160,
			IsFeatured:    false,
			PublishedAt:   ptrTime(time.Now().Add(-18 * time.Hour)),
			CreatedAt:     time.Now().Add(-24 * time.Hour),
			UpdatedAt:     time.Now(),
		},
	}

	// Spor haberleri
	sporArticles := []models.Article{
		{
			Title:         "Galatasaray, Fenerbahçe Derbisine Hazır",
			Slug:          "galatasaray-fenerbahce-derbisine-hazir",
			Content:       "Süper Lig'in 30. haftasında Fenerbahçe ile karşılaşacak olan Galatasaray, hazırlıklarını tamamladı. Teknik direktör Okan Buruk yönetiminde gerçekleştirilen son antrenmanın ardından...",
			Summary:       "Süper Lig'in 30. haftasında Fenerbahçe ile karşılaşacak olan Galatasaray, hazırlıklarını tamamladı.",
			FeaturedImage: "https://images.unsplash.com/photo-1508098682722-e99c643e7f3b?ixlib=rb-4.0.3&auto=format&fit=crop&w=800&q=80",
			AuthorID:      users[1].ID,
			CategoryID:    categories[2].ID,
			Status:        "published",
			ViewCount:     320,
			IsFeatured:    true,
			PublishedAt:   ptrTime(time.Now().Add(-6 * time.Hour)),
			CreatedAt:     time.Now().Add(-10 * time.Hour),
			UpdatedAt:     time.Now(),
		},
		{
			Title:         "Milli Basketbolcu NBA'de Sözleşme İmzaladı",
			Slug:          "milli-basketbolcu-nbade-sozlesme-imzaladi",
			Content:       "Türk basketbolunun yükselen yıldızı, NBA takımlarından Los Angeles Lakers ile 3 yıllık sözleşme imzaladı. Yapılan açıklamaya göre, transferin değeri...",
			Summary:       "Türk basketbolunun yükselen yıldızı, NBA takımlarından Los Angeles Lakers ile 3 yıllık sözleşme imzaladı.",
			FeaturedImage: "https://images.unsplash.com/photo-1518407613690-d9fc990e795f?ixlib=rb-4.0.3&auto=format&fit=crop&w=800&q=80",
			AuthorID:      users[1].ID,
			CategoryID:    categories[2].ID,
			Status:        "published",
			ViewCount:     210,
			IsFeatured:    false,
			PublishedAt:   ptrTime(time.Now().Add(-28 * time.Hour)),
			CreatedAt:     time.Now().Add(-36 * time.Hour),
			UpdatedAt:     time.Now(),
		},
		{
			Title:         "Türkiye, Dünya Kupası Elemelerinde Zafer Kazandı",
			Slug:          "turkiye-dunya-kupasi-elemelerinde-zafer-kazandi",
			Content:       "A Milli Futbol Takımımız, 2026 Dünya Kupası Elemeleri kapsamında konuk ettiği takımı 3-1 mağlup etti. Milli takımımız, bu galibiyetle puanını 10'a yükseltti ve grup liderliğine yerleşti...",
			Summary:       "A Milli Futbol Takımımız, 2026 Dünya Kupası Elemeleri kapsamında konuk ettiği takımı 3-1 mağlup etti.",
			FeaturedImage: "https://images.unsplash.com/photo-1522778119026-d647f0596c20?ixlib=rb-4.0.3&auto=format&fit=crop&w=800&q=80",
			AuthorID:      users[0].ID,
			CategoryID:    categories[2].ID,
			Status:        "published",
			ViewCount:     185,
			IsFeatured:    false,
			PublishedAt:   ptrTime(time.Now().Add(-16 * time.Hour)),
			CreatedAt:     time.Now().Add(-20 * time.Hour),
			UpdatedAt:     time.Now(),
		},
	}

	// Tüm makaleleri birleştir
	articles := append(gundemArticles, ekonomiArticles...)
	articles = append(articles, sporArticles...)

	// Veritabanına ekle
	if err := db.Create(&articles).Error; err != nil {
		return nil, err
	}

	// Makalelere etiket ilişkisi ekle
	for i, article := range articles {
		var articleTags []models.Tag
		// Her makaleye farklı etiketler ekle
		tagIndices := []int{i % len(tags), (i + 3) % len(tags), (i + 5) % len(tags)}
		for _, idx := range tagIndices {
			articleTags = append(articleTags, tags[idx])
		}

		// Many-to-many ilişkisini güncelle
		if err := db.Model(&article).Association("Tags").Replace(articleTags); err != nil {
			return nil, err
		}
	}

	log.Println("Makale seed verileri oluşturuldu")
	return articles, nil
}

// seedComments yorum örnek verilerini oluşturur
func seedComments(db *gorm.DB, users []models.User, articles []models.Article) error {
	var comments []models.Comment

	for _, article := range articles {
		// Her makale için 2-5 arası yorum ekle
		commentCount := 2 + (int(article.ID) % 4) // 2 ile 5 arasında değişen sayıda yorum
		for i := 0; i < commentCount; i++ {
			// Farklı kullanıcılara yorum yaptır
			userIndex := i % len(users)
			comment := models.Comment{
				ArticleID:  article.ID,
				UserID:     users[userIndex].ID,
				Content:    fmt.Sprintf("Bu haber hakkında düşüncelerim: Lorem ipsum dolor sit amet, consectetur adipiscing elit. Yorum %d", i+1),
				IsApproved: true,
				CreatedAt:  time.Now().Add(-time.Duration(i+1) * 12 * time.Hour),
				UpdatedAt:  time.Now(),
			}
			comments = append(comments, comment)
		}
	}

	// Yorumları veritabanına ekle
	if err := db.Create(&comments).Error; err != nil {
		return err
	}

	log.Println("Yorum seed verileri oluşturuldu")
	return nil
}

// seedAdSpaces reklam alanı örnek verilerini oluşturur
func seedAdSpaces(db *gorm.DB) error {
	adSpaces := []models.AdSpace{
		{
			Name:      "Header Reklam",
			Placement: "header",
			Content:   "<div class='text-center'><img src='/uploads/ads/header-banner.jpg' alt='Header Reklam' class='img-fluid'></div>",
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Name:      "Sidebar Üst",
			Placement: "sidebar-top",
			Content:   "<div class='text-center'><img src='/uploads/ads/sidebar-top.jpg' alt='Sidebar Üst Reklam' class='img-fluid'></div>",
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Name:      "Sidebar Alt",
			Placement: "sidebar-bottom",
			Content:   "<div class='text-center'><img src='/uploads/ads/sidebar-bottom.jpg' alt='Sidebar Alt Reklam' class='img-fluid'></div>",
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Name:      "İçerik Arası",
			Placement: "content",
			Content:   "<div class='text-center my-4'><img src='/uploads/ads/content.jpg' alt='İçerik Reklam' class='img-fluid'></div>",
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	// Reklam alanlarını veritabanına ekle
	if err := db.Create(&adSpaces).Error; err != nil {
		return err
	}

	log.Println("Reklam alanı seed verileri oluşturuldu")
	return nil
}

// ptrTime time.Time tipindeki değerin pointerını döndürür
func ptrTime(t time.Time) *time.Time {
	return &t
}
