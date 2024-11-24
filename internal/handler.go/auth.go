package handler

import (
	"dim_kurs/internal/domain/request"
	"dim_kurs/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type IAuth interface {
	Login(c *gin.Context)
	LoginPage(c *gin.Context)
	Register(c *gin.Context)
	RegisterPage(c *gin.Context)
}

type Auth struct {
	authUseCase usecase.IAuth
}

func NewAuth(authUseCase usecase.IAuth) *Auth {
	return &Auth{
		authUseCase: authUseCase,
	}
}

func (h *Auth) Login(c *gin.Context) {
	var req request.Login
	if err := c.ShouldBind(&req); err != nil {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{"error": "Invalid login request"})
		return
	}

	accessToken, err := h.authUseCase.Login(c.Request.Context(), req)
	if err != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{"error": "Invalid credentials"})
		return
	}

	c.SetCookie("Authorization", "Bearer "+accessToken, 3600, "/", "", false, true)

	redirectTo := c.PostForm("redirect_to")
	if redirectTo == "" {
		redirectTo = "/"
	}

	c.Redirect(http.StatusFound, redirectTo)
}

func (h *Auth) LoginPage(c *gin.Context) {
	redirectTo := c.DefaultQuery("redirect_to", "/")

	c.HTML(http.StatusOK, "login.html", gin.H{"redirect_to": redirectTo})
}

func (h *Auth) Register(c *gin.Context) {
	var req request.Register
	if err := c.ShouldBind(&req); err != nil {
		c.HTML(http.StatusBadRequest, "register.html", gin.H{"error": "Invalid register request"})
		return
	}

	accessToken, err := h.authUseCase.Register(c.Request.Context(), req)
	if err != nil {
		c.HTML(http.StatusUnauthorized, "register.html", gin.H{"error": "Invalid credentials"})
		return
	}

	c.SetCookie("Authorization", "Bearer "+accessToken, 3600, "/", "", false, true)

	redirectTo := c.PostForm("redirect_to")
	if redirectTo == "" {
		redirectTo = "/"
	}

	c.Redirect(http.StatusFound, redirectTo)
}

func (h *Auth) RegisterPage(c *gin.Context) {
	redirectTo := c.DefaultQuery("redirect_to", "/")

	c.HTML(http.StatusOK, "register.html", gin.H{"redirect_to": redirectTo})
}
