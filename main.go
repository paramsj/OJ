package main

import (
	"net/http"

	"github.com/labstack/echo/v5"

	"github.com/labstack/echo/v5/middleware"

	"example.com/paramsj/my-go-project/internal/handlers"

	"github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i any) error {
	if err := cv.validator.Struct(i); err != nil {
		// Optionally return the error to let each route control the status code.
		return echo.ErrBadRequest.Wrap(err)
	}
	return nil
}

func main() {
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	e.GET("/", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"message":    "Hello, World!",
			"anotherOne": "This is truly something!",
		})
	})
	e.GET("/health", handlers.Health)
	e.POST("/submit", handlers.HandleSubmission)
	if err := e.Start(":1323"); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
