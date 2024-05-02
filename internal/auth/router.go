package auth

import (
	"sync"
	"tit/pkg/media"

	"github.com/gorilla/mux"

	"gorm.io/gorm"
)

func HandleAuthRouter(db *gorm.DB, r *mux.Router, wg *sync.WaitGroup) {
	defer wg.Done()

	router := r.PathPrefix("/api/auth").Subrouter()
	minioStorage, _ := media.NewMinioBlobStorage("minio:9000", "bqvayHlLY5pHdimH", "zthpQFaAnbVvUwrobLCA0M3jVc6U9mdl", false, "pictures")
	mediaService := media.NewMediaService(minioStorage)

	authRepo := NewAuthRepository(db)
	authService := NewAuthService(authRepo, mediaService)

	router.HandleFunc("/login", HandleLogin(authService)).Methods("POST")
	router.HandleFunc("/register", HandleRegister(authService)).Methods("POST")
}
