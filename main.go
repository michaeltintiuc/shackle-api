package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/joho/godotenv"
	"github.com/michaeltintiuc/shackle-api/pkg/app"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := getEnv("PORT", "8080")
	dbInfo := app.DbInfo{
		Host: os.Getenv("DB_HOST"),
		Port: os.Getenv("DB_PORT"),
		User: os.Getenv("DB_USER"),
		Pass: os.Getenv("DB_PASS"),
		Name: os.Getenv("DB_NAME"),
	}

	application, err := app.NewApp(port, dbInfo, os.Getenv("JWT_SECRET"))
	if err != nil {
		log.Fatalln(err)
	}

	go application.ListenAndServe()

	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	// Block until we receive our signal.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	application.Shutdown()
	os.Exit(0)
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
