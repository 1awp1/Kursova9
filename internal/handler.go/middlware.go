package handler

import (
	"dim_kurs/internal/usecase"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type IMiddleware interface {
	AccountIdentity(c *gin.Context)
	MethodOverride(c *gin.Context)
}

const (
	authorizationHeader = "Authorization"
)

type Middleware struct {
	authUseCase usecase.IAuth
}

func NewMiddleware(authUseCase usecase.IAuth) *Middleware {
	return &Middleware{
		authUseCase: authUseCase,
	}
}

func (m *Middleware) AccountIdentity(c *gin.Context) {
	token, err := c.Cookie("Authorization")
	if err != nil || token == "" {
		redirectToLogin(c)
		return
	}

	parts := strings.SplitN(token, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" || len(parts[1]) == 0 {
		redirectToLogin(c)
		return
	}

	claims, err := m.authUseCase.VerifyToken(c.Request.Context(), parts[1])
	if err != nil {
		redirectToLogin(c)
		return
	}

	c.Set("claims", *claims)
	c.Next()
}

func redirectToLogin(c *gin.Context) {
	originalURL := c.Request.RequestURI
	c.Redirect(http.StatusFound, "api/v1/auth/login?redirect_to="+originalURL)
	c.Abort()
}

func (h *Middleware) MethodOverride(c *gin.Context) {
	method := c.DefaultPostForm("_method", "")

	fmt.Println(method)

	if method != "" {
		c.Request.Method = method
	}
	c.Next()
}
