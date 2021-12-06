package controllers

import (
	"net/http"
	"this/database"
	"this/database/models"
	"this/gate"

	"github.com/gin-gonic/gin"
)

type AddUserForm struct {
	StoreID   int    `form:"StoreID"`
	Firstname string `form:"firstname" binding:"required"`
	Lastname  string `form:"lastname" binding:"required"`
	Usertype  uint   `form:"usertype" binding:"required"`
	Password  string `form:"password" binding:"required"`
}

// @Summary AddUser
// @Description Add User
// @Tags After Authorization
// @accept x-www-form-urlencoded
// @Produce json
// @Security BearerIdAuth
// @param Authorization header string false "Authorization"
// @param StoreID formData int false "Store ID"
// @param firstname formData string true "First Name"
// @param lastname formData string true "Last Name"
// @param usertype formData uint true "User Type"
// @param password formData string true "Password"
// @Success 200 "{ ok }"
// @Failure 401 "Auth failed"
// @Failure 404 "{ error }"
// @Router /user/add/ [post]
func AddUser(c *gin.Context) {
	temp, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Middleware Auth not exists!",
		})
		return
	}
	claims := temp.(*gate.Claims)

	if !(claims.Usertype == "Admin" || claims.Usertype == "Office") {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Usertype " + claims.Usertype + " is unavailable!",
		})
		return
	}

	var form AddUserForm
	if c.ShouldBind(&form) != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Check for empty fields!",
		})
		return
	}

	if form.StoreID != database.WAREHOUSE_INDEX {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Store is not WAREHOUSE!",
		})
		return
	}

	store, ok := database.GetStore(form.StoreID)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Store error!",
		})
		return
	}

	var user *models.User
	err := user.AddUser(*store, form.Firstname, form.Lastname, form.Password, form.Usertype)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "AddUser error!",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ok": true,
	})
}
