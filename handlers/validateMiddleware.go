package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/maantos/jwtAuth/initializers"
	"github.com/maantos/jwtAuth/models"
)

func Test(rw http.ResponseWriter, r *http.Request) {

	user := r.Context().Value(User{}).(models.User)

	e := json.NewEncoder(rw)

	err := e.Encode(user)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

type User struct{}

func MyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("accessToken")
		if err != nil {
			//eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NzcxODQ5MTAsInN1YiI6MTF9.h4hxy67gcgYbvixkT5Rkw4oFrZ-q1DgozQYOGrcZjrQ
			http.Error(w, "token is missing", http.StatusUnauthorized)
			return
		}

		// Parse takes the token string and a function for looking up the key. The latter is especially
		// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
		// head of the token to identify which key to use, but the parsed token (head and claims) is provided
		// to the callback, providing flexibility.
		token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
			return []byte(os.Getenv("SECRET")), nil
		})

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			//check exp date
			if float64(time.Now().Unix()) > claims["exp"].(float64) {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			var user models.User
			initializers.DB.First(&user, claims["sub"])
			if user.ID == 0 {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), User{}, user)
			next.ServeHTTP(w, r.WithContext(ctx))

		} else {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
	})
}
