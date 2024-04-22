package auth

import (
	"encoding/json"
	"tit/pkg/common"

	"log"
	"net/http"
)

func HandleLogin(authService *AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var loginDTO LoginUserDTO

		if err := json.NewDecoder(r.Body).Decode(&loginDTO); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Printf("Error parsing request body to DTO: %v\n", err)
			return
		}

		accessToken, err := authService.LoginUser(loginDTO)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("User not created: %v\n", err)
			return
		}
		response := Token{Token: accessToken}
		common.WriteJSON(w, http.StatusAccepted, response)

	}
}
