package main

import (
	"database/sql"
	"ta/internal/config"
	"ta/internal/handler"
	"ta/internal/middleware"
	"ta/internal/repository"
	"ta/internal/service"
	"ta/pkg/logger"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func main() {
	cfg := config.Load()
	log, err := logger.New()
	if err != nil {
		panic(err)
	}

	log.Info("starting server", zap.String("port", cfg.Port))

	db, err := sql.Open("postgres", cfg.DSN())
	if err != nil {
		log.Fatal("failed to open database connection",
			zap.Error(err),
			zap.String("dsn", cfg.DSN()))
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("failed to ping database",
			zap.Error(err),
			zap.String("host", cfg.DBHost),
			zap.String("port", cfg.DBPort),
			zap.String("database", cfg.DBName),
			zap.String("user", cfg.DBUser))
	}

	log.Info("database connection established")

	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)
	orderRepo := repository.NewOrderRepository(db)

	userService := service.NewUserService(userRepo)
	productService := service.NewProductService(productRepo)
	orderService := service.NewOrderService(orderRepo, productRepo, productService)

	userHandler := handler.NewUserHandler(userService)
	productHandler := handler.NewProductHandler(productService)
	orderHandler := handler.NewOrderHandler(orderService)

	r := gin.Default()
	r.Use(middleware.Tracer())
	r.Use(middleware.Logger(log))
	r.Use(middleware.Recovery(log))

	api := r.Group("/api/v1")
	{
		api.POST("/users", userHandler.Register)
		api.GET("/users/:id", userHandler.GetByID)

		api.POST("/products", productHandler.Create)
		api.GET("/products/:id", productHandler.GetByID)

		api.POST("/orders", orderHandler.Create)
		api.GET("/orders/:id", orderHandler.GetByID)
		api.GET("/users/:user_id/orders", orderHandler.GetByUserID)
	}

	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("failed to start server", zap.Error(err))
	}
}
