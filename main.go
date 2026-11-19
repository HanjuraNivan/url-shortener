package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/HanjuraNivan/url-shortener/api"
	"github.com/HanjuraNivan/url-shortener/services"
	urlshortenersqlc "github.com/HanjuraNivan/url-shortener/sqlc/urlshortenerqueries"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func main() {
	fmt.Println("Spinning up the URL Shortener service...")
	connection, err := connectToPostgres()
	if err != nil {
		fmt.Println("Error connecting to Postgres:", err)
		return
	}
	defer connection.Close(context.Background())
	err = connection.Ping(context.Background())
	if err != nil {
		fmt.Println("Cannot ping database. Please check env variables.")
		return
	}
	fmt.Println("Successfully connected to Postgres!")

	router := gin.Default()
	newServices := services.NewServices(urlshortenersqlc.New(connection))
	api.RegisterHandlers(router, newServices)
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		fmt.Println("Error starting HTTP server:", err)
		return
	}
}

func connectToPostgres() (*pgx.Conn, error) {
	dbPassword := os.Getenv("DB_PASSWORD")
	dbUser := os.Getenv("DB_USER")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	connectionString := fmt.Sprintf("postgresql://%s:%s@%s:%s/postgres", dbUser, dbPassword, dbHost, dbPort)
	conn, err := pgx.Connect(context.Background(), connectionString)
	if err != nil {
		return nil, errors.New("Cannot connect to database. Please check env variables.")
	}
	return conn, nil
}
