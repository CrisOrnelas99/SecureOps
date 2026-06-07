package main

import (
	"log"
	"net/http"

	"risk-service-go/api/config"
	"risk-service-go/api/controller"
	"risk-service-go/api/service"
)

func main() {
	cfg := config.Load()
	riskService := service.NewRiskService()
	riskController := controller.NewRiskController(riskService)

	http.HandleFunc("/health", controller.Health)
	http.HandleFunc("/calculate-risk", riskController.CalculateRisk)

	log.Printf("Risk service running on :%s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, nil))
}
