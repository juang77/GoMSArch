package routes

import (
	"database/sql"

	"github.com/juang77/GoMSArch/profile-service/app/http/controllers"
	"github.com/juang77/GoMSArch/profile-service/config"
	"github.com/juang77/shared/util/middleware"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

// InitRoutes initializes the REST and IPC routes for this service.
func InitRoutes(db *sql.DB, cnf config.Config) *mux.Router {
	router := mux.NewRouter()
	router = setRESTRoutes(db, cnf, router)
	return router
}

func setRESTRoutes(db *sql.DB, cnf config.Config, router *mux.Router) *mux.Router {

	// Subrouter /users
	users := router.PathPrefix("/users").Subrouter()

	// TODO :  https://github.com/gorilla/handlers/blob/master/cors.go#L140
	users.Methods("OPTIONS").Handler(negroni.New(
		negroni.HandlerFunc(middleware.AcceptOPTIONS),
	))

	// Update user /users
	users.Methods("PUT").Handler(negroni.New(
		negroni.HandlerFunc(middleware.AccessControlHandler),
		middleware.RequireTokenAuthenticationHandler(cnf.SecretKey),
		controllers.UpdateUserHandler(db, cnf),
	))

	// Delete User /users
	users.Methods("DELETE").Handler(negroni.New(
		negroni.HandlerFunc(middleware.AccessControlHandler),
		middleware.RequireTokenAuthenticationHandler(cnf.SecretKey),
		controllers.DeleteUserHandler(db, cnf),
	))

	// Create user /sers
	users.Methods("POST").Handler(negroni.New(
		negroni.HandlerFunc(middleware.AccessControlHandler),
		controllers.CreateUserHandler(db, cnf),
	))

	// Get one user /user/{id}
	oneUser := router.PathPrefix("/user/{id}").Subrouter()
	oneUser.Methods("GET").Handler(negroni.New(
		negroni.HandlerFunc(middleware.AccessControlHandler),
		middleware.RequireTokenAuthenticationHandler(cnf.SecretKey),
		controllers.UserByIndexHandler(db),
	))

	return router
}
