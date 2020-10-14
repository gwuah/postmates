package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/electra-systems/core-api/database/models"
	"github.com/electra-systems/core-api/lib/sms"
	"github.com/electra-systems/core-api/utils"
	myValidator "github.com/electra-systems/core-api/utils/validator"
	"gorm.io/gorm"
)

type CreateCustomerRequest struct {
	Phone string `json:"phone" validate:"required"`
}

func (h *Handler) ListCustomers(c *gin.Context) {
	var customers []models.Customer

	if err := h.DB.Find(&customers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed To Retrieve Customers",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "success",
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
			"message": "failed To Retrieve Customer",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "success",
		"customer": customer,
	})
	return
}

func (h *Handler) sendSMS(customer models.Customer) {
	response, err := h.SMS.SendTextMessage(sms.Message{
		To:  utils.GeneratePhoneNumber(customer.Phone),
		Sms: fmt.Sprintf("Your electra code: %s", utils.GenerateOTP()),
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
			"message": "Request failed",
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
				"message": "Customer Creation failed",
			})
			return
		}

		go h.sendSMS(*record)

		c.JSON(http.StatusOK, gin.H{
			"message": "New Customer",
		})

	}
}
