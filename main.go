package main

import (
	"io"
	"os"

	"micromiro/database"
	"micromiro/handlers"
	"micromiro/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func setupLogger() *log.Logger{
	logger := log.New()

	file, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Ошибка открытия файла лога: %v", err)
	}
	logger.SetOutput(file)
	mw := io.MultiWriter(os.Stdout, file)
	logger.SetOutput(mw)

	return logger
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	logger := setupLogger()
	logger.Info("Запуск приложения...")

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

	port := ":8080"
	if os.Getenv("USE_TLS") == "true" {
		router.RunTLS(port, "cert.pem", "key.pem")
		logger.WithField("port", port).Info("Starting server")
	} else {
		router.Run(":8080")
		logger.WithField("port", port).Info("Starting server")
	}
}
