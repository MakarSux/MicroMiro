package handlers

import (
	"database/sql"
	"micromiro/database"
	"micromiro/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
    "golang.org/x/crypto/bcrypt"
)


func Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка хеширования пароля"})
		return
	}

	db, _ := database.ConnectDB()
	defer db.Close()

	query := `INSERT INTO users (username, email, password, role_id, created_at, updated_at) 
              VALUES ($1, $2, $3, 1, $4, $5) RETURNING id`

	var userID int
	err = db.QueryRow(query, req.Username, req.Email, hashedPassword, time.Now(), time.Now()).Scan(&userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания пользователя"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Пользователь успешно зарегистрирован", "user_id": userID})
}

func Login(c *gin.Context) {
    var req models.LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    db, _ := database.ConnectDB()
    defer db.Close()

    var user models.User
    query := `SELECT id, username, email, password, role_id FROM users WHERE email = $1`
    err := db.QueryRow(query, req.Email).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.RoleID)
    if err == sql.ErrNoRows {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    } else if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
        return
    }

    // Проверяем пароль
    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    // Генерируем JWT токен
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": user.ID,
        "email":   user.Email,
        "role_id": user.RoleID,
        "exp":     time.Now().Add(time.Hour * 24).Unix(), // Токен истекает через 24 часа
    })

    secret := os.Getenv("JWT_SECRET")
    if secret == "" {
        secret = "your-secret-key" // Укажи в .env
    }
    tokenString, err := token.SignedString([]byte(secret))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"token": tokenString})
}