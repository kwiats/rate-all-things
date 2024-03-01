package router

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kwiats/rate-all-things/internal/category/model"
	"github.com/kwiats/rate-all-things/internal/category/repository"
	"github.com/kwiats/rate-all-things/internal/category/service"
	"gorm.io/gorm"
)

func HandleCategoryRouter(db *gorm.DB, router *mux.Router) {
	categoryRepo := repository.NewCategoryRepository(db)
	categoryService := service.NewCategoryService(categoryRepo)

	router.HandleFunc("/category", handleCreateCategory(categoryService)).Methods("POST")

}

func handleCreateCategory(categoryService *service.CategoryService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var categoryDTO model.CategoryDTO

		if err := json.NewDecoder(r.Body).Decode(&categoryDTO); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		createdCategory, err := categoryService.CreateCategory(categoryDTO)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(createdCategory)
		if err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}
