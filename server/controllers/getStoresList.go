package controllers

import (
	"net/http"
	"this/database"

	"github.com/gin-gonic/gin"
)

// @Summary GetStoresList
// @Description Get Stores List
// @Tags Before Authorization
// @accept plain
// @Produce json
// @Success 200 "{ ok , data }"
// @Failure 404 "{ error }"
// @Router /stores/list [get]
func GetStoresList(c *gin.Context) {
	var list []string
	for _, store := range database.STORES {
		list = append(list, store.Name)
	}

	c.JSON(http.StatusOK, gin.H{
		"ok":   true,
		"data": list,
	})
}
