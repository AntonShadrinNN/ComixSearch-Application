package httpgin

import "github.com/gin-gonic/gin"

type userRequest struct {
	Keywords string `json:"keywords"`
}

func createSuccessResp(data map[string]string, err error) *gin.H {
	return &gin.H{
		"comices": data,
		"error":   err,
	}
}

func createErrorResp(err error) *gin.H {
	return &gin.H{
		"comices": "",
		"error":   err.Error(),
	}
}
