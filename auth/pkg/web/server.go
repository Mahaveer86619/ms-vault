package web

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Mahaveer86619/ms/auth/pkg/config"
	"github.com/Mahaveer86619/ms/auth/pkg/db"
	"github.com/Mahaveer86619/ms/auth/pkg/handlers"
	mid "github.com/Mahaveer86619/ms/auth/pkg/middleware"
	"github.com/Mahaveer86619/ms/auth/pkg/services"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	authLimiter *mid.RateLimiter
	apiLimiter  *mid.RateLimiter
)

func initSystem() {
	config.InitConfig()
	db.InitDB()

	authLimiter = mid.NewRateLimiter(5, 30*time.Second) // 05 req in 30 sec
	apiLimiter = mid.NewRateLimiter(10, 1*time.Minute)  // 10 req in 1 min
}

func StartServer() {
	initSystem()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())
	e.Use(middleware.CORS())

	registerServices(e)

	serverAddress := fmt.Sprintf(":%s", config.GConfig.Port)
	if err := e.Start(serverAddress); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Error starting server: %v", err)
	}
}

func registerServices(e *echo.Echo) {

	// API group
	apiGroup := e.Group("/api/v1")
	authGroup := e.Group("/api/v1/auth")

	authGroup.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			rateLimitHandler := authLimiter.RateLimitMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { next(c) }))
			rateLimitHandler.ServeHTTP(c.Response(), c.Request())
			return nil
		}
	})
	apiGroup.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			rateLimitHandler := apiLimiter.RateLimitMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { next(c) }))
			rateLimitHandler.ServeHTTP(c.Response(), c.Request())
			return nil
		}
	})

	// Services
	healthService := services.NewHealthService()
	avatarService := services.NewAvatarService()
	authService := services.NewAuthService(avatarService)

	// Handlers
	handlers.NewHealthHandler(apiGroup, healthService)
	handlers.NewAvatarHandler(apiGroup, avatarService)
	handlers.NewAuthHandler(authGroup, authService)

}
