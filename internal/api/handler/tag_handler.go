package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/username/haber/internal/domain"
	"github.com/username/haber/internal/service"
)

// TagHandler etiket işleyicileri
type TagHandler struct {
	tagService service.ITagService
}

// NewTagHandler yeni bir TagHandler oluşturur
func NewTagHandler(tagService service.ITagService) *TagHandler {
	return &TagHandler{
		tagService: tagService,
	}
}

// RegisterRoutes rotaları kayıt eder
func (h *TagHandler) RegisterRoutes(router fiber.Router, authMw fiber.Handler, adminMw fiber.Handler) {
	// Herkese açık rotalar
	router.Get("/tags", h.ListTags)
	router.Get("/tags/:id", h.GetTag)
	router.Get("/tags/slug/:slug", h.GetTagBySlug)
	router.Get("/tags/article/:articleID", h.GetTagsByArticle)

	// Sadece admin rotaları
	adminRoutes := router.Group("/admin/tags", adminMw)
	adminRoutes.Post("/", h.CreateTag)
	adminRoutes.Put("/:id", h.UpdateTag)
	adminRoutes.Delete("/:id", h.DeleteTag)
}

// ListTags etiketleri listeler
// @Summary Etiketleri listele
// @Description Tüm etiketleri sayfalanmış şekilde listeler
// @Tags Etiketler
// @Accept json
// @Produce json
// @Param page query int false "Sayfa numarası (varsayılan: 1)"
// @Param limit query int false "Sayfa başına sonuç sayısı (varsayılan: 20, maksimum: 100)"
// @Success 200 {object} domain.PaginatedResponse{data=[]domain.Tag}
// @Router /tags [get]
func (h *TagHandler) ListTags(c *fiber.Ctx) error {
	// Sayfalama parametrelerini al
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	tags, total, err := h.tagService.ListTags(offset, limit)
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
		"data": tags,
		"meta": fiber.Map{
			"current_page": page,
			"per_page":     limit,
			"total":        total,
			"total_pages":  totalPages,
		},
	})
}

// GetTag ID'ye göre etiket getirir
// @Summary Etiket detayı getir
// @Description ID'ye göre etiket detayı getirir
// @Tags Etiketler
// @Accept json
// @Produce json
// @Param id path int true "Etiket ID"
// @Success 200 {object} domain.Tag
// @Failure 400 {object} domain.ErrorResponse "Geçersiz etiket ID"
// @Failure 404 {object} domain.ErrorResponse "Etiket bulunamadı"
// @Router /tags/{id} [get]
func (h *TagHandler) GetTag(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Geçersiz etiket ID")
	}

	tag, err := h.tagService.GetTagByID(uint(id))
	if err != nil {
		return err
	}

	return c.JSON(tag)
}

// GetTagBySlug slug'a göre etiket getirir
// @Summary Slug ile etiket getir
// @Description Slug'a göre etiket detayı getirir
// @Tags Etiketler
// @Accept json
// @Produce json
// @Param slug path string true "Etiket Slug"
// @Success 200 {object} domain.Tag
// @Failure 404 {object} domain.ErrorResponse "Etiket bulunamadı"
// @Router /tags/slug/{slug} [get]
func (h *TagHandler) GetTagBySlug(c *fiber.Ctx) error {
	slug := c.Params("slug")
	tag, err := h.tagService.GetTagBySlug(slug)
	if err != nil {
		return err
	}

	return c.JSON(tag)
}

// GetTagsByArticle makaleye göre etiketleri getirir
// @Summary Makaleye göre etiketleri getir
// @Description Belirli bir makaleye bağlı tüm etiketleri getirir (KULLANIM DIŞI)
// @Tags Etiketler
// @Accept json
// @Produce json
// @Param articleID path int true "Makale ID"
// @Success 200 {array} domain.Tag
// @Failure 400 {object} domain.ErrorResponse "Geçersiz makale ID"
// @Failure 501 {object} domain.ErrorResponse "Metot desteklenmiyor"
// @Router /tags/article/{articleID} [get]
// @Deprecated
func (h *TagHandler) GetTagsByArticle(c *fiber.Ctx) error {
	return fiber.NewError(fiber.StatusNotImplemented, "Bu metot artık desteklenmiyor. Makale modelini Tags ilişkisi ile kullanınız.")
}

// CreateTag yeni bir etiket oluşturur
// @Summary Yeni etiket oluştur
// @Description Yeni bir etiket oluşturur (Sadece admin)
// @Tags Admin, Etiketler
// @Accept json
// @Produce json
// @Param tag body domain.CreateTagRequest true "Etiket bilgileri"
// @Success 201 {object} domain.Tag
// @Failure 400 {object} domain.ErrorResponse "Geçersiz istek"
// @Failure 401 {object} domain.ErrorResponse "Yetkisiz erişim"
// @Failure 403 {object} domain.ErrorResponse "Yetersiz yetki"
// @Security ApiKeyAuth
// @Router /admin/tags [post]
func (h *TagHandler) CreateTag(c *fiber.Ctx) error {
	var req domain.CreateTagRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Geçersiz istek formatı")
	}

	tag, err := h.tagService.CreateTag(req.Name, req.Slug)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(tag)
}

// UpdateTag etiketi günceller
// @Summary Etiket güncelle
// @Description Mevcut bir etiketi günceller (Sadece admin)
// @Tags Admin, Etiketler
// @Accept json
// @Produce json
// @Param id path int true "Etiket ID"
// @Param tag body domain.UpdateTagRequest true "Etiket bilgileri"
// @Success 200 {object} domain.Tag
// @Failure 400 {object} domain.ErrorResponse "Geçersiz istek"
// @Failure 401 {object} domain.ErrorResponse "Yetkisiz erişim"
// @Failure 403 {object} domain.ErrorResponse "Yetersiz yetki"
// @Failure 404 {object} domain.ErrorResponse "Etiket bulunamadı"
// @Security ApiKeyAuth
// @Router /admin/tags/{id} [put]
func (h *TagHandler) UpdateTag(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Geçersiz etiket ID")
	}

	var req domain.UpdateTagRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Geçersiz istek formatı")
	}

	tag, err := h.tagService.UpdateTag(uint(id), req.Name, req.Slug)
	if err != nil {
		return err
	}

	return c.JSON(tag)
}

// DeleteTag etiketi siler
// @Summary Etiket sil
// @Description Bir etiketi siler (Sadece admin)
// @Tags Admin, Etiketler
// @Accept json
// @Produce json
// @Param id path int true "Etiket ID"
// @Success 204 "Başarıyla silindi"
// @Failure 400 {object} domain.ErrorResponse "Geçersiz istek"
// @Failure 401 {object} domain.ErrorResponse "Yetkisiz erişim"
// @Failure 403 {object} domain.ErrorResponse "Yetersiz yetki"
// @Failure 404 {object} domain.ErrorResponse "Etiket bulunamadı"
// @Security ApiKeyAuth
// @Router /admin/tags/{id} [delete]
func (h *TagHandler) DeleteTag(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Geçersiz etiket ID")
	}

	err = h.tagService.DeleteTag(uint(id))
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}
