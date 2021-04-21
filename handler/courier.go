package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/gwuah/postmates/shared"
	myValidator "github.com/gwuah/postmates/utils/validator"
)

type closestCourierResponse struct {
	Couriers []string `json:"couriers"`
}

func (h *Handler) GetClosestCouriers(c *gin.Context) {
	var data shared.GetClosestCouriersRequest
	if err := c.ShouldBindJSON(&data); err != nil {
		for _, fieldErr := range err.(validator.ValidationErrors) {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": myValidator.FieldError{Err: fieldErr}.String(),
			})
			return
		}
	}

	couriersWithEta, err := h.Services.GetClosestCouriers(data.Origin, 1)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failure",
			"err":     err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    couriersWithEta,
	})

}
