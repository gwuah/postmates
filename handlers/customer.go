package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/gwuah/api/models"
	myValidator "github.com/gwuah/api/utils/validator"
)

// CreateCustomerRequest - create customer request object
type CreateCustomerRequest struct {
	Phone string `json:"phone" validate:"required"`
}

// ListCustomers - returns list of customers
func (h *Handler) ListCustomers(c *gin.Context) {
	var customers []models.Customer

	if err := h.DB.Find(&customers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed To Retrieve Customers",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Success",
		"customers": customers,
	})
}

// ViewCustomer - returns a single customer by ID
func (h *Handler) ViewCustomer(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "customer id not found",
		})
	}

	customer := new(models.Customer)
	result := h.DB.First(customer, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed To Retrieve Customer",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Success",
		"customer": customer,
	})
	return
}

// CreateCustomer - creates a new customer and returns it
func (h *Handler) CreateCustomer(c *gin.Context) {
	req := new(CreateCustomerRequest)
	if err := c.ShouldBindJSON(req); err != nil {
		for _, fieldErr := range err.(validator.ValidationErrors) {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": myValidator.FieldError{Err: fieldErr}.String(),
			})
			return
		}
	}

	customer := models.Customer{
		Phone: req.Phone,
	}

	result := h.DB.Create(&customer)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed To Retrieve Customer",
		})
		return
	}

	if result.RowsAffected > 0 {
		c.JSON(http.StatusOK, gin.H{
			"message":  "Success",
			"customer": customer,
		})
		return
	}

}
