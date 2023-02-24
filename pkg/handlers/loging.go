package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/maantos/jwtAuth/pkg/db"
	"github.com/maantos/jwtAuth/pkg/models"

	"golang.org/x/crypto/bcrypt"
)

func Login(rw http.ResponseWriter, r *http.Request) {

	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	if err != nil {
		http.Error(rw, "error decoding request", http.StatusBadRequest)
		return
	}
	var user models.User
	db.DB.First(&user, "email = ?", body.Email)

	if user.ID == 0 {
		http.Error(rw, "invalid email or password", http.StatusBadRequest)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		http.Error(rw, "invalid email or password", http.StatusBadRequest)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Minute * 15).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	fmt.Print(tokenString)
	if err != nil {
		fmt.Println(err)
		http.Error(rw, "failed to create token", http.StatusBadRequest)
		return
	}

	cookie := &http.Cookie{
		Name:     "accessToken",
		Value:    tokenString,
		MaxAge:   3600,
		Path:     "/",
		HttpOnly: true,
		// Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(rw, cookie)
	fmt.Fprintf(rw, "aaa")
}
