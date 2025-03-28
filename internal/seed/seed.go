package seed

import (
	"fmt"
	"log"
	"time"

	"github.com/username/haber/internal/domain"
	"github.com/username/haber/pkg/auth"
	"gorm.io/gorm"
)

// SeedDB veritabanını örnek verilerle doldurur
func SeedDB(db *gorm.DB) error {
	// Veritabanını kontrol et, boş değilse seed işlemini atla
	var userCount int64
	db.Model(&domain.User{}).Count(&userCount)
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
func seedUsers(db *gorm.DB) ([]domain.User, error) {
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

	users := []domain.User{
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
func seedCategories(db *gorm.DB) ([]domain.Category, error) {
	categories := []domain.Category{
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
func seedTags(db *gorm.DB) ([]domain.Tag, error) {
	tags := []domain.Tag{
		{
			Name:      "Politika",
			Slug:      "politika",
			CreatedAt: time.Now(),
		},
		{
			Name:      "Ekonomi",
			Slug:      "ekonomi",
			CreatedAt: time.Now(),
		},
		{
			Name:      "Spor",
			Slug:      "spor",
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
			Name:      "Eğitim",
			Slug:      "egitim",
			CreatedAt: time.Now(),
		},
		{
			Name:      "Dünya",
			Slug:      "dunya",
			CreatedAt: time.Now(),
		},
		{
			Name:      "Yaşam",
			Slug:      "yasam",
			CreatedAt: time.Now(),
		},
		{
			Name:      "Analiz",
			Slug:      "analiz",
			CreatedAt: time.Now(),
		},
		{
			Name:      "Röportaj",
			Slug:      "roportaj",
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
func seedArticles(db *gorm.DB, users []domain.User, categories []domain.Category, tags []domain.Tag) ([]domain.Article, error) {
	// Yardımcı fonksiyon
	ptrTime := func(t time.Time) *time.Time {
		return &t
	}

	// Makaleleri oluştur
	// Not: Bu sadece iki örnek makale içeriyor, gerçek uygulamada daha fazla ve kategorilere göre çeşitli makaleler eklenebilir
	articles := []domain.Article{
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
	}

	// Makaleleri veritabanına ekle
	if err := db.Create(&articles).Error; err != nil {
		return nil, err
	}

	// Makalelere etiketleri bağla
	for i := range articles {
		// Her makaleye 3 etiket ekleyelim (makale indeksine göre farklı etiketler)
		var articleTags []domain.Tag
		for j := 0; j < 3; j++ {
			tagIndex := (i + j) % len(tags)
			articleTags = append(articleTags, tags[tagIndex])
		}

		// Many-to-many ilişkiyi kur
		if err := db.Model(&articles[i]).Association("Tags").Append(articleTags); err != nil {
			return nil, err
		}
	}

	log.Println("Makale seed verileri oluşturuldu")
	return articles, nil
}

// seedComments yorum örnek verilerini oluşturur
func seedComments(db *gorm.DB, users []domain.User, articles []domain.Article) error {
	comments := []domain.Comment{
		{
			ArticleID:  articles[0].ID,
			UserID:     users[2].ID,
			Content:    "Bu haberi okuyunca çok üzüldüm. Çocukların koşulları iyileştirilmeli.",
			IsApproved: true,
			CreatedAt:  time.Now().Add(-12 * time.Hour),
			UpdatedAt:  time.Now().Add(-12 * time.Hour),
		},
		{
			ArticleID:  articles[0].ID,
			UserID:     users[1].ID,
			Content:    "Konu hakkında daha detaylı bilgi verilmeli.",
			IsApproved: true,
			CreatedAt:  time.Now().Add(-8 * time.Hour),
			UpdatedAt:  time.Now().Add(-8 * time.Hour),
		},
		{
			ArticleID:  articles[1].ID,
			UserID:     users[2].ID,
			Content:    "Bu vergi düzenlemesi orta sınıfı nasıl etkileyecek?",
			IsApproved: true,
			CreatedAt:  time.Now().Add(-6 * time.Hour),
			UpdatedAt:  time.Now().Add(-6 * time.Hour),
		},
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
	adSpaces := []domain.AdSpace{
		{
			Name:      "Üst Banner",
			Placement: "header",
			Content:   "<div class=\"ad-banner\">Banner Reklam Alanı</div>",
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Name:      "Yan Sidebar",
			Placement: "sidebar",
			Content:   "<div class=\"ad-sidebar\">Sidebar Reklam Alanı</div>",
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Name:      "Makale İçi",
			Placement: "article-content",
			Content:   "<div class=\"ad-content\">İçerik Reklam Alanı</div>",
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Name:      "Alt Alan",
			Placement: "footer",
			Content:   "<div class=\"ad-footer\">Alt Reklam Alanı</div>",
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
