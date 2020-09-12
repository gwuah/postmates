package handler

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	DB *gorm.DB
}

func New(db *gorm.DB) *Handler {
	return &Handler{db}
}

func (h *Handler) Register(v1 *gin.RouterGroup) {

	customers := v1.Group("/customers")
	customers.GET("/", h.ListCustomers)
	customers.GET("/:id", h.ViewCustomer)
	customers.POST("/", h.CreateCustomer)

	deliveries := v1.Group("/deliveries")
	deliveries.GET("/", h.ListDeliveries)
	deliveries.GET("/:id", h.ViewDeliveries)
	deliveries.POST("/", h.CreateDelivery)

	orders := v1.Group("/orders")
	orders.GET("/", h.ListOrders)
	orders.GET("/:id", h.ViewOrder)
	orders.POST("/", h.CreateOrder)

	electrons := v1.Group("/electrons")
	electrons.GET("/", h.ListElectrons)
	electrons.GET("/:id", h.ViewElectron)
	electrons.POST("/", h.CreateElectron)

}
