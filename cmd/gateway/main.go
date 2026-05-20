package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joaodddev/go-auth-gateway/internal/auth"
	"github.com/joaodddev/go-auth-gateway/internal/config"
	"github.com/joaodddev/go-auth-gateway/internal/handler"
	"github.com/joaodddev/go-auth-gateway/internal/middleware"
	"github.com/joaodddev/go-auth-gateway/internal/proxy"
	"github.com/joaodddev/go-auth-gateway/pkg/logger"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}

	// Initialize logger
	if err := logger.Init("info"); err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}
	defer logger.Sync()

	// Initialize Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       0,
	})

	// Test Redis connection
	ctx := context.Background()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		logger.Log.Warn("Redis connection failed", zap.Error(err))
	} else {
		logger.Log.Info("Redis connected successfully")
	}

	// Initialize JWT manager
	jwtManager := auth.NewJWTManager(cfg.JWT.Secret, cfg.JWT.ExpirationHours)

	// Initialize auth middleware
	authMiddleware := auth.NewAuthMiddleware(jwtManager)

	// Initialize rate limiter
	rateLimiter := middleware.NewRateLimiter(redisClient, cfg.RateLimit.Requests, cfg.RateLimit.Duration)

	// Initialize reverse proxy
	serviceURLs := map[string]string{
		"users":    cfg.Services.UserServiceURL,
		"orders":   cfg.Services.OrderServiceURL,
		"products": cfg.Services.ProductServiceURL,
	}

	reverseProxy, err := proxy.NewReverseProxy(serviceURLs)
	if err != nil {
		logger.Log.Fatal("Failed to create reverse proxy", zap.Error(err))
	}

	// Initialize handlers
	authHandler := handler.NewAuthHandler(jwtManager)

	// Setup routes
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	})

	// Auth routes (public)
	mux.HandleFunc("/login", authHandler.Login)
	mux.HandleFunc("/validate", authMiddleware.Authenticate(authHandler.Validate))

	// Protected service routes
	mux.HandleFunc("/users/", authMiddleware.Authenticate(
		rateLimiter.Limit(reverseProxy.Route("users")),
	))
	mux.HandleFunc("/orders/", authMiddleware.Authenticate(
		rateLimiter.Limit(reverseProxy.Route("orders")),
	))
	mux.HandleFunc("/products/", authMiddleware.Authenticate(
		rateLimiter.Limit(reverseProxy.Route("products")),
	))

	// Apply CORS middleware globally
	handler := middleware.CORS(mux)

	// Create server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Server.Port),
		Handler:      handler,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Graceful shutdown
	go func() {
		logger.Log.Info("Server starting", zap.String("port", cfg.Server.Port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Fatal("Server failed", zap.Error(err))
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Log.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Log.Info("Server exited properly")
}
