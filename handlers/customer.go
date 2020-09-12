package handler

import (
	"github.com/gin-gonic/gin"
)

type CreateCustomerRequest struct {
	Phone string `json:"phone"`
}

func (h *Handler) ListCustomers(c *gin.Context) {

}

func (h *Handler) ViewCustomer(c *gin.Context) {

}

func (h *Handler) CreateCustomer(c *gin.Context) {

	// var data

	// if c.BindJSON(&data) != nil {
	// 	c.JSON(500, gin.H{
	// 		"message": "Error",
	// 	})
	// 	return
	// }
}
