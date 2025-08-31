package lib

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type StandardResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func JSONSuccess(c echo.Context, data interface{}) error {
	resp := StandardResponse{
		Code:    http.StatusOK,
		Message: "success",
		Data:    data,
	}
	return c.JSON(http.StatusOK, resp)
}

func JSONCreated(c echo.Context, data interface{}) error {
	resp := StandardResponse{
		Code:    http.StatusCreated,
		Message: "created",
		Data:    data,
	}
	return c.JSON(http.StatusCreated, resp)
}

func JSONError(c echo.Context, code int, msg string) error {
	resp := StandardResponse{
		Code:    code,
		Message: msg,
	}
	return c.JSON(code, resp)
}
