package main

import (
	"log"
	// "os/user"

	"micromiro/database"
	"micromiro/handlers"
	"micromiro/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	db, err := database.ConnectDB()
	if err != nil {
		log.Fatalf("failed to start the server: %v", err)
	}
	defer db.Close()

	router := gin.Default()

	v1 := router.Group("/api/v1")
	{
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})	
		v1.POST("/register", handlers.Register)
		v1.POST("/login", handlers.Login)

		protected := v1.Group("/protected")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.GET("/profile", func(c *gin.Context){
				userID, _ := c.Get("user_id")
				email, _ := c.Get("email")
				c.JSON(200, gin.H{
					"user_id": userID,
					"email": email,
				})
			})
		}
	}

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
