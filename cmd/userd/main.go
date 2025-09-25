package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/your-org/custos/internal/application/usecase/auth"
	"github.com/your-org/custos/internal/config"
	"github.com/your-org/custos/internal/domain/service"
	"github.com/your-org/custos/internal/infrastructure/persistence/mysql"
	"github.com/your-org/custos/internal/interface/http/handler"
	"github.com/your-org/custos/internal/interface/http/middleware"
	"github.com/your-org/custos/internal/interface/http/router"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := mysql.NewDatabase(cfg.Database.DSN(), cfg.App.Env == "development")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.AutoMigrate(); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	userRepo := mysql.NewUserRepository(db.DB())

	tokenService := service.NewTokenService(cfg.JWT.SecretKey, cfg.JWT.AccessTokenTTL)
	authService := service.NewAuthService(userRepo, tokenService)

	registerUC := auth.NewRegisterUseCase(authService)
	loginUC := auth.NewLoginUseCase(authService)

	authHandler := handler.NewAuthHandler(registerUC, loginUC)
	userHandler := handler.NewUserHandler()
	healthHandler := handler.NewHealthHandler()
	authMW := middleware.NewAuthMiddleware(tokenService)

	routerHandler := router.NewRouter(authHandler, userHandler, healthHandler, authMW)
	ginEngine := routerHandler.SetupRoutes()

	srv := &http.Server{
		Addr:         ":" + cfg.App.Port,
		Handler:      ginEngine,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	go func() {
		log.Printf("Server starting on port %s", cfg.App.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}