package main

import (
	"database/sql"
	"fmt"
	
	"subscription-service/internal/config"
	"subscription-service/internal/handlers"
	"subscription-service/internal/repository"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"go.uber.org/zap"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "subscription-service/docs"
)






func main() {
	
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	
	cfg, err := config.LoadConfig()
	if err != nil {
		sugar.Fatalf("Failed to load config: %v", err)
	}

	
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.Name)
	
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		sugar.Fatalf("Failed to open DB connection: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		sugar.Fatalf("Failed to ping DB: %v", err)
	}
	sugar.Info("Successfully connected to PostgreSQL")

	
	repo := repository.NewSubscriptionRepository(db)
	handler := handlers.NewHandler(repo, sugar)

	r := gin.Default()

	
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	
	api := r.Group("/api/v1")
	{
		api.POST("/subscriptions", handler.CreateSubscription)
		api.GET("/subscriptions", handler.GetAllSubscriptions) 
		api.GET("/subscriptions/:id", handler.GetSubscription)
		api.PUT("/subscriptions/:id", handler.UpdateSubscription)
		api.DELETE("/subscriptions/:id", handler.DeleteSubscription)
		
		api.GET("/subscriptions/cost", handler.GetTotalCost)
	}

	
	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	sugar.Infof("Starting server on %s", addr)
	if err := r.Run(addr); err != nil {
		sugar.Fatalf("Server failed: %v", err)
	}
}
