package controllers

import (
	"database/sql"
	"net/http"

	"github.com/juang77/GoMSArch/product-service/config"
	db "github.com/juang77/GoMSArch/product-service/database"
	util "github.com/juang77/GoMSArch/shared/util"
	"github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
)

func GetProductsHandler(connection *sql.DB, cnf config.Config) negroni.HandlerFunc {
	return negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		logrus.Info("List Products")

		products, err := db.GetProducts(connection)
		if err != nil {
			util.SendError(w, err)
			return
		}

		util.SendOK(w, products)
	})
}
