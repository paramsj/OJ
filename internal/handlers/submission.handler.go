package handlers

import (
	"net/http"

	"example.com/paramsj/my-go-project/internal/executor"
	"github.com/labstack/echo/v5"
)

type RunReq struct {
	ID       string `json:"id" validate:"required"`
	Code     string `json:"code" validate:"required"`
	Language string `json:"language" validate:"required"`
}

func HandleSubmission(c *echo.Context) error {

	req := new(RunReq)

	if err := c.Bind(req); err != nil {
		return err
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	if req.Language != "cpp" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "only cpp is supported (got " + req.Language + ")",
		})
	}

	result, err := executor.ExecuteCpp(req.Code)

	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, result)
}
