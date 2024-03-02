package router

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/kwiats/rate-all-things/internal/category/model"
	"github.com/kwiats/rate-all-things/internal/category/repository"
	"github.com/kwiats/rate-all-things/internal/category/service"
	"gorm.io/gorm"
	"log"
	"net/http"
	"sync"
)

func HandleCustomFieldRouter(db *gorm.DB, router *mux.Router, wg *sync.WaitGroup) {
	defer wg.Done()

	customFieldRouter := router.PathPrefix("/custom-field").Subrouter()
	customFieldRepo := repository.NewCustomFieldRepository(db)
	customFieldService := service.NewCustomFieldService(customFieldRepo)

	customFieldRouter.HandleFunc("/", handleCreateCustomField(customFieldService)).Methods("POST")

}

func handleCreateCustomField(customFieldService *service.CustomFieldService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var customFieldDTO model.CustomFieldDTO

		if err := json.NewDecoder(r.Body).Decode(&customFieldDTO); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Printf("Error parsing request body to DTO: %v\n", err)
			return
		}

		createdCategory, err := customFieldService.CreateCustomField(customFieldDTO)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("Category not created: %v\n", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(createdCategory)
		if err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			log.Printf("Error encoding JSON response: %v\n", err)
			return
		}
	}
}
