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
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("database close error:%v", err)
		}
	}()

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

	public := r.Group("/api/v1")

	public.GET("/products", productHandler.GetAll)
	public.GET("/products/:id", productHandler.GetByID)
	public.GET("/products/category/:id", productHandler.GetByCategory)

	public.GET("/categories", categoryHandler.GetAllCategories)
	public.GET("categories/:id", categoryHandler.GetByIDCategory)

	admin := r.Group("/admin")
	{
		admin.POST("/products", productAdminHandler.Add)
		admin.PUT("/products/:id", productAdminHandler.Update)
		admin.PATCH("/products/:id/archive", productAdminHandler.Archiv)
		admin.PATCH("/products/:id/restore", productAdminHandler.Restore)
		admin.GET("/products/archived", productAdminHandler.GetArchived)

		admin.POST("/categories", categoryHandler.CreateCategory)
		admin.PUT("/categories/:id", categoryHandler.UpDateCategory)
		admin.DELETE("/categories/:id", categoryHandler.DeleteCategory)
		admin.GET("/categories/:id", categoryHandler.GetByIDCategory)

	}

	port := cfg.Server.Port
	log.Printf("server rгn on port %s", port)
	if err := r.Run(port); err != nil {
		log.Fatal("Server startup error")
	}
}
