package user

import (
	"sync"
	"tit/internal/app/middleware"

	"tit/pkg/media"

	"github.com/gorilla/mux"

	"gorm.io/gorm"
)

func HandleUserRouter(db *gorm.DB, r *mux.Router, wg *sync.WaitGroup) {
	defer wg.Done()

	router := r.PathPrefix("/api/user").Subrouter()
	userRepo := NewUserRepository(db)
	localStorage := media.NewBlobStorage("uploads")
	storage := media.NewMediaService(localStorage)
	userService := NewUserService(userRepo, storage)

	router.HandleFunc("/", HandleCreate(userService)).Methods("POST")
	router.HandleFunc("/", middleware.AuthMiddleware(HandleGet(userService))).Methods("GET")
	router.HandleFunc("/{id}", middleware.AuthMiddleware(HandleGetById(userService))).Methods("GET")
	router.HandleFunc("/{username}", middleware.AuthMiddleware(HandleGetByUsername(userService))).Methods("GET")
	router.HandleFunc("/{id}", middleware.AuthMiddleware(HandleDelete(userService))).Methods("GET")
}
