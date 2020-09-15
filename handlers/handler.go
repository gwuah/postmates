package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/gwuah/api/middleware"
	"github.com/gwuah/api/repository"
	"github.com/gwuah/api/utils/jwt"
	"github.com/gwuah/api/utils/secure"
	"github.com/gwuah/api/wss"
	"gorm.io/gorm"
)

type Handler struct {
	DB   *gorm.DB
	Repo *repository.Repository
	JWT  jwt.Service
	Sec  *secure.Service
}

func New(DB *gorm.DB, jwt jwt.Service, sec *secure.Service) *Handler {
	repo := repository.New(DB)
	return &Handler{DB, repo, jwt, sec}
}

func (h *Handler) Register(v1 *gin.RouterGroup) {

	wss := wss.New()

	v1.GET("/customer/realtime/:id", wss.HandleWebsocketConnection("customer"))
	v1.GET("/electron/realtime/:id", wss.HandleWebsocketConnection("electron"))

	v1.POST("/signup", h.Signup)
	v1.POST("/login", h.Login)
	v1.POST("/otp/verify", h.VerifyOTP)
	v1.GET("/refresh/:token", h.Refresh)

	customers := v1.Group("/customers", middleware.JWT(h.JWT))
	customers.GET("/", h.ListCustomers)
	customers.GET("/:id", h.ViewCustomer)
	customers.POST("/", h.CreateCustomer)

	deliveries := v1.Group("/deliveries", middleware.JWT(h.JWT))
	deliveries.GET("/", h.ListDeliveries)
	deliveries.GET("/:id", h.ViewDelivery)
	deliveries.POST("/", h.CreateDelivery)

	orders := v1.Group("/orders", middleware.JWT(h.JWT))
	orders.GET("/", h.ListOrders)
	orders.GET("/:id", h.ViewOrder)
	orders.POST("/", h.CreateOrder)

	electrons := v1.Group("/electrons", middleware.JWT(h.JWT))
	electrons.GET("/", h.ListElectrons)
	electrons.GET("/:id", h.ViewElectron)
	electrons.POST("/", h.CreateElectron)

}
