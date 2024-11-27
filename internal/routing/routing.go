package routing

import (
	"dim_kurs/internal/handler.go"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	handler.IAuth
	handler.IMiddleware
	handler.IUser
}

func InitRoutes(router *gin.Engine, handlers Handlers) {

	router.Use(handlers.IMiddleware.MethodOverride)

	api := router.Group("/api/v1")
	{
		api.GET("/", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "index.html", nil)
		})

		profile := api.Group("/profile", handlers.IMiddleware.AccountIdentity)
		{
			profile.GET("/", handlers.IUser.ProfilePage)
			profile.POST("/", handlers.IUser.Logout)
			profile.POST("/update", handlers.IUser.Update)
		}
		// admin := api.Group("/admin")
		// {
		// 	admin.GET("/", handlers.IUser.ProfilePage)
		// 	admin.POST("/", handlers.IUser.Logout)
		// 	admin.PUT("/update", handlers.IUser.Update)
		// }
		auth := api.Group("/auth")
		{
			auth.GET("/login", handlers.LoginPage)
			auth.POST("/login", handlers.Login)
			auth.GET("/register", handlers.RegisterPage)
			auth.POST("/register", handlers.Register)
		}
	}
}
