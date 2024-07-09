package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/lib/pq"

	"github.com/joho/godotenv"
)

type DBconn struct {
	DB *sql.DB
}

func ConnectDB() (*sql.DB, error) {

	err := godotenv.Load(os.Getenv("KO_DATA_PATH") + "/.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	var host string = os.Getenv("DB_HOST")
	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		log.Fatalf("Error converting DB_PORT to integer: %v", err)
	}
	var user string = os.Getenv("DB_USER")
	var password string = os.Getenv("DB_PASSWORD")
	var dbname string = os.Getenv("DB_NAME")

	connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname = %s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", connString)
	if err != nil {
		log.Printf("failed to connect to database: %v", err)
		return &sql.DB{}, err
	}
	return db, nil
}
