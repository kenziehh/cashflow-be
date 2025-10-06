package middleware

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Logger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Process request
		err := c.Next()

		// Log details
		duration := time.Since(start)
		status := c.Response().StatusCode()
		method := c.Method()
		path := c.Path()
		ip := c.IP()

		// Color codes for status
		statusColor := getStatusColor(status)
		
		logMsg := fmt.Sprintf("[%s] %s | %s%3d%s | %13v | %15s | %-7s %s",
			time.Now().Format("2006/01/02 - 15:04:05"),
			"FIBER",
			statusColor,
			status,
			"\033[0m",
			duration,
			ip,
			method,
			path,
		)

		log.Println(logMsg)

		return err
	}
}

func getStatusColor(status int) string {
	switch {
	case status >= 200 && status < 300:
		return "\033[32m" // Green
	case status >= 300 && status < 400:
		return "\033[36m" // Cyan
	case status >= 400 && status < 500:
		return "\033[33m" // Yellow
	default:
		return "\033[31m" // Red
	}
}