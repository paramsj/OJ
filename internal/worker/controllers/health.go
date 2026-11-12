package controllers

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

func GetHealth(c *echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"message": "Worker is running fine!"})
}
