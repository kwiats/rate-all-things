package user

import (
	"log"
	"sync"
	"tit/internal/app/middleware"
	"tit/internal/config"

	"tit/pkg/media"

	"github.com/gorilla/mux"

	"gorm.io/gorm"
)

func HandleUserRouter(db *gorm.DB, r *mux.Router, wg *sync.WaitGroup) {
	defer wg.Done()

	router := r.PathPrefix("/api/user").Subrouter()
	userRepo := NewUserRepository(db)
	settings, err := config.NewConfiguration()
	if err != nil {
		log.Fatalf("failed to load config file: %v", err)
	}
	localStorage, err := media.NewMinioBlobStorage(settings.Envs.MinioHost, settings.Envs.MinioAccessKey, settings.Envs.MinioSecretAccessKey, settings.Envs.MinioIsSecure, "pictures")
	if err != nil {
		log.Fatal(err)
	}
	storage := media.NewMediaService(localStorage)
	userService := NewUserService(userRepo, storage)

	router.HandleFunc("/", HandleCreate(userService)).Methods("POST")
	router.HandleFunc("/", middleware.AuthMiddleware(HandleGet(userService))).Methods("GET")
		router.HandleFunc("/search/", middleware.AuthMiddleware(HandleSearch(userService))).Methods("GET")

	router.HandleFunc("/{id}", middleware.AuthMiddleware(HandleGetById(userService))).Methods("GET")
	router.HandleFunc("/{username}", middleware.AuthMiddleware(HandleGetByUsername(userService))).Methods("GET")
	router.HandleFunc("/{id}", middleware.AuthMiddleware(HandleDelete(userService))).Methods("GET")
}
