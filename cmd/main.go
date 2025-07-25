package main

import (
	"log"

	"github.com/Vladislav5108/inventorytracker/configs"
	"github.com/Vladislav5108/inventorytracker/internal/infrastructure/db/postgres"
	"github.com/Vladislav5108/inventorytracker/internal/transport/myhttp"
	"github.com/Vladislav5108/inventorytracker/internal/usecase"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := configs.Load()
	if err != nil {
		log.Fatal("Config error:", err)
	}

	log.Println("Server port:", cfg.Server.Port)
	log.Println("DB host:", cfg.Postgres.Host)

	db, err := postgres.NewDB(postgres.Config{
		Host:     cfg.Postgres.Host,
		Port:     cfg.Postgres.Port,
		User:     cfg.Postgres.User,
		Password: cfg.Postgres.Password,
		DBName:   cfg.Postgres.DBName,
		SSLMode:  cfg.Postgres.SSLMode,
		Timeout:  cfg.Postgres.Timeout,
	})
	if err != nil {
		log.Fatal("DB Connection error:", err)
	}
	defer db.Close()

	log.Println("Postgres connection!")

	ProductRepo := postgres.NewProductRepo(db)
	CategoryRepo := postgres.NewCategoryRepo(db)

	log.Println("Repositories initialized")

	ProductUsecase := usecase.NewProductUseCase(ProductRepo)
	CategoryUsecase := usecase.NewCategoryUseCase(CategoryRepo)
	ProductAdminUsecase := usecase.NewAdminProductUseCase(ProductRepo, ProductRepo)

	productHandler := myhttp.NewProductHandler(ProductUsecase)
	productAdminHandler := myhttp.NewAdminProductHandler(ProductAdminUsecase)
	categoryHandler := myhttp.NewCategoryHandler(CategoryUsecase)

	r := gin.Default()

	r.GET("/products", productHandler.GetAll)
	r.GET("/products/:id", productHandler.GetByID)
	r.GET("/categories", categoryHandler.GetAllCategories)

	admin := r.Group("/admin")
	{
		admin.POST("/products", productAdminHandler.Add)
		admin.PUT("/products/:id", productAdminHandler.Update)
	}

	port := cfg.Server.Port
	log.Printf("server rгn on port %s", port)
	if err := r.Run(port); err != nil {
		log.Fatal("Server startup error")
	}
}
