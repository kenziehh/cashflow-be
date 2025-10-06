package http

import (
	"strings"

	"devsecops-be/internal/domain/auth/dto"
	"devsecops-be/internal/domain/auth/service"
	"devsecops-be/pkg/errx"
	"devsecops-be/pkg/response"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AuthHandler struct {
	service  service.AuthService
	validate *validator.Validate
}

func NewAuthHandler(service service.AuthService) *AuthHandler {
	return &AuthHandler{
		service:  service,
		validate: validator.New(),
	}
}

// Register godoc
// @Summary Register new user
// @Description Register a new user with email, password, and name
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Register request"
// @Success 201 {object} response.Response{data=dto.AuthResponse}
// @Failure 400 {object} response.Response
// @Failure 409 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.NewBadRequestError("Invalid request body")
	}

	if err := h.validate.Struct(req); err != nil {
		return errx.NewBadRequestError(err.Error())
	}

	result, err := h.service.Register(c.Context(), &req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(response.SuccessResponse("User registered successfully", result))
}

// Login godoc
// @Summary User login
// @Description Login with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login request"
// @Success 200 {object} response.Response{data=dto.AuthResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.NewBadRequestError("Invalid request body")
	}

	if err := h.validate.Struct(req); err != nil {
		return errx.NewBadRequestError(err.Error())
	}

	result, err := h.service.Login(c.Context(), &req)
	if err != nil {
		return err
	}

	return c.JSON(response.SuccessResponse("Login successful", result))
}

// Logout godoc
// @Summary User logout
// @Description Logout and invalidate token
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")

	if err := h.service.Logout(c.Context(), token); err != nil {
		return err
	}

	return c.JSON(response.SuccessResponse("Logout successful", nil))
}

// GetProfile godoc
// @Summary Get user profile
// @Description Get current user profile
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=dto.UserProfile}
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /auth/me [get]
func (h *AuthHandler) GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	id, err := uuid.Parse(userID)
	if err != nil {
		return errx.NewBadRequestError("Invalid user ID")
	}

	profile, err := h.service.GetProfile(c.Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(response.SuccessResponse("Profile retrieved successfully", profile))
}
