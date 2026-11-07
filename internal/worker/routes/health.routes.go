package routes

import (
	"agentV2/internal/worker/controllers"

	"github.com/labstack/echo/v5"
)

func RegisterWorkerHealth(e *echo.Echo) {
	e.GET("/health", controllers.GetHealth)
}
