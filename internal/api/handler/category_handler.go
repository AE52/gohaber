package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/username/haber/internal/domain"
	"github.com/username/haber/internal/service"
)

// CategoryHandler kategori işleyicileri
type CategoryHandler struct {
	categoryService service.ICategoryService
}

// NewCategoryHandler yeni bir CategoryHandler oluşturur
func NewCategoryHandler(categoryService service.ICategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

// RegisterRoutes rotaları kayıt eder
func (h *CategoryHandler) RegisterRoutes(router fiber.Router, authMw fiber.Handler, adminMw fiber.Handler) {
	// Herkese açık rotalar
	router.Get("/categories", h.ListCategories)
	router.Get("/categories/:id", h.GetCategory)
	router.Get("/categories/slug/:slug", h.GetCategoryBySlug)

	// Sadece admin rotaları
	adminRoutes := router.Group("/admin/categories", adminMw)
	adminRoutes.Post("/", h.CreateCategory)
	adminRoutes.Put("/:id", h.UpdateCategory)
	adminRoutes.Delete("/:id", h.DeleteCategory)
}

// ListCategories kategorileri listeler
func (h *CategoryHandler) ListCategories(c *fiber.Ctx) error {
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

	// Ana kategorileri mi almak istiyoruz?
	parentID := c.Query("parent_id", "0")
	var filter map[string]interface{}
	if parentID == "0" {
		filter = map[string]interface{}{"parent_id": nil}
	} else {
		pid, err := strconv.ParseUint(parentID, 10, 32)
		if err == nil {
			filter = map[string]interface{}{"parent_id": uint(pid)}
		}
	}

	// Kategorileri getir
	categories, total, err := h.categoryService.ListCategories(offset, limit, filter)
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
		"data": categories,
		"meta": fiber.Map{
			"current_page": page,
			"per_page":     limit,
			"total":        total,
			"total_pages":  totalPages,
		},
	})
}

// GetCategory ID'ye göre kategori getirir
func (h *CategoryHandler) GetCategory(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Geçersiz kategori ID")
	}

	category, err := h.categoryService.GetCategoryByID(uint(id))
	if err != nil {
		return err
	}

	return c.JSON(category)
}

// GetCategoryBySlug slug'a göre kategori getirir
func (h *CategoryHandler) GetCategoryBySlug(c *fiber.Ctx) error {
	slug := c.Params("slug")
	category, err := h.categoryService.GetCategoryBySlug(slug)
	if err != nil {
		return err
	}

	return c.JSON(category)
}

// CreateCategory yeni bir kategori oluşturur
func (h *CategoryHandler) CreateCategory(c *fiber.Ctx) error {
	var req domain.CreateCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Geçersiz istek formatı")
	}

	category, err := h.categoryService.CreateCategory(&req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(category)
}

// UpdateCategory kategoriyi günceller
func (h *CategoryHandler) UpdateCategory(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Geçersiz kategori ID")
	}

	var req domain.UpdateCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Geçersiz istek formatı")
	}

	category, err := h.categoryService.UpdateCategory(uint(id), &req)
	if err != nil {
		return err
	}

	return c.JSON(category)
}

// DeleteCategory kategoriyi siler
func (h *CategoryHandler) DeleteCategory(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Geçersiz kategori ID")
	}

	err = h.categoryService.DeleteCategory(uint(id))
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}
