package main

import (
	"net/http"
	"os"
	"strconv"

	"./config"
	"./database"
	"github.com/gorilla/mux"
	logging "github.com/op/go-logging"
)

var log = logging.MustGetLogger("main")

var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

func main() {

	//load Configuration
	var configManager config.Manager
	configuration, err := configManager.Load("", "yaml")
	if err != nil {
		log.Error("Error while loading configuration: ", err)
		os.Exit(3)
	}

	//Load database manager
	var dbManager database.Manager
	err = dbManager.Init(configuration.Database)
	if err != nil {
		log.Error("Error while loading database manager: ", err)
		os.Exit(3)
	}

	//Load cache manager

	//handler for API
	var h Handler
	h.Init(dbManager)

	// Declare a new router
	r := mux.NewRouter()
	//declare routes
	r.HandleFunc("/employee", h.addEmployee).Methods("POST")
	r.HandleFunc("/employee/{id}", h.getEmployeeByID).Methods("GET")
	r.HandleFunc("/employee/{sex}/sex", h.getEmployeeBySex).Methods("GET")

	http.ListenAndServe(":"+strconv.Itoa(configuration.Server.Port), r)
	log.Info("Server runnning on port: ", strconv.Itoa(configuration.Server.Port))
}
