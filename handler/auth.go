package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/gwuah/api/database/models"
	myValidator "github.com/gwuah/api/utils/validator"
)

type SignupRequest struct {
	Phone string `json:"phone" validate:"required"`
}

type LoginRequest struct {
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Password string `json:"password" validate:"required"`
}

type VerifyOTPRequest struct {
	Phone string `json:"phone" validate:"required"`
	Code  int    `json:"code" validate:"required"`
}

func (h *Handler) Signup(c *gin.Context) {
	req := new(SignupRequest)
	if err := c.ShouldBindJSON(req); err != nil {
		for _, fieldErr := range err.(validator.ValidationErrors) {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": myValidator.FieldError{Err: fieldErr}.String(),
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "failed to parse request",
		})
		return
	}

	customer := models.Customer{
		Phone: req.Phone,
	}

	result := h.DB.Create(&customer)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed To Signup Customer",
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

func (h *Handler) Login(c *gin.Context) {
	req := new(LoginRequest)
	if err := c.ShouldBindJSON(req); err != nil {
		for _, fieldErr := range err.(validator.ValidationErrors) {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": myValidator.FieldError{Err: fieldErr}.String(),
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "failed to parse request",
		})
		return
	}

	if req.Phone == "" && req.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "validation failed, either phone or email not found",
		})
		return
	}

	customer := new(models.Customer)

	if req.Phone != "" && req.Email == "" {
		result := h.DB.Where("phone = ?", req.Phone).First(customer)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "failed to fetch customer",
			})
			return
		}
	}

	if req.Email != "" && req.Phone == "" {
		result := h.DB.Where("email = ?", req.Email).First(customer)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "failed to fetch customer",
			})
			return
		}
	}

	token, err := h.JWT.GenerateToken(customer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to generate token",
		})
		return
	}

	refreshToken := h.Sec.Token(token)

	result := h.DB.Model(customer).Update("token", refreshToken)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to store refresh token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Success",
		"token":        token,
		"refreshToken": refreshToken,
	})
	return
}

func (h *Handler) VerifyOTP(c *gin.Context) {
	req := new(VerifyOTPRequest)
	if err := c.ShouldBindJSON(req); err != nil {
		for _, fieldErr := range err.(validator.ValidationErrors) {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": myValidator.FieldError{Err: fieldErr}.String(),
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "failed to parse request",
		})
		return
	}

	customer := new(models.Customer)
	result := h.DB.First(customer)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to retrieve customer",
		})
		return
	}

	if req.Code != customer.Code {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "code does not match",
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

	result = h.DB.Model(customer).Updates(map[string]interface{}{
		"active": true, "token": refreshToken,
	})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to activate customer",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Success",
		"token":        token,
		"refreshToken": refreshToken,
	})
	return
}

func (h *Handler) Refresh(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "token not found",
		})
		return
	}

	customer := new(models.Customer)
	result := h.DB.Where("token = ?", token).First(customer)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to fetch customer",
		})
		return
	}

	newToken, err := h.JWT.GenerateToken(customer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to generate token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
		"token":   newToken,
	})
	return
}
