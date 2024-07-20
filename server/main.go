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
	db.ConnectToDB()
}

func main() {
	r := router.Router()
	port := os.Getenv("PORT")
	fmt.Println("Server is getting started ...")

	// http.ListenAndServeTLS to make it https
	// change your origin in .env too

	log.Fatal(http.ListenAndServe(port, r))
	fmt.Printf("Listening at %v\n", port)

}
