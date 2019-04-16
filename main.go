package main

import (
	"net/http"
	"os"
	"strconv"

	"./cache"
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
	configuration, err := configManager.Load("./config.yml", "yaml")
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
	var cacheManager cache.Manager
	err = cacheManager.Init(configuration.Cache)
	if err != nil {
		log.Error("Error while loading cache manager: ", err)
		os.Exit(3)
	}

	//handler for API
	var h Handler
	h.Init(dbManager, cacheManager)

	// Declare a new router
	r := LoadRouter(h)

	http.ListenAndServe(":"+strconv.Itoa(configuration.Server.Port), r)
	log.Info("Server runnning on port: ", strconv.Itoa(configuration.Server.Port))
}

//LoadRouter returns a router instance
func LoadRouter(h Handler) *mux.Router {
	router := mux.NewRouter()
	//declare routes
	router.HandleFunc("/employee", h.addEmployee).Methods("POST")
	router.HandleFunc("/employee/{id}", h.getEmployeeByID).Methods("GET")
	router.Path("/employee/sex/{sex}").Queries("lastid", "{lastid}").Queries("limit", "{limit}").HandlerFunc(h.getEmployeeBySex).Methods("GET")
	return router
}
