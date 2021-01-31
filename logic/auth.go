package logic

import (
	"encoding/json"
	"github.com/golobby/container"
	"go-jwt-server/dal"
	"go-jwt-server/models"
	"go-jwt-server/types"
	"log"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {

}

func Logout(w http.ResponseWriter, r *http.Request) {

}

func Register(w http.ResponseWriter, r *http.Request) {
	var db *types.DatabaseT
	err := container.Make(&db)

	if err != nil {
		panic("Failed to resolve db context")
	}

	var user models.User
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Failed decoding JSON body", 400)
		return
	}

	_, err = dal.AddUser(db, &user)
	log.Printf("Successfully added a user with is %d", user.Id)

	err = json.NewEncoder(w).Encode(user)

	if err != nil {
		http.Error(w, "Failed encoding JSON result", 400)
		return
	}
}
