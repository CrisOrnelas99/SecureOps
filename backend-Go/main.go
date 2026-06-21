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
	if err := cfg.Validate(); err != nil {
		log.Fatalf("configuration validation failed: %v", err)
	}

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

	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(middleware.RequestContext())
	engine.Use(middleware.SecurityHeaders())
	engine.Use(middleware.GormMiddleware(gormDB))
	engine.Use(middleware.Cors(cfg.CorsAllowedOrigin))
	engine.Use(middleware.RequestFilter())
	// Register all routes centrally in the controller package
	controller.RegisterRoutes(engine, jwtManager, userRepository, controller.RouteHandlers{
		RegisterAuth:        authController.Register,
		LoginAuth:           authController.Login,
		GetAssets:           assetController.GetAssets,
		GetAsset:            assetController.GetAsset,
		CreateAsset:         assetController.CreateAsset,
		UpdateAsset:         assetController.UpdateAsset,
		DeleteAsset:         assetController.DeleteAsset,
		AssignVulnerability: assetController.AssignVulnerability,
		RemoveVulnerability: assetController.RemoveVulnerability,
		GetVulnerabilities:  vulnerabilityController.GetVulnerabilities,
		GetVulnerability:    vulnerabilityController.GetVulnerability,
		CreateVulnerability: vulnerabilityController.CreateVulnerability,
		UpdateVulnerability: vulnerabilityController.UpdateVulnerability,
		DeleteVulnerability: vulnerabilityController.DeleteVulnerability,
	})

	log.Printf("Go backend running on :%s", cfg.Port)
	if err := engine.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
