package router

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
	"github.com/kwiats/rate-all-things/internal/user/repository"
	"github.com/kwiats/rate-all-things/internal/user/service"
	"github.com/kwiats/rate-all-things/pkg/schema"
	"github.com/kwiats/rate-all-things/server/middleware"
	"gorm.io/gorm"
)

func HandleUserRouter(db *gorm.DB, r *mux.Router, wg *sync.WaitGroup) {
	defer wg.Done()

	router := r.PathPrefix("/user").Subrouter()
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)

	router.HandleFunc("/", handleCreateUser(userService)).Methods("POST")
	router.HandleFunc("/", middleware.AuthMiddleware(handleGetUsers(userService))).Methods("GET")
	router.HandleFunc("/{id}", middleware.AuthMiddleware(handleGetUserById(userService))).Methods("GET")
	router.HandleFunc("/{username}", middleware.AuthMiddleware(handleGetUserByUsername(userService))).Methods("GET")
	router.HandleFunc("/{id}", middleware.AuthMiddleware(handleDeleteUser(userService))).Methods("GET")
}

func handleCreateUser(userService *service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var userCreateDTO schema.UserCreateDTO

		if err := json.NewDecoder(r.Body).Decode(&userCreateDTO); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Printf("Error parsing request body to DTO: %v\n", err)
			return
		}

		createdUser, err := userService.CreateUser(userCreateDTO)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("User not created: %v\n", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(createdUser)
		if err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			log.Printf("Error encoding JSON response: %v\n", err)
			return
		}
	}
}

func handleGetUsers(userService *service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := userService.GetUsers()
		if err != nil {
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(users)
		if err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			log.Printf("Error encoding JSON response: %v\n", err)
			return
		}

	}
}

func handleGetUserById(userService *service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]
		userId, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Uncorrected user id", http.StatusBadRequest)
			log.Printf("Error parsing user ID: %v\n", err)
			return
		}
		users, err := userService.GetUserById(uint(userId))
		if err != nil {
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(users)
		if err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			log.Printf("Error encoding JSON response: %v\n", err)
			return
		}

	}
}

func handleDeleteUser(userService *service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]

		var forceDelete bool
		forceDeleteStr := r.URL.Query().Get("force")
		if forceDeleteStr != "" {
			var err error
			forceDelete, err = strconv.ParseBool(forceDeleteStr)
			if err != nil {
				http.Error(w, "Cannot parse query params", http.StatusBadRequest)
				log.Printf("Cannot parse query param: %v\n", err)
				return
			}
		}

		userId, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Uncorrected user id", http.StatusBadRequest)
			log.Printf("Error parsing user ID: %v\n", err)
			return
		}

		isDeleted, err := userService.DeleteUser(uint(userId), forceDelete)
		if err != nil {
			http.Error(w, fmt.Sprintf("Category with ID %d not found: %v", userId, err), http.StatusNotFound)
			log.Printf("Category not found: %v\n", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(isDeleted)
		if err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			log.Printf("Error encoding JSON response: %v\n", err)
			return
		}

	}
}

func handleGetUserByUsername(userService *service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		username := vars["username"]

		user, err := userService.GetUserByUsername(username)
		if err != nil {
			http.Error(w, fmt.Sprintf("Category with Username %v not found: %v", username, err), http.StatusNotFound)
			log.Printf("Category not found: %v\n", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(user)
		if err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			log.Printf("Error encoding JSON response: %v\n", err)
			return
		}

	}
}
