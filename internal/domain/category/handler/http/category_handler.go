package http

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/kenziehh/cashflow-be/internal/domain/category/service"
	"github.com/kenziehh/cashflow-be/pkg/response"
)

type CategoryHandler struct {
	service  service.CategoryService
	validate *validator.Validate
}

func NewCategoryHandler(service service.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		service:  service,
		validate: validator.New(),
	}
}

// GetAllCategories godoc
// @Summary      Get all categories
// @Description  Retrieve a list of all categories
// @Security     BearerAuth
// @Tags         Categories
// @Accept       json
// @Produce      json
// @Success      200  {object}  dto.GetAllCategoryResponse
// @Failure      500  {object}  map[string]string
// @Router       /categories [get]
func (h *CategoryHandler) GetAllCategories(c *fiber.Ctx) error {
	ctx := c.Context()

	categories, err := h.service.GetAllCategories(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch categories",
		})
	}

	return c.JSON(response.SuccessResponse("Categories retrieved successfully", categories))
}
