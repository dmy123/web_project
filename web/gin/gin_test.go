package gin

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"testing"
)

func TestUserController(t *testing.T) {
	g := gin.Default()
	c := &UserController{}
	g.GET("/user", c.GetUser)
	g.POST("/user", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello %s", "world")
	})
	g.GET("/static", func(context *gin.Context) {
		// 读文件
		// 谐响应
	})
	_ = g.Run(":8082")

	http.ListenAndServe(":8083", g)
}
