package routing

import (
	"dim_kurs/internal/handler.go"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	handler.IAuth
}

func InitRoutes(router *gin.Engine, handlers Handlers) {
	api := router.Group("/api/v1")
	{
		base := api.Group("/")
		{
			base.GET("", func(c *gin.Context) {
				c.HTML(http.StatusOK, "index.html", nil)
			}) //TODO
		}
		auth := api.Group("/auth")
		{
			auth.GET("/login", handlers.LoginPage)
			auth.POST("/login", handlers.Login)
			auth.GET("/register", handlers.RegisterPage)
			auth.POST("/register", handlers.Register)
		}
	}
}
