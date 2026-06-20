package main

import (
	"context"
	"log"
	"time"

	"github.com/gin-gonic/gin"

	"secureops/backend-go/api/config"
	"secureops/backend-go/api/controller"
	controllerasset "secureops/backend-go/api/controller/asset"
	controllerauth "secureops/backend-go/api/controller/auth"
	controllervulnerability "secureops/backend-go/api/controller/vulnerability"
	"secureops/backend-go/api/middleware"
	repositoryasset "secureops/backend-go/api/repository/asset"
	repositoryuser "secureops/backend-go/api/repository/user"
	repositoryvulnerability "secureops/backend-go/api/repository/vulnerability"
	"secureops/backend-go/api/security"
	serviceasset "secureops/backend-go/api/service/asset"
	serviceauth "secureops/backend-go/api/service/auth"
	servicevulnerability "secureops/backend-go/api/service/vulnerability"
	"secureops/backend-go/api/utils"
)

func main() {
	cfg := config.Load()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	gormDB, err := utils.Connect(ctx, cfg)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	defer func() {
		if err := utils.Close(gormDB); err != nil {
			log.Printf("database close failed: %v", err)
		}
	}()

	if err := utils.RunMigrations(ctx, gormDB); err != nil {
		log.Fatalf("database migration failed: %v", err)
	}

	jwtManager := security.NewJWTManager(cfg.JWTSecret, cfg.JWTExpiration, cfg.JWTIssuer, cfg.JWTAudience)

	userRepository := repositoryuser.NewUserRepository(gormDB)
	assetRepository := repositoryasset.NewAssetRepository(gormDB)
	vulnerabilityRepository := repositoryvulnerability.NewVulnerabilityRepository(gormDB)
	authService := serviceauth.NewAuthService(jwtManager, userRepository)
	assetService := serviceasset.NewAssetService(assetRepository)
	vulnerabilityService := servicevulnerability.NewVulnerabilityService(vulnerabilityRepository)

	authController := controllerauth.NewAuthController(authService)
	assetController := controllerasset.NewAssetController(assetService)
	vulnerabilityController := controllervulnerability.NewVulnerabilityController(vulnerabilityService)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.RequestContext())
	router.Use(middleware.SecurityHeaders())
	router.Use(middleware.GormMiddleware(gormDB))
	router.Use(config.CorsConfig())
	router.Use(middleware.RequestFilter())
	// Register all routes centrally in the controller package
	controller.RegisterRoutes(router, jwtManager, userRepository, authController, assetController, vulnerabilityController)

	log.Printf("Go backend running on :%s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
