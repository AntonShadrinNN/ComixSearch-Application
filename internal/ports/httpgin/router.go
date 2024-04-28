package httpgin

import (
	app "comixsearch/internal/app"
	"context"

	"github.com/gin-gonic/gin"
)

// AppRouter sets up routes for getting comices in application.
func AppRouter(ctx context.Context, r gin.IRouter, a app.App) {
	r.POST("/search", search(ctx, a))
}
