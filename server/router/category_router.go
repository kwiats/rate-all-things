package router

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

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
	router.HandleFunc("/category/{id}", handleGetCategory(categoryService)).Methods("GET")
	router.HandleFunc("/category", handleGetCategories(categoryService)).Methods("GET")

}

func handleCreateCategory(categoryService *service.CategoryService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var categoryDTO model.CategoryDTO

		if err := json.NewDecoder(r.Body).Decode(&categoryDTO); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Printf("Error parsing request body to DTO: %v\n", err)
			return
		}

		createdCategory, err := categoryService.CreateCategory(categoryDTO)
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

func handleGetCategory(categoryService *service.CategoryService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]
		idCategory, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Uncorrected category id", http.StatusBadRequest)
			log.Printf("Error parsing category ID: %v\n", err)
			return
		}

		category, err := categoryService.GetCategory(uint(idCategory))
		if err != nil {
			http.Error(w, fmt.Sprintf("Category with ID %d not found: %v", idCategory, err), http.StatusNotFound)
			log.Printf("Category not found: %v\n", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(category)
		if err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			log.Printf("Error encoding JSON response: %v\n", err)
			return
		}

	}
}

func handleGetCategories(categoryService *service.CategoryService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		category, err := categoryService.GetCategories()
		if err != nil {
			http.Error(w, fmt.Sprintf("Categories not found", err), http.StatusNotFound)
			log.Printf("Categories not found: %v\n", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(category)
		if err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			log.Printf("Error encoding JSON response: %v\n", err)
			return
		}

	}
}
