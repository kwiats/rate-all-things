package router

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
	"github.com/kwiats/rate-all-things/internal/review/repository"
	"github.com/kwiats/rate-all-things/internal/review/service"
	"github.com/kwiats/rate-all-things/pkg/media"
	"github.com/kwiats/rate-all-things/pkg/schema"
	"github.com/kwiats/rate-all-things/server/middleware"
	"gorm.io/gorm"
)

func HandleReviewItemRouter(db *gorm.DB, r *mux.Router, wg *sync.WaitGroup) {
	defer wg.Done()

	router := r.PathPrefix("/review-item").Subrouter()
	reviewRepo := repository.NewReviewRepository(db)

	localBlobStorage := media.NewBlobStorage("media/files/review-item/")
	mediaService := media.NewMediaService(localBlobStorage)

	reviewService := service.NewReviewService(reviewRepo, mediaService)

	router.HandleFunc("/", handleGetReviewItems(reviewService)).Methods("GET")
	router.HandleFunc("/", middleware.AuthMiddleware(handleCreateReviewItems(reviewService))).Methods("POST")
	router.HandleFunc("/{id}/", middleware.AuthMiddleware(handleGetReviewItem(reviewService))).Methods("GET")
}

func handleGetReviewItems(service *service.ReviewService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reviews, err := service.ReviewItems()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("User not created: %v\n", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(reviews)
		if err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			log.Printf("Error encoding JSON response: %v\n", err)
			return
		}
	}
}

func handleGetReviewItem(service *service.ReviewService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]
		reviewId, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Uncorrected review item id", http.StatusBadRequest)
			log.Printf("Error parsing review item ID: %v\n", err)
			return
		}

		reviews, err := service.ReviewItem(uint(reviewId))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("User not created: %v\n", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(reviews)
		if err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			log.Printf("Error encoding JSON response: %v\n", err)
			return
		}
	}
}
func handleCreateReviewItems(service *service.ReviewService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			http.Error(w, "Failed to parse multipart form: "+err.Error(), http.StatusBadRequest)
			log.Printf("Error parsing multipart form: %v\n", err)
			return
		}

		var review schema.CreateReviewItem

		if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Printf("Error parsing request body to DTO: %v\n", err)
			return
		}
		createdReviewItem, err := service.CreateReviewItem(1, review)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("Review item not created: %v\n", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(createdReviewItem)
		if err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			log.Printf("Error encoding JSON response: %v\n", err)
			return
		}
	}
}
