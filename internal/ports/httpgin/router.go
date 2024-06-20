package httpgin

import (
	app "comixsearch/internal/app"
	"context"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// AppRouter sets up routes for getting comices in application.
func AppRouter(ctx context.Context, r gin.IRouter, a app.App) {
	r.POST("/search", search(ctx, a))
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
