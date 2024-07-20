package controllers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/rishabhs-s/db"
	"github.com/rishabhs-s/models"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {

	var newUser models.User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("Error in hasing user password")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	slog.Info(" password hashed")

	newUser.Password = string(hashedPassword)

	if err := db.Db.Create(&newUser).Error; err != nil {
		slog.Error("Error in creating user")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	slog.Info(fmt.Sprintf("User %s successfully created", newUser.Username))
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}
