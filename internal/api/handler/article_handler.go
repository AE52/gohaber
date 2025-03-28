package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/username/haber/internal/domain"
	"github.com/username/haber/internal/service"
)

// ArticleHandler makale işleyicileri
type ArticleHandler struct {
	articleService service.IArticleService
}

// NewArticleHandler yeni bir ArticleHandler oluşturur
func NewArticleHandler(articleService service.IArticleService) *ArticleHandler {
	return &ArticleHandler{
		articleService: articleService,
	}
}

// RegisterRoutes rotaları kayıt eder
func (h *ArticleHandler) RegisterRoutes(router fiber.Router, authMw fiber.Handler, adminMw fiber.Handler) {
	// Herkese açık rotalar
	router.Get("/articles", h.ListArticles)
	router.Get("/articles/:id", h.GetArticle)
	router.Get("/articles/slug/:slug", h.GetArticleBySlug)
	router.Get("/articles/featured", h.GetFeaturedArticles)
	router.Get("/articles/category/:categoryID", h.GetArticlesByCategory)
	router.Get("/articles/tag/:tagID", h.GetArticlesByTag)
	router.Get("/articles/author/:authorID", h.GetArticlesByAuthor)

	// Admin ve editörler için rotalar
	adminRoutes := router.Group("/admin/articles", adminMw)
	adminRoutes.Post("/", h.CreateArticle)
	adminRoutes.Put("/:id", h.UpdateArticle)
	adminRoutes.Delete("/:id", h.DeleteArticle)
}

// ListArticles makaleleri listeler
// @Summary Makaleleri listele
// @Description Tüm makaleleri sayfalanmış şekilde listeler
// @Tags Makaleler
// @Accept json
// @Produce json
// @Param page query int false "Sayfa numarası (varsayılan: 1)"
// @Param limit query int false "Sayfa başına sonuç sayısı (varsayılan: 10, maksimum: 100)"
// @Param status query string false "Makale durumu (varsayılan: published)"
// @Success 200 {object} domain.PaginatedResponse{data=[]domain.Article}
// @Router /articles [get]
func (h *ArticleHandler) ListArticles(c *fiber.Ctx) error {
	// Sayfalama parametrelerini al
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	offset := (page - 1) * limit

	// Filtreleri al
	filters := make(map[string]interface{})
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	} else {
		filters["status"] = domain.ArticleStatusPublished
	}

	// Makaleleri getir
	articles, total, err := h.articleService.ListArticles(offset, limit, filters)
	if err != nil {
		return err
	}

	// Toplam sayfa sayısını hesapla
	totalPages := (int(total) + limit - 1) / limit
	if totalPages < 1 {
		totalPages = 1
	}

	// Sonuçları döndür
	return c.JSON(fiber.Map{
		"data": articles,
		"meta": fiber.Map{
			"current_page": page,
			"per_page":     limit,
			"total":        total,
			"total_pages":  totalPages,
		},
	})
}

// GetArticle ID'ye göre makale getirir
// @Summary Makale detayı getir
// @Description ID'ye göre makale detayı getirir
// @Tags Makaleler
// @Accept json
// @Produce json
// @Param id path int true "Makale ID"
// @Success 200 {object} domain.Article
// @Failure 400 {object} domain.ErrorResponse "Geçersiz makale ID"
// @Failure 404 {object} domain.ErrorResponse "Makale bulunamadı"
// @Router /articles/{id} [get]
func (h *ArticleHandler) GetArticle(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Geçersiz makale ID")
	}

	article, err := h.articleService.GetArticleByID(uint(id))
	if err != nil {
		return err
	}

	return c.JSON(article)
}

// GetArticleBySlug slug'a göre makale getirir
// @Summary Slug ile makale getir
// @Description Slug'a göre makale detayı getirir
// @Tags Makaleler
// @Accept json
// @Produce json
// @Param slug path string true "Makale Slug"
// @Success 200 {object} domain.Article
// @Failure 404 {object} domain.ErrorResponse "Makale bulunamadı"
// @Router /articles/slug/{slug} [get]
func (h *ArticleHandler) GetArticleBySlug(c *fiber.Ctx) error {
	slug := c.Params("slug")
	article, err := h.articleService.GetArticleBySlug(slug)
	if err != nil {
		return err
	}

	return c.JSON(article)
}

// GetFeaturedArticles öne çıkan makaleleri getirir
// @Summary Öne çıkan makaleleri getir
// @Description Öne çıkan makaleleri listeler
// @Tags Makaleler
// @Accept json
// @Produce json
// @Param limit query int false "Maksimum makale sayısı (varsayılan: 5, maksimum: 20)"
// @Success 200 {array} domain.Article
// @Router /articles/featured [get]
func (h *ArticleHandler) GetFeaturedArticles(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "5"))
	if limit < 1 || limit > 20 {
		limit = 5
	}

	articles, err := h.articleService.GetFeaturedArticles(limit)
	if err != nil {
		return err
	}

	return c.JSON(articles)
}

// GetArticlesByCategory kategoriye göre makaleleri getirir
// @Summary Kategoriye göre makaleleri getir
// @Description Belirli bir kategoriye ait makaleleri sayfalanmış şekilde listeler
// @Tags Makaleler
// @Accept json
// @Produce json
// @Param categoryID path int true "Kategori ID"
// @Param page query int false "Sayfa numarası (varsayılan: 1)"
// @Param limit query int false "Sayfa başına sonuç sayısı (varsayılan: 10, maksimum: 100)"
// @Success 200 {object} domain.PaginatedResponse{data=[]domain.Article}
// @Failure 400 {object} domain.ErrorResponse "Geçersiz kategori ID"
// @Router /articles/category/{categoryID} [get]
func (h *ArticleHandler) GetArticlesByCategory(c *fiber.Ctx) error {
	categoryID, err := strconv.ParseUint(c.Params("categoryID"), 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Geçersiz kategori ID")
	}

	// Sayfalama parametrelerini al
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	offset := (page - 1) * limit

	articles, total, err := h.articleService.GetArticlesByCategory(uint(categoryID), offset, limit)
	if err != nil {
		return err
	}

	// Toplam sayfa sayısını hesapla
	totalPages := (int(total) + limit - 1) / limit
	if totalPages < 1 {
		totalPages = 1
	}

	// Sonuçları döndür
	return c.JSON(fiber.Map{
		"data": articles,
		"meta": fiber.Map{
			"current_page": page,
			"per_page":     limit,
			"total":        total,
			"total_pages":  totalPages,
		},
	})
}

// GetArticlesByTag etikete göre makaleleri getirir
// @Summary Etikete göre makaleleri getir
// @Description Belirli bir etikete sahip makaleleri sayfalanmış şekilde listeler
// @Tags Makaleler
// @Accept json
// @Produce json
// @Param tagID path int true "Etiket ID"
// @Param page query int false "Sayfa numarası (varsayılan: 1)"
// @Param limit query int false "Sayfa başına sonuç sayısı (varsayılan: 10, maksimum: 100)"
// @Success 200 {object} domain.PaginatedResponse{data=[]domain.Article}
// @Failure 400 {object} domain.ErrorResponse "Geçersiz etiket ID"
// @Router /articles/tag/{tagID} [get]
func (h *ArticleHandler) GetArticlesByTag(c *fiber.Ctx) error {
	tagID, err := strconv.ParseUint(c.Params("tagID"), 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Geçersiz etiket ID")
	}

	// Sayfalama parametrelerini al
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	offset := (page - 1) * limit

	articles, total, err := h.articleService.GetArticlesByTag(uint(tagID), offset, limit)
	if err != nil {
		return err
	}

	// Toplam sayfa sayısını hesapla
	totalPages := (int(total) + limit - 1) / limit
	if totalPages < 1 {
		totalPages = 1
	}

	// Sonuçları döndür
	return c.JSON(fiber.Map{
		"data": articles,
		"meta": fiber.Map{
			"current_page": page,
			"per_page":     limit,
			"total":        total,
			"total_pages":  totalPages,
		},
	})
}

// GetArticlesByAuthor yazara göre makaleleri getirir
// @Summary Yazara göre makaleleri getir
// @Description Belirli bir yazara ait makaleleri sayfalanmış şekilde listeler
// @Tags Makaleler
// @Accept json
// @Produce json
// @Param authorID path int true "Yazar ID"
// @Param page query int false "Sayfa numarası (varsayılan: 1)"
// @Param limit query int false "Sayfa başına sonuç sayısı (varsayılan: 10, maksimum: 100)"
// @Success 200 {object} domain.PaginatedResponse{data=[]domain.Article}
// @Failure 400 {object} domain.ErrorResponse "Geçersiz yazar ID"
// @Router /articles/author/{authorID} [get]
func (h *ArticleHandler) GetArticlesByAuthor(c *fiber.Ctx) error {
	authorID, err := strconv.ParseUint(c.Params("authorID"), 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Geçersiz yazar ID")
	}

	// Sayfalama parametrelerini al
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	offset := (page - 1) * limit

	articles, total, err := h.articleService.GetArticlesByAuthor(uint(authorID), offset, limit)
	if err != nil {
		return err
	}

	// Toplam sayfa sayısını hesapla
	totalPages := (int(total) + limit - 1) / limit
	if totalPages < 1 {
		totalPages = 1
	}

	// Sonuçları döndür
	return c.JSON(fiber.Map{
		"data": articles,
		"meta": fiber.Map{
			"current_page": page,
			"per_page":     limit,
			"total":        total,
			"total_pages":  totalPages,
		},
	})
}

// CreateArticle yeni bir makale oluşturur
// @Summary Yeni makale oluştur
// @Description Yeni bir makale oluşturur (Sadece admin ve editörler)
// @Tags Admin, Makaleler
// @Accept json
// @Produce json
// @Param article body domain.CreateArticleRequest true "Makale bilgileri"
// @Success 201 {object} domain.Article
// @Failure 400 {object} domain.ErrorResponse "Geçersiz istek"
// @Failure 401 {object} domain.ErrorResponse "Yetkisiz erişim"
// @Failure 403 {object} domain.ErrorResponse "Yetersiz yetki"
// @Security ApiKeyAuth
// @Router /admin/articles [post]
func (h *ArticleHandler) CreateArticle(c *fiber.Ctx) error {
	var req domain.CreateArticleRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Geçersiz istek formatı")
	}

	// Kullanıcı ID'sini al
	userID := c.Locals("user_id").(uint)

	// Makaleyi oluştur
	article, err := h.articleService.CreateArticle(&req, userID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(article)
}

// UpdateArticle makaleyi günceller
// @Summary Makale güncelle
// @Description Mevcut bir makaleyi günceller (Sadece admin ve editörler)
// @Tags Admin, Makaleler
// @Accept json
// @Produce json
// @Param id path int true "Makale ID"
// @Param article body domain.UpdateArticleRequest true "Makale bilgileri"
// @Success 200 {object} domain.Article
// @Failure 400 {object} domain.ErrorResponse "Geçersiz istek"
// @Failure 401 {object} domain.ErrorResponse "Yetkisiz erişim"
// @Failure 403 {object} domain.ErrorResponse "Yetersiz yetki"
// @Failure 404 {object} domain.ErrorResponse "Makale bulunamadı"
// @Security ApiKeyAuth
// @Router /admin/articles/{id} [put]
func (h *ArticleHandler) UpdateArticle(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Geçersiz makale ID")
	}

	var req domain.UpdateArticleRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Geçersiz istek formatı")
	}

	// Makaleyi güncelle
	article, err := h.articleService.UpdateArticle(uint(id), &req)
	if err != nil {
		return err
	}

	return c.JSON(article)
}

// DeleteArticle makaleyi siler
// @Summary Makale sil
// @Description Bir makaleyi siler (Sadece admin ve editörler)
// @Tags Admin, Makaleler
// @Accept json
// @Produce json
// @Param id path int true "Makale ID"
// @Success 204 "Başarıyla silindi"
// @Failure 400 {object} domain.ErrorResponse "Geçersiz istek"
// @Failure 401 {object} domain.ErrorResponse "Yetkisiz erişim"
// @Failure 403 {object} domain.ErrorResponse "Yetersiz yetki"
// @Failure 404 {object} domain.ErrorResponse "Makale bulunamadı"
// @Security ApiKeyAuth
// @Router /admin/articles/{id} [delete]
func (h *ArticleHandler) DeleteArticle(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Geçersiz makale ID")
	}

	err = h.articleService.DeleteArticle(uint(id))
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}
