package routes

import (
	"database/sql"

	"github.com/gorilla/mux"
	"github.com/juang77/GoMSArch/authentication-service/app/http/controllers"
	"github.com/juang77/GoMSArch/authentication-service/config"
	"github.com/juang77/GoMSArch/shared/util/middleware"
	"github.com/urfave/negroni"
)

// InitRoutes instantiates a new gorilla/mux router
func InitRoutes(db *sql.DB, cnf config.Config) *mux.Router {
	router := mux.NewRouter()
	router = setAuthenticationRoutes(db, cnf, router)
	return router
}

// setAuthenticationRoutes specifies all routes for the authentication service
func setAuthenticationRoutes(db *sql.DB, cnf config.Config, router *mux.Router) *mux.Router {

	// Subrouter /token-auth
	tokenAUTH := router.PathPrefix("/token-auth").Subrouter()

	// User Login POST /token-auth
	tokenAUTH.Methods("POST").Handler(negroni.New(
		negroni.HandlerFunc(middleware.AccessControlHandler),
		controllers.LoginHandler(db, cnf),
	))

	// OPTIONS /token-auth
	tokenAUTH.Methods("OPTIONS").Handler(negroni.New(
		negroni.HandlerFunc(middleware.AcceptOPTIONS),
	))

	return router
}
