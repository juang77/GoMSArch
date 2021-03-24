package main

import (
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/juang77/LaikaGoMS/product-service/app/http/routes"
	"github.com/juang77/LaikaGoMS/product-service/config"
	db "github.com/juang77/LaikaGoMS/product-service/database"
	negronilogrus "github.com/meatballhat/negroni-logrus"
	"github.com/urfave/negroni"

	// go-sql-driver/mysql is needed for the database connection
	_ "github.com/go-sql-driver/mysql"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	log.SetLevel(log.DebugLevel)

	//Get config
	// cnf := config.Config{}
	// cnf.DBHost = "192.168.0.13"
	// cnf.DBUsername = "conectapis"
	// cnf.DBPassword = "Nicolas8032367."
	// cnf.Port = 5005
	// cnf.DBPort = 3306
	// cnf.Database = "laikadblinux"

	// PORT:5004
	// USERS_SERVICE_URL:users
	// DB_USERNAME:conectapis
	// DB_PASSWORD:Nicolas8032367.
	// DB_HOST:192.168.0.13
	// DB_PORT:3306
	// DB:laikadblinux
	// SECRET_KEY:Juancho

	cnf := config.LoadConfig()

	connection, err := db.OpenConnection(cnf)
	if err != nil {
		log.Fatal(err)
	}

	defer db.CloseConnection(connection)

	routes := routes.InitRoutes(connection, cnf)
	n := negroni.Classic()
	n.Use(negronilogrus.NewMiddleware())
	n.UseHandler(routes)

	// Start and listen on port in cbf.Port
	log.Info("Starting server on port " + strconv.Itoa(cnf.Port))
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(cnf.Port), n))

}
