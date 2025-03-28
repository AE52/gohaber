package handler

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/username/haber/internal/domain"
	"github.com/username/haber/internal/service"
	"github.com/username/haber/pkg/auth"
)

// UserHandler kullanıcı işleyicileri
type UserHandler struct {
	userService service.IUserService
}

// NewUserHandler yeni bir UserHandler oluşturur
func NewUserHandler(userService service.IUserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// RegisterRoutes rotaları kayıt eder
func (h *UserHandler) RegisterRoutes(router fiber.Router, authMw fiber.Handler, adminMw fiber.Handler) {
	// Kimlik doğrulama gerektiren rotalar
	protectedRoutes := router.Group("/users", authMw)
	protectedRoutes.Get("/me", h.GetCurrentUser)
	protectedRoutes.Put("/me", h.UpdateCurrentUser)
	protectedRoutes.Put("/me/password", h.UpdatePassword)

	// Sadece admin rotaları
	adminRoutes := router.Group("/admin/users", adminMw)
	adminRoutes.Get("/", h.ListUsers)
	adminRoutes.Get("/:id", h.GetUser)
	adminRoutes.Post("/", h.CreateUser)
	adminRoutes.Put("/:id", h.UpdateUser)
	adminRoutes.Delete("/:id", h.DeleteUser)
}

// GetCurrentUser mevcut kullanıcı bilgilerini getirir
func (h *UserHandler) GetCurrentUser(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		return err
	}

	return c.JSON(user)
}

// UpdateCurrentUser mevcut kullanıcı bilgilerini günceller
func (h *UserHandler) UpdateCurrentUser(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	var req domain.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Geçersiz istek formatı")
	}

	// Mevcut kullanıcıyı getir
	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		return err
	}

	// Kullanıcı bilgilerini güncelle
	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.FullName != "" {
		user.FullName = req.FullName
	}
	if req.ProfileImage != "" {
		user.ProfileImage = req.ProfileImage
	}
	// Sadece admin rolü değiştirebilir
	if req.Role != "" && c.Locals("user_role") == domain.RoleAdmin {
		user.Role = req.Role
	}
	user.UpdatedAt = time.Now()

	err = h.userService.UpdateUser(user)
	if err != nil {
		return err
	}

	return c.JSON(user)
}

// UpdatePassword kullanıcı şifresini günceller
func (h *UserHandler) UpdatePassword(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	var req domain.UpdatePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Geçersiz istek formatı")
	}

	if req.NewPassword != req.ConfirmPassword {
		return fiber.NewError(fiber.StatusBadRequest, "Şifreler eşleşmiyor")
	}

	err := h.userService.UpdatePassword(userID, req.CurrentPassword, req.NewPassword)
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}

// ListUsers kullanıcıları listeler
func (h *UserHandler) ListUsers(c *fiber.Ctx) error {
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

	// Filtreler
	filters := make(map[string]interface{})
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}
	if role := c.Query("role"); role != "" {
		filters["role"] = role
	}

	users, total, err := h.userService.ListUsers(offset, limit, filters)
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
		"data": users,
		"meta": fiber.Map{
			"current_page": page,
			"per_page":     limit,
			"total":        total,
			"total_pages":  totalPages,
		},
	})
}

// GetUser ID'ye göre kullanıcı getirir
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Geçersiz kullanıcı ID")
	}

	user, err := h.userService.GetUserByID(uint(id))
	if err != nil {
		return err
	}

	return c.JSON(user)
}

// CreateUser yeni bir kullanıcı oluşturur
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var req domain.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Geçersiz istek formatı")
	}

	// Şifre kontrolü
	if req.Password != req.ConfirmPassword {
		return fiber.NewError(fiber.StatusBadRequest, "Şifreler eşleşmiyor")
	}

	// Yeni kullanıcı objesi oluştur
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		return err
	}

	now := time.Now()
	user := &domain.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		FullName:     req.FullName,
		Role:         req.Role,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	// Kullanıcıyı oluştur
	err = h.userService.CreateUser(user)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}

// UpdateUser kullanıcıyı günceller
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Geçersiz kullanıcı ID")
	}

	var req domain.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Geçersiz istek formatı")
	}

	// Mevcut kullanıcıyı getir
	user, err := h.userService.GetUserByID(uint(id))
	if err != nil {
		return err
	}

	// Kullanıcı bilgilerini güncelle
	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.FullName != "" {
		user.FullName = req.FullName
	}
	if req.ProfileImage != "" {
		user.ProfileImage = req.ProfileImage
	}
	if req.Role != "" {
		user.Role = req.Role
	}
	user.UpdatedAt = time.Now()

	err = h.userService.UpdateUser(user)
	if err != nil {
		return err
	}

	return c.JSON(user)
}

// DeleteUser kullanıcıyı siler
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Geçersiz kullanıcı ID")
	}

	err = h.userService.DeleteUser(uint(id))
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}
