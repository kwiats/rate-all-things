package user

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"tit/pkg/common"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

func HandleCreate(userService *UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Printf("Error parsing multipart form: %v\n", err)
			return
		}

		file, handler, err := r.FormFile("profile_picture")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Printf("Error retrieving the file from form data: %v\n", err)
			return
		}
		defer file.Close()
		var userCreateDTO UserCreateDTO
		decoder := schema.NewDecoder()
		if err := decoder.Decode(&userCreateDTO, r.PostForm); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Printf("Error decoding form data to DTO: %v\n", err)
			return
		}

		createdUser, err := userService.CreateUser(userCreateDTO, handler)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("User not created: %v\n", err)
			return
		}

		common.WriteJSON(w, http.StatusCreated, createdUser)

	}
}

func HandleGet(userService *UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := userService.GetUsers()
		if err != nil {
			return
		}

		common.WriteJSON(w, http.StatusOK, users)

	}
}

func HandleSearch(userService *UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q")
		if query == "" {
			http.Error(w, "Query parameter 'q' is required", http.StatusBadRequest)
			return
		}

		users, err := userService.SearchUsers(query)
		if err != nil {
			return
		}

		common.WriteJSON(w, http.StatusOK, users)

	}
}

func HandleGetById(userService *UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]
		userId, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Uncorrected user id", http.StatusBadRequest)
			log.Printf("Error parsing user ID: %v\n", err)
			return
		}
		users, err := userService.GetUserById(uint(userId))
		if err != nil {
			return
		}


		common.WriteJSON(w, http.StatusOK, users)
	}
}

func HandleDelete(userService *UserService) http.HandlerFunc {
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

		userId, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Uncorrected user id", http.StatusBadRequest)
			log.Printf("Error parsing user ID: %v\n", err)
			return
		}

		isDeleted, err := userService.DeleteUser(uint(userId), forceDelete)
		if err != nil {
			http.Error(w, fmt.Sprintf("Category with ID %d not found: %v", userId, err), http.StatusNotFound)
			log.Printf("Category not found: %v\n", err)
			return
		}

		common.WriteJSON(w, http.StatusOK, isDeleted)

	}
}

func HandleGetByUsername(userService *UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		username := vars["username"]

		user, err := userService.GetUserByUsername(username)
		if err != nil {
			http.Error(w, fmt.Sprintf("Category with Username %v not found: %v", username, err), http.StatusNotFound)
			log.Printf("Category not found: %v\n", err)
			return
		}

		common.WriteJSON(w, http.StatusOK, user)

	}
}
