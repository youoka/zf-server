package statics

import "github.com/gin-gonic/gin"

func Index(c *gin.Context) {
	c.HTML(200, "index.html", nil)
}
func Login(c *gin.Context) {
	c.HTML(200, "login.html", nil)
}
func PushOrder(c *gin.Context) {
	c.HTML(200, "push_order.html", nil)
}
func PullOrder(c *gin.Context) {
	c.HTML(200, "pull_order.html", nil)
}
