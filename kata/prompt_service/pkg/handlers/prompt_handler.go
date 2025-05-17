package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/walterfan/prompt-service/pkg/database"
	"github.com/walterfan/prompt-service/pkg/models"

	"github.com/gin-gonic/gin"
)

func CreatePrompt(c *gin.Context) {
	var prompt models.Prompt
	if err := c.ShouldBindJSON(&prompt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Create(&prompt).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, prompt)
}

func GetPrompt(c *gin.Context) {
	id := c.Param("id")
	var prompt models.Prompt
	if err := database.DB.First(&prompt, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Prompt not found"})
		return
	}
	c.JSON(http.StatusOK, prompt)
}

func UpdatePrompt(c *gin.Context) {
	id := c.Param("id")
	var prompt models.Prompt
	if err := database.DB.First(&prompt, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Prompt not found"})
		return
	}

	var input models.Prompt
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.DB.Model(&prompt).Updates(input)
	c.JSON(http.StatusOK, prompt)
}

func DeletePrompt(c *gin.Context) {
	id := c.Param("id")
	if err := database.DB.Delete(&models.Prompt{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}

func SearchPrompts(c *gin.Context) {
	keyword := c.Query("q")
	pageNum := c.DefaultQuery("pageNum", "1")
	pageSize := c.DefaultQuery("pageSize", "20")

	var prompts []models.Prompt
	query := database.DB.Model(&models.Prompt{})

	if keyword != "" {
		kw := "%" + strings.ToLower(keyword) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ? OR LOWER(tags) LIKE ?", kw, kw, kw)
	}


	// Convert pageNum and pageSize to integers
	var pageSizeValue int = 20
	if s, err := strconv.Atoi(pageSize); err == nil && s > 0 {
		pageSizeValue = s
	}

	var offset int
	if p, err := strconv.Atoi(pageNum); err == nil && p > 0 {
		offset = (p - 1) * pageSizeValue
	}


	query = query.Order("updated_at desc").Limit(pageSizeValue).Offset(offset)

	if err := query.Find(&prompts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, prompts)
}
