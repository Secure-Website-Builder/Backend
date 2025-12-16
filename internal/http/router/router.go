package router

import (
	"github.com/Secure-Website-Builder/Backend/internal/http/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRouter(
	categoryHandler *handlers.CategoryHandler,
	productHandler *handlers.ProductHandler,
	categoryProductHandler *handlers.CategoryProductHandler,
	cartHandler *handlers.CartHandler,
) *gin.Engine {

	r := gin.Default()

	stores := r.Group("/stores/:store_id")
	{
		stores.GET("/categories", categoryHandler.ListCategories)
		stores.GET("/products", productHandler.ListProducts)
		stores.GET("/cart", cartHandler.GetCart)
	}

	productGroup := r.Group("/stores/:store_id/products")
	{
		productGroup.GET("/:product_id", productHandler.GetProduct)
	}

	categoryGroup := r.Group("/stores/:store_id/categories")
	{
		categoryGroup.GET("/:id/top-products", categoryProductHandler.GetTopProducts)
	}

	return r
}
