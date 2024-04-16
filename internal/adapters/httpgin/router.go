package httpgin

import (
	app "comixsearch/internal/app"

	"github.com/gin-gonic/gin"
)

func AppRouter(r gin.IRouter, a app.SearchApp) {
	r.GET("/Search", search(a))
}
