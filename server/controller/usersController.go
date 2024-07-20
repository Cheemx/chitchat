package controller

import (
	"api/db"
	"api/middleware"
	"api/model"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func Signup(w http.ResponseWriter, r *http.Request) {
	// setting headers
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	// Get the username and password from req body
	var user model.User
	_ = json.NewDecoder(r.Body).Decode(&user)

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		log.Fatal("error while hashing the password")
	}
	user.Password = string(hashedPassword)

	// Create the user
	_, inserr := db.UserCollection.InsertOne(context.TODO(), user)
	if inserr != nil {
		log.Fatal("Failed to create user")
	}

	// Respond
	json.NewEncoder(w).Encode(user)
}

func Login(w http.ResponseWriter, r *http.Request) {
	// setting headers
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "POST")

	// Get the username and pass off request body
	var body model.User
	errordecode := json.NewDecoder(r.Body).Decode(&body)

	if errordecode != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Error while decoding body")
		return
	}
	// Look up for requested user

	filter := bson.M{"username": body.Username}
	var res model.User
	err := db.UserCollection.FindOne(context.TODO(), filter).Decode(&res)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode("No user found")
			return
		}
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Error while retrieving the user")
		return
	}

	fmt.Printf("found document %v\n", res)

	// Compare sent in pass with saved user pass hash
	e := bcrypt.CompareHashAndPassword([]byte(res.Password), []byte(body.Password))

	if e != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("Invalid Username or password while comparing")
		return
	}

	// Generate a jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": res.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	// get the complete encoded token as a string using the secret key
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Failed to create token")
		log.Println(err)
		return
	}

	// Send it back
	cookie := http.Cookie{
		Name:     "Authorization",
		Value:    tokenString,
		Path:     "",
		MaxAge:   3600 * 24 * 30,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &cookie)

	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

func Validate(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserContextKey).(model.User)
	if !ok {
		http.Error(w, "No User found in Context", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(user)
	json.NewEncoder(w).Encode("I'm in Validate function")
}
