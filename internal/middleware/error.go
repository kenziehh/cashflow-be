package middleware

import (
	"errors"
	"fmt"
	"log"
	"runtime/debug"

	"devsecops-be/pkg/errx"
	"devsecops-be/pkg/response"

	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	// Default response
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	// Log dasar
	fmt.Printf("[ERROR] %v\n", err)

	// Cek apakah error adalah AppError (custom error)
	if appErr, ok := errx.IsAppError(err); ok {
		code = appErr.Code
		message = appErr.Message
		log.Printf("[AppError] %s | Status: %d | Path: %s", appErr.Message, appErr.Code, c.Path())
		return c.Status(code).JSON(response.ErrorResponse(message))
	}

	// Cek apakah error bawaan Fiber (404, 405, dll)
	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		code = fiberErr.Code
		message = fiberErr.Message
	}

	// Handle error spesifik berdasarkan status code
	switch {
	case code == fiber.StatusNotFound:
		return c.Status(code).JSON(response.ErrorResponse("Route not found"))
	case code == fiber.StatusMethodNotAllowed:
		return c.Status(code).JSON(response.ErrorResponse("Method not allowed"))
	case code == fiber.StatusRequestTimeout:
		return c.Status(code).JSON(response.ErrorResponse("Request timeout"))
	case code >= fiber.StatusBadRequest && code < fiber.StatusInternalServerError:
		return c.Status(code).JSON(response.ErrorResponse(message))
	default:
		// Log error tidak terduga dengan stacktrace untuk debugging
		// requestID := c.Locals("requestid")
		log.Printf("[UnexpectedError] %v | Path: %s\nStacktrace:\n%s", err, c.Path(), debug.Stack())

		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse(
			"Something went wrong",
		))
	}
}
