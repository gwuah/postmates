package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/gwuah/api/models"
)

type CreateCustomerRequest struct {
	Phone string `json:"phone" validate:"minLength:15"`
}

func (h *Handler) ListCustomers(c *gin.Context) {

}

func (h *Handler) ViewCustomer(c *gin.Context) {
	id := c.Param("id")
	customer := new(models.Customer)
	result := h.DB.First(customer, id)
	if result.Error != nil {
		c.JSON(500, gin.H{
			"message": "Error",
		})
		return
	}

	c.JSON(200, customer)
	return
}

func (h *Handler) CreateCustomer(c *gin.Context) {

	req := new(CreateCustomerRequest)
	if c.BindJSON(req) != nil {
		c.JSON(500, gin.H{
			"message": "Error",
		})
		return
	}

	customer := models.Customer{
		Phone: req.Phone,
	}
	result := h.DB.Create(&customer)
	if result.RowsAffected > 0 {
		c.JSON(200, gin.H{
			"message":  "success",
			"customer": customer,
		})
		return
	}

}
