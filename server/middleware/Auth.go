package middleware

import (
	"api/db"
	"api/model"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type contextKey string

const UserContextKey contextKey = "user"

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Get the cookie off req body
		tokenString, err := r.Cookie("Authorization")

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode("Failed to parse the tokenstring")
			return
		}

		// Decode/ Validate it
		secret := os.Getenv("SECRET")
		t, err := jwt.Parse(tokenString.Value, func(t *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode("Invalid token")
			return
		}
		claims, ok := t.Claims.(jwt.MapClaims)
		if ok && t.Valid {
			// Check the expiration
			if float64(time.Now().Unix()) > claims["exp"].(float64) {
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode("Token Expired")
				return
			}

			//Find the user with token username
			var user model.User
			log.Printf("claims[\"sub\"]: %v", claims["sub"])
			id, e := primitive.ObjectIDFromHex(claims["sub"].(string))
			if e != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode("ID ka idea galat hai")
				return
			}
			filter := bson.M{"_id": id}
			err := db.UserCollection.FindOne(context.TODO(), filter).Decode(&user)

			if err != nil {
				if errors.Is(err, mongo.ErrNoDocuments) {
					w.WriteHeader(http.StatusUnauthorized)
					json.NewEncoder(w).Encode("Error while retrieving user with given token")
					return
				}
				http.Error(w, "Database error", http.StatusInternalServerError)
				return
			}

			// Attach to request
			ctx := context.WithValue(r.Context(), UserContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))

			json.NewEncoder(w).Encode("I'm in RequireAuth")
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode("Invalid token claims")
		}

	})
}
