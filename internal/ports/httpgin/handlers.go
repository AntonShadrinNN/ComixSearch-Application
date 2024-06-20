package httpgin

import (
	app "comixsearch/internal/app"
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// search retrieves comices based on user input.
func search(ctx context.Context, a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			limit int
			err   error
		)
		if query, ok := c.GetQuery("limit"); ok {
			if limit, err = strconv.Atoi(query); err != nil {
				c.JSON(http.StatusBadRequest, createErrorResp(err))
				return
			}
		} else {
			limit = 10
		}
		var reqBody userRequest
		if err := c.Bind(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, createErrorResp(err))
			return
		}

		keywords := strings.Fields(reqBody.Keywords)
		data, err := a.Search(ctx, keywords, limit)

		if err != nil {
			c.Status(http.StatusInternalServerError)
			c.JSON(http.StatusInternalServerError, createErrorResp(err))
			return
		}

		c.JSON(http.StatusOK, createSuccessResp(data, nil))
	}
}
