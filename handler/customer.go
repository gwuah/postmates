package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/gwuah/api/database/models"
	"github.com/gwuah/api/lib/otp"
	"github.com/gwuah/api/lib/sms"
	myValidator "github.com/gwuah/api/utils/validator"
	"gorm.io/gorm"
)

type CreateCustomerRequest struct {
	Phone string `json:"phone" validate:"required"`
}

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

func (h *Handler) ViewCustomer(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Customer ID not found",
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

func (h *Handler) sendSMS(customer models.Customer) {
	response, err := h.SMS.SendTextMessage(sms.Message{
		To:  customer.Phone,
		Sms: fmt.Sprintf("Your electra code: %d", otp.GenerateOTP()),
	})

	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(response)
}

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

	existingCustomer, err := h.Repo.FindCustomerByPhone(req.Phone)

	if err != nil && err != gorm.ErrRecordNotFound {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Request Failed",
		})
		return
	}

	if existingCustomer != nil {
		go h.sendSMS(*existingCustomer)
		c.JSON(http.StatusOK, gin.H{
			"message": "Customer Exists Already",
		})

	} else {

		record, err := h.Repo.CreateCustomerWithPhone(req.Phone)

		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Customer Creation Failed",
			})
			return
		}

		go h.sendSMS(*record)

		c.JSON(http.StatusOK, gin.H{
			"message": "New Customer",
		})

	}
}
