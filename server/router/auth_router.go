package router

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/kwiats/rate-all-things/internal/auth/repository"
	"github.com/kwiats/rate-all-things/internal/auth/service"
	"github.com/kwiats/rate-all-things/pkg/schema"
	"gorm.io/gorm"
	"log"
	"net/http"
	"sync"
)

func HandleAuthRouter(db *gorm.DB, r *mux.Router, wg *sync.WaitGroup) {
	defer wg.Done()

	router := r.PathPrefix("/auth").Subrouter()
	authRepo := repository.NewAuthRepository(db)
	authService := service.NewAuthService(authRepo)

	router.HandleFunc("/login", handleLogin(authService)).Methods("POST")
}

func handleLogin(authService *service.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var loginDTO schema.LoginUserDTO

		if err := json.NewDecoder(r.Body).Decode(&loginDTO); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Printf("Error parsing request body to DTO: %v\n", err)
			return
		}

		accessToken, err := authService.LoginUser(loginDTO)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("User not created: %v\n", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(accessToken)
		if err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			log.Printf("Error encoding JSON response: %v\n", err)
			return
		}
	}
}
