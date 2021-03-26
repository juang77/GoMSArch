package controllers

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"

	"errors"

	"github.com/juang77/GoMSArch/profile-service/app/models"
	"github.com/juang77/GoMSArch/profile-service/config"
	db "github.com/juang77/GoMSArch/profile-service/database"
	"github.com/juang77/GoMSArch/shared/util"

	"github.com/gorilla/mux"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/urfave/negroni"
)

// CreateUserHandler creates a new user in the database. Password is saved as a hash.
func CreateUserHandler(connection *sql.DB, cnf config.Config) negroni.HandlerFunc {
	return negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		user := &models.UserCreate{}
		err := util.RequestToJSON(r, user)
		if err != nil {
			util.SendBadRequest(w, errors.New("Bad json"))
			return
		}

		if err := user.Validate(); err == nil {
			if err := user.ValidatePassword(); err == nil {

				hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
				user.Hash = string(hash)

				createdID, err := db.CreateUser(connection, user)

				if err != nil {
					util.SendBadRequest(w, err)
					return
				}
				createdUser, err := db.GetUserByID(connection, createdID)
				if err != nil {
					util.SendBadRequest(w, err)
					return
				}

				// create JWT object with claims
				expiration := time.Now().Add(time.Hour * 24 * 31).Unix()
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
					"sub": createdUser.ID,
					"iat": time.Now().Unix(),
					"exp": expiration,
				})

				// Load secret key from config and generate a signed token
				secretKey := cnf.SecretKey
				tokenString, err := token.SignedString([]byte(secretKey))
				if err != nil {
					util.SendError(w, err)
					return
				}

				type Token struct {
					Token     string               `json:"token"`
					ExpiresOn string               `json:"expires_on"`
					User      *models.UserResponse `json:"user"`
				}

				util.SendOK(w, &Token{
					Token:     tokenString,
					ExpiresOn: strconv.Itoa(int(expiration)),
					User:      &createdUser,
				})

			} else {
				util.SendBadRequest(w, err)
			}
		} else {
			util.SendBadRequest(w, err)
		}
	})
}

// DeleteUserHandler removes a user from the database. User can only deletes it's own record.
func DeleteUserHandler(connection *sql.DB, cnf config.Config) negroni.HandlerFunc {
	return negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		var queryToken = r.URL.Query().Get("token")

		if len(queryToken) < 1 {
			queryToken = r.Header.Get("token")
		}

		if len(queryToken) < 1 {
			util.SendBadRequest(w, errors.New("token is mandatory"))
			return
		}

		user := &models.UserResponse{}
		err := util.RequestToJSON(r, user)
		if err != nil {
			util.SendBadRequest(w, errors.New("Bad json"))
			return
		}

		secretKey := cnf.SecretKey
		tok, err := jwt.Parse(queryToken, func(t *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})

		claims := tok.Claims.(jwt.MapClaims)
		var ID = claims["sub"].(float64)

		if int64(ID) != user.ID {
			util.SendBadRequest(w, errors.New("you can only delete your own user object"))
			return
		}

		db.DeleteUser(connection, user)
		if err != nil {
			util.SendBadRequest(w, err)
			return
		}
		util.SendOK(w, string(""))

	})
}

// UpdateUserHandler updates an user based on it's user ID. User is only allowed to update it's own record. Verification is being done based on the JWT in the request.
func UpdateUserHandler(connection *sql.DB, cnf config.Config) negroni.HandlerFunc {
	return negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		var queryToken = r.URL.Query().Get("token")

		if len(queryToken) < 1 {
			queryToken = r.Header.Get("token")
		}

		if len(queryToken) < 1 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(string("token is mandatory")))
			return
		}

		user := &models.UserResponse{}
		err := util.RequestToJSON(r, user)
		if err != nil {
			util.SendBadRequest(w, errors.New("bad json"))
			return
		}

		secretKey := cnf.SecretKey
		tok, err := jwt.Parse(queryToken, func(t *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})
		if err != nil {
			util.SendBadRequest(w, err)
			return
		}

		claims := tok.Claims.(jwt.MapClaims)
		var ID = claims["sub"].(float64) // gets the ID

		if int64(ID) != user.ID {
			util.SendBadRequest(w, errors.New("you can only change your own user object"))
			return
		}

		if err := user.Validate(); err == nil {

			db.UpdateUser(connection, user)

			util.SendOK(w, user)

		} else {
			util.SendBadRequest(w, err)
		}
	})
}

// UserByIndexHandler retrieves an user from the database based on its id. This handler expects the id being passed in the route variable in the current request.
func UserByIndexHandler(connection *sql.DB) negroni.HandlerFunc {
	return negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

		vars := mux.Vars(r)
		strID := vars["id"]

		id, err := strconv.Atoi(strID)
		if err != nil {
			util.SendBadRequest(w, err)
			return
		}

		user, err := db.GetUserByID(connection, int64(id))

		if err != nil {
			util.SendBadRequest(w, err)
			return
		}
		util.SendOK(w, user)
	})
}
