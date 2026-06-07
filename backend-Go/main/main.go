package main

import (
	"context"
	"log"
	"time"

	"github.com/gin-gonic/gin"

	"secureops/backend-go/api/config"
	"secureops/backend-go/api/controller"
	"secureops/backend-go/api/database"
	"secureops/backend-go/api/middleware"
	"secureops/backend-go/api/repository"
	"secureops/backend-go/api/security"
	"secureops/backend-go/api/service"
)

func main() {
	cfg := config.Load()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	gormDB, err := database.Connect(ctx, cfg)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	defer func() {
		if err := database.Close(gormDB); err != nil {
			log.Printf("database close failed: %v", err)
		}
	}()

	if err := database.EnsureSchema(ctx, gormDB); err != nil {
		log.Fatalf("database schema setup failed: %v", err)
	}

	jwtService := security.NewJwtService(cfg.JWTSecret, cfg.JWTExpiration)

	userRepository := repository.NewUserRepository(gormDB)
	assetRepository := repository.NewAssetRepository(gormDB)
	vulnerabilityRepository := repository.NewVulnerabilityRepository(gormDB)
	wafEventRepository := repository.NewWafEventRepository(gormDB)

	restClient := config.RestClientConfig(cfg)
	authService := service.NewAuthService(userRepository, jwtService)
	assetRiskService := service.NewAssetRiskService(assetRepository)
	assetService := service.NewAssetService(assetRepository, vulnerabilityRepository, restClient, assetRiskService)
	vulnerabilityService := service.NewVulnerabilityService(vulnerabilityRepository)

	authController := controller.NewAuthController(authService)
	assetController := controller.NewAssetController(assetService)
	vulnerabilityController := controller.NewVulnerabilityController(vulnerabilityService)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(config.CorsConfig())
	router.Use(middleware.WafFilter(wafEventRepository))

	router.GET("/api/health", controller.Health)
	router.POST("/api/auth/register", authController.Register)
	router.POST("/api/auth/login", authController.Login)

	protected := router.Group("/api")
	protected.Use(config.SecurityConfig(jwtService, userRepository))
	{
		protected.GET("/test/secure", controller.SecureTest)

		protected.GET("/assets", assetController.GetAssets)
		protected.GET("/assets/:id", assetController.GetAsset)
		protected.POST("/assets", assetController.CreateAsset)
		protected.PUT("/assets/:id", assetController.UpdateAsset)
		protected.DELETE("/assets/:id", assetController.DeleteAsset)
		protected.POST("/assets/:id/vulnerabilities/:vulnerabilityId", assetController.AssignVulnerability)
		protected.DELETE("/assets/:id/vulnerabilities/:vulnerabilityId", assetController.RemoveVulnerability)
		protected.POST("/assets/:id/calculate-risk", assetController.CalculateRisk)

		protected.GET("/vulnerabilities", vulnerabilityController.GetVulnerabilities)
		protected.GET("/vulnerabilities/:id", vulnerabilityController.GetVulnerability)
		protected.POST("/vulnerabilities", vulnerabilityController.CreateVulnerability)
		protected.PUT("/vulnerabilities/:id", vulnerabilityController.UpdateVulnerability)
		protected.DELETE("/vulnerabilities/:id", vulnerabilityController.DeleteVulnerability)
	}

	log.Printf("Go backend running on :%s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
