package routes

import (
	"agentV2/internal/worker/controllers"

	"github.com/labstack/echo/v5"
)

func RegisterSubmissionRoutes(e *echo.Echo) {
	e.GET("/submit", controllers.Submit)
}
