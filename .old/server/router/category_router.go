package router

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/kwiats/rate-all-things/pkg/schema"
	"github.com/kwiats/rate-all-things/server/middleware"

	"github.com/gorilla/mux"
	"github.com/kwiats/rate-all-things/internal/category/repository"
	"github.com/kwiats/rate-all-things/internal/category/service"
	"gorm.io/gorm"
)

func HandleCategoryRouter(db *gorm.DB, r *mux.Router, wg *sync.WaitGroup) {
	defer wg.Done()

	router := r.PathPrefix("/category").Subrouter()
	categoryRepo := repository.NewCategoryRepository(db)
	categoryService := service.NewCategoryService(categoryRepo)

	router.HandleFunc("/", middleware.AuthMiddleware(handleCreateCategory(categoryService))).Methods("POST")
	router.HandleFunc("/{id}/", middleware.AuthMiddleware(handleGetCategory(categoryService))).Methods("GET")
	router.HandleFunc("/", middleware.AuthMiddleware(handleGetCategories(categoryService))).Methods("GET")
	router.HandleFunc("/{id}/", middleware.AuthMiddleware(handleDeleteCategory(categoryService))).Methods("DELETE")
	router.HandleFunc("/{id}/", middleware.AuthMiddleware(handleUpdateCategory(categoryService))).Methods("PATCH")
}

func handleCreateCategory(categoryService *service.CategoryService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var categoryDTO schema.CreateCategoryDTO

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
func handleDeleteCategory(categoryService *service.CategoryService) http.HandlerFunc {
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

		idCategory, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Uncorrected category id", http.StatusBadRequest)
			log.Printf("Error parsing category ID: %v\n", err)
			return
		}

		isDeleted, err := categoryService.DeleteCategory(uint(idCategory), forceDelete)
		if err != nil {
			http.Error(w, fmt.Sprintf("Category with ID %d not found: %v", idCategory, err), http.StatusNotFound)
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
func handleUpdateCategory(categoryService *service.CategoryService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]
		idCategory, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Uncorrected category id", http.StatusBadRequest)
			log.Printf("Error parsing category ID: %v\n", err)
			return
		}

		var updatedCategory schema.UpdateCategoryDTO

		if err := json.NewDecoder(r.Body).Decode(&updatedCategory); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Printf("Error parsing request body to DTO: %v\n", err)
			return
		}

		category, err := categoryService.UpdateCategory(uint(idCategory), updatedCategory)
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
			http.Error(w, fmt.Sprintf("Categories not found %v", err), http.StatusNotFound)
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
