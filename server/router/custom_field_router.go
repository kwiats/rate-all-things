package router

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kwiats/rate-all-things/internal/category/model"
	"github.com/kwiats/rate-all-things/internal/category/repository"
	"github.com/kwiats/rate-all-things/internal/category/service"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
	"sync"
)

func HandleCustomFieldRouter(db *gorm.DB, router *mux.Router, wg *sync.WaitGroup) {
	defer wg.Done()

	customFieldRouter := router.PathPrefix("/custom-field").Subrouter()
	customFieldRepo := repository.NewCustomFieldRepository(db)
	customFieldService := service.NewCustomFieldService(customFieldRepo)

	customFieldRouter.HandleFunc("/", handleCreateCustomField(customFieldService)).Methods("POST")
	customFieldRouter.HandleFunc("/{id}/", handleGetCustomField(customFieldService)).Methods("GET")
	customFieldRouter.HandleFunc("/", handleGetCustomFields(customFieldService)).Methods("GET")
	customFieldRouter.HandleFunc("/{id}/", handleDeleteCustomField(customFieldService)).Methods("DELETE")
	customFieldRouter.HandleFunc("/{id}/", handleUpdateCustomField(customFieldService)).Methods("PATCH")
}

func handleCreateCustomField(customFieldService *service.CustomFieldService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var customFieldDTO model.CustomFieldDTO

		if err := json.NewDecoder(r.Body).Decode(&customFieldDTO); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Printf("Error parsing request body to DTO: %v\n", err)
			return
		}

		createdCustomField, err := customFieldService.CreateCustomField(customFieldDTO)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("Custom field not created: %v\n", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(createdCustomField)
		if err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			log.Printf("Error encoding JSON response: %v\n", err)
			return
		}
	}
}
func handleDeleteCustomField(customFieldService *service.CustomFieldService) http.HandlerFunc {
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

		idCustomField, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Uncorrected custom field id", http.StatusBadRequest)
			log.Printf("Error parsing custom field ID: %v\n", err)
			return
		}

		isDeleted, err := customFieldService.DeleteCustomField(uint(idCustomField), forceDelete)
		if err != nil {
			http.Error(w, fmt.Sprintf("Custom field with ID %d not found: %v", idCustomField, err), http.StatusNotFound)
			log.Printf("Custom field not found: %v\n", err)
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

func handleGetCustomField(customFieldService *service.CustomFieldService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]
		idCustomField, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Uncorrected custom field id", http.StatusBadRequest)
			log.Printf("Error parsing custom field ID: %v\n", err)
			return
		}

		customField, err := customFieldService.GetCustomField(uint(idCustomField))
		if err != nil {
			http.Error(w, fmt.Sprintf("Custom field with ID %d not found: %v", idCustomField, err), http.StatusNotFound)
			log.Printf("Custom field not found: %v\n", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(customField)
		if err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			log.Printf("Error encoding JSON response: %v\n", err)
			return
		}

	}
}
func handleUpdateCustomField(customFieldService *service.CustomFieldService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]
		idCustomField, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Uncorrected custom field id", http.StatusBadRequest)
			log.Printf("Error parsing custom field ID: %v\n", err)
			return
		}

		var updatedCustomField model.CustomFieldDTO

		if err := json.NewDecoder(r.Body).Decode(&updatedCustomField); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Printf("Error parsing request body to DTO: %v\n", err)
			return
		}

		customField, err := customFieldService.UpdateCustomField(uint(idCustomField), updatedCustomField)
		if err != nil {
			http.Error(w, fmt.Sprintf("Custom field with ID %d not found: %v", idCustomField, err), http.StatusNotFound)
			log.Printf("Custom field not found: %v\n", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(customField)
		if err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			log.Printf("Error encoding JSON response: %v\n", err)
			return
		}

	}
}

func handleGetCustomFields(customFieldService *service.CustomFieldService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		customFields, err := customFieldService.GetCustomFields()
		if err != nil {
			http.Error(w, fmt.Sprintf("Custom fields not found %v", err), http.StatusNotFound)
			log.Printf("Custom field not found: %v\n", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(customFields)
		if err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			log.Printf("Error encoding JSON response: %v\n", err)
			return
		}

	}
}
