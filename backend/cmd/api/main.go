package main

import (
	"fmt"
	"log"

	"backend/internal/server"
)

func main() {
	app := server.NewServer()

	fmt.Println("Server is running on http://localhost:3000")
	log.Fatal(app.Listen(":3000"))
}