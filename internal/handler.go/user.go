package handler

import (
	"dim_kurs/internal/domain/request"
	"dim_kurs/internal/usecase"
	"dim_kurs/pkg/token"
	"net/http"

	"github.com/gin-gonic/gin"
)

type IUser interface {
	ProfilePage(c *gin.Context)
	Update(c *gin.Context)
	Logout(c *gin.Context)
}

type User struct {
	userUseCase usecase.IUser
}

func NewUser(userUseCase usecase.IUser) *User {
	return &User{
		userUseCase: userUseCase,
	}
}

func (h *User) ProfilePage(c *gin.Context) {
	tknClaims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "empty claims"})
		return
	}
	claims, ok := tknClaims.(token.AuthInfo)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "token isn't string"})
		return
	}

	user, err := h.userUseCase.GetProfile(c, claims)
	if err != nil {
		c.HTML(http.StatusUnauthorized, "userprofile.html", gin.H{"error": "Invalid credentials"})
		return
	}

	c.HTML(http.StatusOK, "userprofile.html", user)
}

func (h *User) Update(c *gin.Context) {
	tknClaims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "empty claims"})
		return
	}
	claims, ok := tknClaims.(token.AuthInfo)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "token isn't string"})
		return
	}

	var req request.UpdateUser
	redirectTo := c.DefaultQuery("redirect_to", "/api/v1")

	if err := c.ShouldBind(&req); err != nil {
		c.HTML(http.StatusUnauthorized, "userprofile.html", gin.H{"error": "Invalid credentials"})
		return
	}

	if err := h.userUseCase.Update(c, claims, req); err != nil {
		c.HTML(http.StatusUnauthorized, "userprofile.html", gin.H{"error": "Invalid credentials"})
		return
	}

	c.Redirect(http.StatusFound, redirectTo)
}

func (h *User) Logout(c *gin.Context) {
	tknClaims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "empty claims"})
		return
	}
	claims, ok := tknClaims.(token.AuthInfo)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "token isn't string"})
		return
	}

	redirectTo := c.DefaultQuery("redirect_to", "/api/v1")

	if err := h.userUseCase.Logout(c, claims); err != nil {
		c.HTML(http.StatusUnauthorized, "userprofile.html", gin.H{"error": "Invalid credentials"})
		return
	}

	c.Redirect(http.StatusFound, redirectTo)
}
