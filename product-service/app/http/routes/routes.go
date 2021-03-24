package routes

import (
	"database/sql"

	"github.com/gorilla/mux"
	"github.com/juang77/GoMSArch/product-service/app/http/controllers"
	"github.com/juang77/GoMSArch/product-service/config"
	"github.com/juang77/GoMSArch/shared/util/middleware"
	"github.com/urfave/negroni"
)

// InitRoutes instantiates a new gorilla/mux router
func InitRoutes(db *sql.DB, cnf config.Config) *mux.Router {
	router := mux.NewRouter()
	router = setRESTRoutes(db, cnf, router)

	return router
}

// setRESTRoutes specifies all public routes for the comment service
func setRESTRoutes(db *sql.DB, cnf config.Config, router *mux.Router) *mux.Router {

	users := router.PathPrefix("/products").Subrouter()
	users.Methods("OPTIONS").Handler(negroni.New(
		negroni.HandlerFunc(middleware.AcceptOPTIONS),
	))

	users.Methods("GET").Handler(negroni.New(
		negroni.HandlerFunc(middleware.AccessControlHandler),
		controllers.GetProductsHandler(db, cnf),
	))

	return router
}
