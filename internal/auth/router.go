package auth

import (
	"sync"

	"github.com/gorilla/mux"

	"gorm.io/gorm"
)

func HandleAuthRouter(db *gorm.DB, r *mux.Router, wg *sync.WaitGroup) {
	defer wg.Done()

	router := r.PathPrefix("/api/auth").Subrouter()
	authRepo := NewAuthRepository(db)
	authService := NewAuthService(authRepo)

	router.HandleFunc("/login", HandleLogin(authService)).Methods("POST")
	router.HandleFunc("/register", HandleRegister(authService)).Methods("POST")
}
