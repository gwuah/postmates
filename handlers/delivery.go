package handler

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gwuah/api/database/models"
	"github.com/gwuah/api/repository"
)

type Coord struct {
	Longitude float64 `json:"longitude" validate:"required"`
	Latitude  float64 `json:"latitude" validate:"required"`
}

type DeliveryRequest struct {
	Origin      Coord  `json:"origin"`
	Destination Coord  `json:"destination"`
	ProductId   uint   `json:"productId"`
	Notes       string `json:"notes"`
	CustomerID  uint   `json:"customerId"`
}

func (h *Handler) ListDeliveries(c *gin.Context) {

}

func (h *Handler) ViewDeliveries(c *gin.Context) {

}

func (h *Handler) CreateDelivery(c *gin.Context) {
	var data DeliveryRequest
	repo := repository.New(h.DB)

	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Bad Request",
		})
		return
	}

	product, result := repo.FindProduct(data.ProductId)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Product Doesn't Exist",
		})
		return
	}

	if strings.ToLower(product.Name) == "express" {
		order := models.Order{}

		if err := h.DB.Create(&order).Error; err != nil {
			log.Println("Failed To Create Order", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Request Failed",
			})
			return
		}

		delivery := models.Delivery{
			OriginLatitude:       data.Origin.Latitude,
			OriginLongitude:      data.Origin.Longitude,
			DestinationLatitude:  data.Destination.Latitude,
			DestinationLongitude: data.Destination.Longitude,
			Notes:                data.Notes,
			OrderID:              order.ID,
			ProductID:            data.ProductId,
			CustomerID:           data.CustomerID,
		}

		if err := h.DB.Create(&delivery).Error; err != nil {
			log.Printf("Failed to create delivery for order [ %d ] ", order.ID)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Request Failed",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":  "Success",
			"delivery": delivery,
		})

	} else if strings.ToLower(product.Name) == "pool" {
		c.JSON(http.StatusOK, gin.H{
			"message": "Work on going for this product",
		})

	}

}
