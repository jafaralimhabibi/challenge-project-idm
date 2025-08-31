package router

import (
	ctx "challenge-project/controller/auth"

	"github.com/labstack/echo/v4"
)

func UserAuth(api *echo.Group) {
	api.OPTIONS("/login/v1", ctx.Login)
	api.POST("/login/v1", ctx.Login)

	api.OPTIONS("/register/v1", ctx.Register)
	api.POST("/register/v1", ctx.Register)
}
