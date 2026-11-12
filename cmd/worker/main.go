package main

import (
	"agentV2/internal/worker/controllers"
	"agentV2/internal/worker/routes"
	"agentV2/internal/worker/vm"
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func main() {
	newmgr := vm.NewManager()
	ctrl := controllers.NewController(newmgr)
	e := echo.New()

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	routes.RegisterWorkerHealth(e)
	routes.RegisterSubmissionRoutes(e, ctrl)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	sc := echo.StartConfig{
		Address:         ":1323",
		GracefulTimeout: 5 * time.Second,
	}
	if err := sc.Start(ctx, e); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}

}
