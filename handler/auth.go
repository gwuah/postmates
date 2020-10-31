package handler

import (
	"net/http"

	"github.com/electra-systems/core-api/database/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) Refresh(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "token not found",
		})
		return
	}

	customer := new(models.Customer)
	result := h.DB.Where("token = ?", token).First(customer)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to fetch customer",
		})
		return
	}

	newToken, err := h.JWT.GenerateToken(customer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to generate token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"token":   newToken,
	})
	return
}
