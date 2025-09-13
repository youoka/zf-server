package common

import (
	"github.com/gin-gonic/gin"
)

func ApiErr(c *gin.Context, operationID string, code int, msg string) {
	c.JSON(200, gin.H{
		"operationID": operationID,
		"code":        code,
		"msg":         msg,
	})
}
func ApiSuccess(c *gin.Context, operationID string, msg string, data any) {
	c.JSON(200, gin.H{
		"operationID": operationID,
		"code":        200,
		"msg":         msg,
		"data":        data,
	})
}
