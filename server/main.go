package main

import (
	"api/controller"
	"api/db"
	"api/router"
	"fmt"
	"log"
	"net/http"
	"os"
)

func init() {
	controller.LoadEnvVariables()
	// er := godotenv.Load(".env")

	// if er != nil {
	// 	log.Fatalf("Error loading .env file %v", er)
	// }
	db.ConnectToDB()
}

func main() {
	r := router.Router()
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("Error parsing the port number")
	}
	fmt.Println("Server is getting started ...")

	// http.ListenAndServeTLS to make it https
	// change your origin in .env too

	log.Fatal(http.ListenAndServe(port, r))
	fmt.Printf("Listening at %v\n", port)

}
