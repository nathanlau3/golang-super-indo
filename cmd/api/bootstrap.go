package main

import (
	"database/sql"
	"log"

	"super-indo-api/pkg/config"
	authhandler "super-indo-api/api/auth"
	producthandler "super-indo-api/api/product"
	authusecase "super-indo-api/internal/auth/usecase"
	productusecase "super-indo-api/internal/product/usecase"
	"super-indo-api/pkg/infrastructure/adapter"
	"super-indo-api/pkg/infrastructure/jwt"
	"super-indo-api/pkg/infrastructure/postgres"
	"super-indo-api/pkg/infrastructure/redis"
	"super-indo-api/pkg/middleware"

	"github.com/gin-gonic/gin"
	goredis "github.com/redis/go-redis/v9"
)

type App struct {
	Config         *config.Config
	DB             *sql.DB
	Redis          *goredis.Client
	JWT            *jwt.JWTService
	ProductHandler *producthandler.ProductHandler
	AuthHandler    *authhandler.AuthHandler
	AuthMiddleware gin.HandlerFunc
}

func Bootstrap() *App {
	cfg := config.LoadConfig()

	pgDB, err := postgres.NewConnection(cfg)
	if err != nil {
		log.Fatalf("gagal koneksi ke database: %v", err)
	}

	if err := postgres.RunMigration(pgDB, postgres.MigrationSQL); err != nil {
		log.Fatalf("gagal menjalankan migration products: %v", err)
	}
	if err := postgres.RunMigration(pgDB, postgres.UserMigrationSQL); err != nil {
		log.Fatalf("gagal menjalankan migration users: %v", err)
	}

	rdb := redis.NewConnection(cfg)

	jwtSvc, err := jwt.NewJWTService(cfg.JWTPrivateKeyPath, cfg.JWTPublicKeyPath)
	if err != nil {
		log.Fatalf("gagal inisialisasi JWT service: %v", err)
	}

	productRepo := adapter.NewProductRepository(pgDB, rdb)
	userRepo := adapter.NewUserRepository(pgDB)

	createProduct := productusecase.NewCreateProduct(productRepo)
	getProducts := productusecase.NewGetProducts(productRepo)
	getProductByID := productusecase.NewGetProductByID(productRepo)

	register := authusecase.NewRegister(userRepo)
	login := authusecase.NewLogin(userRepo, jwtSvc)

	return &App{
		Config:         cfg,
		DB:             pgDB,
		Redis:          rdb,
		JWT:            jwtSvc,
		ProductHandler: producthandler.NewProductHandler(createProduct, getProducts, getProductByID),
		AuthHandler:    authhandler.NewAuthHandler(register, login),
		AuthMiddleware: middleware.AuthMiddleware(jwtSvc),
	}
}
