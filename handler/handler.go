package handler

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/gwuah/api/lib/sms"
	"github.com/gwuah/api/lib/ws"
	"github.com/gwuah/api/repository"
	"github.com/gwuah/api/services"
	"github.com/gwuah/api/utils/jwt"
	"github.com/gwuah/api/utils/secure"
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
}

func New(DB *gorm.DB, jwt jwt.Service, sec *secure.Service, redisDB *redis.Client) *Handler {
	repo := repository.New(DB, redisDB)
	services := services.New(repo)
	hub := ws.NewHub()
	SMS := sms.New(os.Getenv("TERMII_API_KEY"))
	go hub.Run()

	return &Handler{
		DB:                   DB,
		Repo:                 repo,
		JWT:                  jwt,
		Services:             services,
		maxMessageTypeLength: 30,
		Hub:                  hub,
		RedisDB:              redisDB,
		SMS:                  SMS,
	}
}

func (h *Handler) Register(v1 *gin.RouterGroup) {

	v1.GET("/customer/realtime/:id", h.handleConnection("customer"))
	v1.GET("/electron/realtime/:id", h.handleConnection("electron"))

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
