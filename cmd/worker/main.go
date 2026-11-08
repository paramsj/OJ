package main

import (
	"agentV2/internal/worker/routes"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	routes.RegisterWorkerHealth(e)
	routes.RegisterSubmissionRoutes(e)
	if err := e.Start(":1323"); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
