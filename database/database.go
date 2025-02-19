package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	_ "github.com/lib/pq"
)

func ConnectDB() (*sql.DB, error) {

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	
	if host == "" || port == "" || user == "" || password == "" || dbname == "" {
		log.Printf("Database configuration: host=%q, port=%q, user=%q, dbname=%q",
			host, port, user, dbname)
		return nil, fmt.Errorf("отсутствуют необходимые переменные окружения для базы данных")
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
	host, port, user, password, dbname)
	
	log.Printf("Пытаемся подключиться к базе данных с параметрами: %s", connStr)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Printf("ошибка подключения к базе данных: %v", err)
		return nil, err
	}

	if err = db.Ping(); err != nil {
		log.Printf("ошибка проверки соединения с базой данных: %v", err)
		return nil, err
	}
	
	return db, nil
}


