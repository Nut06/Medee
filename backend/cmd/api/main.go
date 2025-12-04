package main

import (
	"fmt"
	"log"
	"os"

	"backend/internal/server"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("No .env was found %v", err)
	}
	
	fmt.Printf("DB port is at: %s\n", os.Getenv("DB_PORT"))
	app := server.NewServer()

	fmt.Println("Server is running on http://localhost:3000")
	log.Fatal(app.Listen(":3000"))
}