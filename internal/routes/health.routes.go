package routes

import (
	"agentV2/internal/controllers"

	"github.com/labstack/echo/v5"
)

func RegisterHealthRoutes(e *echo.Echo) {
	e.Get("/health", controllers.GetHealth)
}
