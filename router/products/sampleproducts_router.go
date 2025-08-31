package router

import (
	ctx "challenge-project/controller/products"

	"github.com/labstack/echo/v4"
)

func ProductRoute(api *echo.Group) {
	api.OPTIONS("/product/get-list/all/v1", ctx.ListProducts)
	api.GET("/product/get-list/all/v1", ctx.ListProducts)

	api.OPTIONS("/product/view/:id/v1", ctx.GetProduct)
	api.GET("/product/view/:id/v1", ctx.GetProduct)

	api.OPTIONS("/product/create/v1", ctx.CreateProduct)
	api.POST("/product/create/v1", ctx.CreateProduct)

	api.OPTIONS("/product/update/:id/v1", ctx.UpdateProduct)
	api.PUT("/product/update/:id/v1", ctx.UpdateProduct)

	api.OPTIONS("/product/delete/:id/v1", ctx.DeleteProduct)
	api.DELETE("/product/delete/:id/v1", ctx.DeleteProduct)
}
