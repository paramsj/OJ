package routes

import (
	"agentV2/internal/agent/controllers"

	"github.com/labstack/echo/v5"
)

func RegisterSubmissionRoutes(e *echo.Echo) {
	e.POST("/submit", controllers.CreateSubmission)
}
