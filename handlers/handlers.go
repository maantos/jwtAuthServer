package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/maantos/jwtAuth/initializers"
	"github.com/maantos/jwtAuth/models"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(rw http.ResponseWriter, r *http.Request) {
	fmt.Print("Handling Signup")
	var user struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil {
		http.Error(rw, "error decoding request", http.StatusBadRequest)
		return
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)

	if err != nil {
		http.Error(rw, "failed to hash password", http.StatusInternalServerError)
	}

	object := models.User{Email: user.Email, Password: string(hashedPass)}

	result := initializers.DB.Create(&object)

	if result.Error != nil {
		http.Error(rw, "failed to create user", http.StatusInternalServerError)
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)

	data, _ := json.Marshal(object.ID)
	rw.Write(data)
}
