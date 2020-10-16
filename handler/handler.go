package handler

import (
	"os"

	"github.com/electra-systems/core-api/lib/billing"
	"github.com/electra-systems/core-api/lib/eta"
	"github.com/electra-systems/core-api/lib/sms"
	"github.com/electra-systems/core-api/lib/ws"
	"github.com/electra-systems/core-api/repository"
	"github.com/electra-systems/core-api/services"
	"github.com/electra-systems/core-api/utils/jwt"
	"github.com/electra-systems/core-api/utils/secure"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

type Handler struct {
	DB                   *gorm.DB
	Repo                 *repository.Repository
	JWT                  jwt.Service
	Sec                  *secure.Service
	Services             *services.Services
	maxMessageTypeLength int
	Hub                  *ws.Hub
	RedisDB              *redis.Client
	SMS                  *sms.SMS
	Eta                  *eta.Eta
}

func New(DB *gorm.DB, jwt jwt.Service, sec *secure.Service, redisDB *redis.Client) *Handler {
	SMS := sms.New(os.Getenv("TERMII_API_KEY"))
	eta := eta.New(os.Getenv("GMAPS_TOKEN"))
	billing := billing.New()
	hub := ws.NewHub()
	go hub.Run()
	repo := repository.New(DB, redisDB)
	services := services.New(repo, eta, hub, billing)

	return &Handler{
		DB:                   DB,
		Repo:                 repo,
		JWT:                  jwt,
		Services:             services,
		maxMessageTypeLength: 30,
		Hub:                  hub,
		RedisDB:              redisDB,
		SMS:                  SMS,
		Eta:                  eta,
	}
}

func (h *Handler) Register(v1 *gin.RouterGroup) {
	v1.GET("/customer/realtime/:id", h.handleConnection("customer"))
	v1.GET("/courier/realtime/:id", h.handleConnection("courier"))
	v1.POST("/get-closest-couriers", h.GetClosestCouriers)
	v1.POST("/get-delivery-cost", h.GetDeliveryCost)
	v1.POST("/customer-rate-trip", h.handleCustomerRating)
	v1.POST("/courier-rate-trip", h.handleCourierRating)

	v1.POST("/signup", h.Signup)
	v1.POST("/login", h.Login)
	v1.POST("/otp/verify", h.VerifyOTP)
	v1.GET("/refresh/:token", h.Refresh)

	//  middleware.JWT(h.JWT)
	customers := v1.Group("/customers")
	customers.GET("/", h.ListCustomers)
	customers.GET("/:id", h.ViewCustomer)
	customers.POST("/", h.CreateCustomer)

}
