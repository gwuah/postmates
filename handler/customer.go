package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gwuah/postmates/database/models"
	"github.com/gwuah/postmates/lib/sms"
	"github.com/gwuah/postmates/utils"
	"gorm.io/gorm"
)

type CreateCustomerRequest struct {
	Phone string `json:"phone" validate:"required"`
}

type LoginRequest struct {
	Phone string `json:"phone" validate:"required"`
	Code  int    `json:"code" validate:"required"`
}

func (h *Handler) ListCustomers(c *gin.Context) {
	var customers []models.Customer

	if err := h.DB.Find(&customers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to retrieve customers",
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
			"message": "customer id not found",
		})
	}

	customer := new(models.Customer)
	result := h.DB.First(customer, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed To retrieve customer",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "success",
		"customer": customer,
	})
	return
}

func (h *Handler) sendSMS(customer models.Customer, token string) {
	response, err := h.SMS.SendTextMessage(sms.Message{
		To:  utils.GeneratePhoneNumber(customer.Phone),
		Sms: fmt.Sprintf("Your electra code: %s", token),
	})

	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(response)
}

func (h *Handler) SignupCustomer(c *gin.Context) {
	var data CreateCustomerRequest
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "failed to parse request",
		})
	}

	existingCustomer, err := h.Repo.FindCustomerByPhone(data.Phone)

	if err != nil && err != gorm.ErrRecordNotFound {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Request failed",
		})
		return
	}

	code := utils.GenerateOTP()

	if existingCustomer != nil {

		_, err = h.Repo.UpdateCustomer(existingCustomer.ID, map[string]interface{}{
			"Code": code,
		})

		go h.sendSMS(*existingCustomer, code)

		c.JSON(http.StatusOK, gin.H{
			"message": "customer already exists, token has been sent",
		})

	} else {

		record, err := h.Repo.CreateCustomerWithPhoneAndCode(data.Phone, utils.ConvertToInt(code))

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Customer Creation failed",
			})
			return
		}

		go h.sendSMS(*record, code)

		c.JSON(http.StatusOK, gin.H{
			"message": "new customer",
			"data": gin.H{
				"customer": record,
			},
		})

	}
}

func (h *Handler) LoginCustomer(c *gin.Context) {

	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "failed to parse request",
		})
		return
	}

	query := fmt.Sprintf("phone = '%s' AND code = '%d'", req.Phone, req.Code)
	customer, err := h.Repo.FindCustomerByQuery(query)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "wrong token or phone number",
		})
		return
	}

	token, err := h.JWT.GenerateToken(customer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to generate token",
		})
		return
	}

	refreshToken := h.Sec.Token(token)

	_, err = h.Repo.UpdateCustomer(customer.ID, map[string]interface{}{
		"Token": token,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to store refresh token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "success",
		"token":        token,
		"refreshToken": refreshToken,
	})
	return
}
