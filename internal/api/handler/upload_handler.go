package handler

import (
	"path/filepath"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/username/haber/internal/domain"
	"github.com/username/haber/internal/service"
)

// UploadHandler dosya yükleme işleyicileri
type UploadHandler struct {
	uploadService service.IUploadService
}

// NewUploadHandler yeni bir UploadHandler oluşturur
func NewUploadHandler(uploadService service.IUploadService) *UploadHandler {
	return &UploadHandler{
		uploadService: uploadService,
	}
}

// RegisterRoutes rotaları kayıt eder
func (h *UploadHandler) RegisterRoutes(router fiber.Router, authMw fiber.Handler, adminMw fiber.Handler) {
	// Yalnızca giriş yapmış kullanıcılar yükleme yapabilir
	uploadRoutes := router.Group("/uploads", authMw)
	uploadRoutes.Post("/", h.UploadFile)
	uploadRoutes.Delete("/:id", h.DeleteFile)

	// Herkese açık
	router.Get("/uploads/:id", h.GetFile)
}

// UploadFile dosya yükleme işlemini gerçekleştirir
func (h *UploadHandler) UploadFile(c *fiber.Ctx) error {
	// Dosya bilgilerini al
	file, err := c.FormFile("file")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Dosya yüklenemedi")
	}

	// Kullanıcı ID'sini al
	userID := c.Locals("user_id").(uint)

	// Yükleme klasörünü belirle (form parametresinden)
	folder := c.FormValue("folder", "general")

	// Dosyayı yükle
	media, err := h.uploadService.UploadFile(file, folder, userID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(media)
}

// GetFile dosyayı getirir
func (h *UploadHandler) GetFile(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Geçersiz medya ID")
	}

	// Medya bilgilerini getir
	media, err := h.uploadService.GetMedia(uint(id))
	if err != nil {
		return err
	}

	// Dosya yolunu belirle
	filePath := filepath.Join("./uploads", media.ObjectName)

	// Dosyayı gönder
	return c.SendFile(filePath)
}

// DeleteFile dosyayı siler
func (h *UploadHandler) DeleteFile(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Geçersiz medya ID")
	}

	// Kullanıcı ID'sini al
	userID := c.Locals("user_id").(uint)

	// Medya bilgilerini getir (silmeden önce yetki kontrolü için)
	media, err := h.uploadService.GetMedia(uint(id))
	if err != nil {
		return err
	}

	// Yetki kontrolü yap (admin veya dosyanın sahibi olmalı)
	if media.UserID != userID && c.Locals("user_role") != domain.RoleAdmin {
		return fiber.NewError(fiber.StatusForbidden, "Bu dosyayı silme yetkiniz yok")
	}

	// Dosyayı sil
	err = h.uploadService.DeleteMedia(uint(id))
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}
