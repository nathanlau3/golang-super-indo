package main

import (
	"log"
	"os"

	"super-indo-api/pkg/config"
	producthandler "super-indo-api/api/product"
	"super-indo-api/internal/product/usecase"
	"super-indo-api/pkg/infrastructure/adapter"
	"super-indo-api/pkg/infrastructure/postgres"
	"super-indo-api/pkg/infrastructure/redis"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	// koneksi ke database
	pgDB, err := postgres.NewConnection(cfg)
	if err != nil {
		log.Fatalf("gagal koneksi ke database: %v", err)
	}
	defer pgDB.Close()

	// jalankan migration
	if err := postgres.RunMigration(pgDB, postgres.MigrationSQL); err != nil {
		log.Fatalf("gagal menjalankan migration: %v", err)
	}

	// koneksi redis
	rdb := redis.NewConnection(cfg)

	// dependency injection — composition root
	productRepo := adapter.NewProductRepository(pgDB, rdb)

	createProduct := usecase.NewCreateProduct(productRepo)
	getProducts := usecase.NewGetProducts(productRepo)
	getProductByID := usecase.NewGetProductByID(productRepo)

	productHandler := producthandler.NewProductHandler(createProduct, getProducts, getProductByID)

	// handle sub-command
	if len(os.Args) > 1 && os.Args[1] == "seed" {
		log.Println("menjalankan seeder...")
		postgres.SeedProducts(pgDB)
		log.Println("seeder selesai")
		return
	}

	// setup router — tiap module register route sendiri
	r := gin.Default()
	productHandler.RegisterRoutes(r.Group("/product"))

	port := cfg.AppPort
	if port == "" {
		port = "8080"
	}

	log.Printf("server berjalan di port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("gagal menjalankan server: %v", err)
	}
}
