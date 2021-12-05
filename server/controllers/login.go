package controllers

import (
	"net/http"
	"strconv"
	"this/database"
	"this/database/models"
	"this/gate"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type LoginForm struct {
	Account  string `form:"account" binding:"required"`
	Password string `form:"password" binding:"required"`
}

// @Summary Login
// @Description Login
// @Tags Before Authorization
// @accept x-www-form-urlencoded
// @Produce json
// @param account formData string true "Account"
// @param password formData string true "Password"
// @Success 200 "{ ok , data , token }"
// @Failure 404 "{ error }"
// @Router /auth/login [post]
func Login(c *gin.Context) {
	var form LoginForm
	if c.ShouldBind(&form) != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Check for empty fields!",
		})
		return
	}

	var user *models.User
	user, err := user.GetUserByAccountPassword(*database.MAIN, form.Account, form.Password) // Main DB 验证
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Account/password error!",
		})
		return
	}

	var usertype *models.Usertype
	usertype, errType := usertype.GetById(*database.MAIN, user.TypeID)
	if errType != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Usertype error!",
		})
		return
	}

	if !(usertype.Name == "Admin" || usertype.Name == "Office" || usertype.Name == "Account") {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Usertype " + usertype.Name + " is unavailable!",
		})
		return
	}

	now := time.Now()
	claims := gate.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:       strconv.FormatInt(int64(user.ID), 10),
			Issuer:   "20211108-daka",
			Subject:  "user",
			Audience: []string{"20211108-daka", "user"},
			ExpiresAt: &jwt.NumericDate{
				Time: now.AddDate(0, 0, 15), // 15 day
			},
			NotBefore: &jwt.NumericDate{
				Time: now,
			},
			IssuedAt: &jwt.NumericDate{
				Time: now,
			},
		},
		UserID:   user.ID,
		Usertype: usertype.Name,
	}

	token, err := gate.CreateToken(&claims)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Token generate error!",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ok": true,
		"data": gin.H{
			"user": gin.H{
				"name": user.Account,
				"type": usertype.Name,
			},
		},
		"token": token,
	})
}
