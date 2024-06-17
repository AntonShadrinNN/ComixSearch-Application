package httpgin

import "github.com/gin-gonic/gin"

type userRequest struct {
	Keywords string `json:"keywords" example:"earth"`
}

type Response struct {
	Comices map[string]string `example:"earth:http://xkcd/earth"`
	Error   error
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
