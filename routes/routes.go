package routes

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"marketplace-app/auth"
	"marketplace-app/bank_account"
	"marketplace-app/middlewares"
	"marketplace-app/product"
	"marketplace-app/repositories"

	_ "github.com/lib/pq"
)

func SetupRouter(db *sql.DB) *gin.Engine {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to your marketplace app!",
		})
	})

	// Version 1 router group
	v1 := r.Group("/v1")

	// user group
	userRepo := repositories.NewUserRepository(db)
	authHandler := auth.NewHandler(userRepo)

	user := v1.Group("/user")
	user.POST("/register", authHandler.Register)
	user.POST("/login", authHandler.Login)

	// product group
	productRepo := repositories.NewProductRepository(db)
	productHandler := product.NewProductHandler(productRepo)

	productGroup := v1.Group("/product")
	productGroup.Use(middlewares.AuthMiddleware())
	productGroup.PATCH("/:id/stock", productHandler.UpdateStock)
	productGroup.PATCH("/:id", productHandler.UpdateProduct)
	productGroup.DELETE("/:id", productHandler.DeleteProduct)
	productGroup.POST("/", productHandler.CreateProduct)
	productGroup.GET("/", productHandler.GetProducts)

	image := v1.Group("/image")
	image.Use(middlewares.AuthMiddleware())
	image.POST("/", productHandler.UploadImage)

	// bank account
	bankAccountRepo := repositories.NewBankAccountRepository(db)
	bankAccountHandler := bank_account.NewBankAccountHandler(bankAccountRepo)

	bankAccount := v1.Group("/bank/account")
	bankAccount.Use(middlewares.AuthMiddleware())
	bankAccount.PATCH("/:id", bankAccountHandler.UpdateBankAccount)
	bankAccount.DELETE("/:id", bankAccountHandler.DeleteBankAccount)
	bankAccount.POST("/", bankAccountHandler.CreateBankAccount)
	bankAccount.GET("/", bankAccountHandler.GetBankAccount)

	return r
}
