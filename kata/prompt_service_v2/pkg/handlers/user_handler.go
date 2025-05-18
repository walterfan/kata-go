package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/walterfan/prompt-service/pkg/database"
	"github.com/walterfan/prompt-service/pkg/models"

	"github.com/gin-gonic/gin"
)

func CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func GetUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.DB.Model(&user).Updates(input)
	c.JSON(http.StatusOK, user)
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if err := database.DB.Delete(&models.User{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}

func SearchUsers(c *gin.Context) {
	keyword := c.Query("q")
	pageNum := c.DefaultQuery("pageNum", "1")
	pageSize := c.DefaultQuery("pageSize", "20")

	var users []models.User
	query := database.DB.Model(&models.User{})

	if keyword != "" {
		kw := "%" + strings.ToLower(keyword) + "%"
		query = query.Where("LOWER(username) LIKE ? OR LOWER(email) LIKE ?", kw, kw)
	}

	// Pagination
	var pageSizeValue int = 20
	if s, err := strconv.Atoi(pageSize); err == nil && s > 0 {
		pageSizeValue = s
	}

	var offset int
	if p, err := strconv.Atoi(pageNum); err == nil && p > 0 {
		offset = (p - 1) * pageSizeValue
	}

	query = query.Order("updated_at desc").Limit(pageSizeValue).Offset(offset)

	if err := query.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}