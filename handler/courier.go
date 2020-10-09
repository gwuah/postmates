package handler

import (
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/gwuah/api/shared"
	myValidator "github.com/gwuah/api/utils/validator"
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

	ids := h.Services.GetClosestCouriers(data.Origin, 2)
	couriers, err := h.Repo.GetAllCouriers(ids)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failure",
			"err":     err,
		})
		return
	}

	var couriersWithEta []shared.CourierWithEta

	for _, courier := range couriers {
		duration, distance, err := h.Eta.GetDistanceAndDuration(shared.Coord{
			Latitude:  courier.Latitude,
			Longitude: courier.Longitude,
		}, shared.Coord{
			Latitude:  data.Origin.Latitude,
			Longitude: data.Origin.Longitude,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "failure",
				"err":     err,
			})
			return
		}

		couriersWithEta = append(couriersWithEta, shared.CourierWithEta{
			Courier:  courier,
			Distance: float64(distance),
			Duration: float64(duration),
		})

		sort.Slice(couriersWithEta, func(i, j int) bool {
			return couriersWithEta[i].Duration < couriersWithEta[j].Duration
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    couriersWithEta,
	})

}
