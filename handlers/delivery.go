package handler

import (
	"fmt"
	"net/http"

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
	ProductId   int    `json:"productId"`
	Notes       string `json:"notes"`
	CustomerID  int    `json:"customerId"`
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

	delivery := models.Delivery{
		OriginLatitude:       data.Origin.Latitude,
		OriginLongitude:      data.Origin.Longitude,
		DestinationLatitude:  data.Destination.Latitude,
		DestinationLongitude: data.Destination.Longitude,
		Notes:                data.Notes,
		CustomerID:           data.CustomerID,
	}

	fmt.Println(delivery)

	c.JSON(http.StatusOK, gin.H{
		"message":  "Success",
		"delivery": product,
	})

}
