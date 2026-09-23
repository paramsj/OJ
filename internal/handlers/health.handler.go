package handlers

import (
	"github.com/labstack/echo/v5"
	"net/http"
)

func Health(c *echo.Context) error {
	return c.String(http.StatusOK, "GG WELL PLAYED\n")
}
