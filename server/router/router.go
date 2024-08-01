package router

import (
	"api/controller"
	"api/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()
	// manager := controller.NewManager()
	router.Use(middleware.CorsMiddleware)

	router.Handle("/", http.FileServer(http.Dir("../frontend")))
	router.HandleFunc("/signup", controller.Signup).Methods("POST")
	router.HandleFunc("/login", controller.Login).Methods("POST")
	router.Handle("/validate", middleware.RequireAuth(http.HandlerFunc(controller.Validate)))
	router.Handle("/ws", middleware.RequireAuth(http.HandlerFunc(controller.ServeWS)))

	// middleware for /validate route
	// authMiddleware := router.PathPrefix("/validate").Subrouter()
	// authMiddleware.Use(middleware.RequireAuth)
	// authMiddleware.HandleFunc("", controller.Validate).Methods("GET")

	return router
}
