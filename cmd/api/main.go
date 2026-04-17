package main

import (
	"log"
	"os"

	"super-indo-api/pkg/infrastructure/postgres"

	"github.com/gin-gonic/gin"
)

func main() {
	app := Bootstrap()
	defer app.DB.Close()

	if len(os.Args) > 1 && os.Args[1] == "seed" {
		log.Println("menjalankan seeder...")
		postgres.SeedProducts(app.DB)
		log.Println("seeder selesai")
		return
	}

	r := gin.Default()

	app.AuthHandler.RegisterRoutes(r.Group("/auth"))
	app.ProductHandler.RegisterRoutes(r.Group("/product"), app.AuthMiddleware)

	port := app.Config.AppPort
	if port == "" {
		port = "8080"
	}

	log.Printf("server berjalan di port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("gagal menjalankan server: %v", err)
	}
}
