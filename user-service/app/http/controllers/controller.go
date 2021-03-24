package controllers

import (
	"database/sql"
	"net/http"

	util "github.com/juang77/GoMSArch/shared/util"
	"github.com/juang77/GoMSArch/user-service/config"
	db "github.com/juang77/GoMSArch/user-service/database"
	"github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
)

func GetUsersHandler(connection *sql.DB, cnf config.Config) negroni.HandlerFunc {
	return negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		logrus.Info("List Users")

		users, err := db.GetUsers(connection)
		if err != nil {
			util.SendError(w, err)
			return
		}

		util.SendOK(w, users)
	})
}
