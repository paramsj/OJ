package routes

import (
	"agentV2/internal/agent/controllers"

	"github.com/labstack/echo/v5"
)

func RegisterHealthRoutes(e *echo.Echo) {
	e.GET("/health", controllers.GetHealth)
}
