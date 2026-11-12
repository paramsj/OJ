package routes

import (
	"agentV2/internal/worker/controllers"

	"github.com/labstack/echo/v5"
)

func RegisterSubmissionRoutes(e *echo.Echo, ctrl *controllers.Controller) {
	e.GET("/start", ctrl.Start)
	// e.POST("/submit", ctrl.Submit)
	e.DELETE("/delVM/:id", ctrl.Delete)
}
