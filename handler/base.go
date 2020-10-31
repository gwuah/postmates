package handler

import (
	"log"
	"net/http"

	"github.com/electra-systems/core-api/shared"
	myValidator "github.com/electra-systems/core-api/utils/validator"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func (h *Handler) handleCustomerRating(c *gin.Context) {
	var data shared.CustomerRatingRequest
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failure",
			"err":     err,
		})
		return
	}

	response, err := h.Services.RateDelivery(shared.RatingRequest{
		IsCustomerRating: true,
		CustomerRating:   data,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failure",
			"err":     err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    response,
	})

}

func (h *Handler) handleCourierRating(c *gin.Context) {
	var data shared.CourierRatingRequest
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failure",
			"err":     err,
		})
		return
	}

	response, err := h.Services.RateDelivery(shared.RatingRequest{
		IsCustomerRating: false,
		CourierRating:    data,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failure",
			"err":     err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    response,
	})

}

func (h *Handler) GetDeliveryCost(c *gin.Context) {
	var quoteRequest shared.GetDeliveryCostRequest
	if err := c.ShouldBindJSON(&quoteRequest); err != nil {
		for _, fieldErr := range err.(validator.ValidationErrors) {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": myValidator.FieldError{Err: fieldErr}.String(),
			})
			return
		}
	}

	response, err := h.Services.GetDeliveryCost(quoteRequest)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failure",
			"err":     err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    response,
	})

}
