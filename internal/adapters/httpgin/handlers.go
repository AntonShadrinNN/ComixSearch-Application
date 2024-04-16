package httpgin

import (
	app "comixsearch/internal/app"
	"context"

	"github.com/gin-gonic/gin"
)

func search(a app.SearchApp) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := []string{",", ""}
		ctx := context.Background()
		b := true
		a.Search(ctx, key, b)
	}
}
