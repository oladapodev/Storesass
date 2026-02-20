// @title Storefront SaaS API
// @version 1.0
// @description A multi-tenant storefront SaaS backend API
// @host localhost:8080
// @BasePath /api/v1
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/yourusername/storefront-saas-go/internal/config"
	"github.com/yourusername/storefront-saas-go/internal/db"
	_ "github.com/yourusername/storefront-saas-go/docs"
	"github.com/yourusername/storefront-saas-go/internal/handler"
	"github.com/yourusername/storefront-saas-go/internal/middleware"
	"github.com/yourusername/storefront-saas-go/internal/repository"
	"github.com/yourusername/storefront-saas-go/internal/service"
)

func main() {
	cfg := config.LoadConfig()

	database, err := db.InitDB(cfg)
	if err != nil {
		log.Fatalf("failed to init db: %v", err)
	}

	redisClient, err := db.InitRedis(cfg)
	if err != nil {
		log.Printf("redis not available: %v", err)
	}

	storeRepo := repository.NewStoreRepository(database)
	productRepo := repository.NewProductRepository(database)

	storeSvc := service.NewStoreService(storeRepo)
	productSvc := service.NewProductService(productRepo, redisClient)

	storeHandler := handler.NewStoreHandler(storeSvc)
	productHandler := handler.NewProductHandler(productSvc, storeSvc)

	router := gin.Default()
	router.Use(middleware.CORSMiddleware())

	v1 := router.Group("/api/v1")
	{
		stores := v1.Group("/stores")
		{
			stores.GET("", storeHandler.ListActiveStores)
			stores.GET("/:slug", storeHandler.GetStoreBySlug)
			stores.POST("", storeHandler.CreateStore)
			stores.GET("/:slug/products", productHandler.ListProductsByStore)
		}

		products := v1.Group("/products")
		{
			products.GET("", productHandler.ListHotProducts)
			products.GET("/search", productHandler.SearchActiveProducts)
			products.GET("/:id", productHandler.GetProductByID)
			products.POST("", productHandler.CreateProduct)
		}
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Serve frontend static files
	router.StaticFile("/", "./web/dist/index.html")
	router.Static("/assets", "./web/dist/assets")
	router.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		method := c.Request.Method
		if (method == http.MethodGet || method == http.MethodHead) &&
			!strings.HasPrefix(path, "/api/") &&
			!strings.HasPrefix(path, "/swagger/") {
			c.File("./web/dist/index.html")
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	})

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: router,
	}

	go func() {
		log.Printf("Starting server on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}

	log.Println("Server exited cleanly")
}
