package response

import "github.com/labstack/echo/v4"

type APIResponse struct {
	Status  string      `json:"status"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func Success(c echo.Context, code int, message string, data interface{}) error {
	return c.JSON(code, APIResponse{
		Status:  "success",
		Code:    code,
		Message: message,
		Data:    data,
	})
}

func Error(c echo.Context, code int, message string) error {
	return c.JSON(code, APIResponse{
		Status:  "error",
		Code:    code,
		Message: message,
	})
}
