package main

import (
	"net/http"
	"os"
	"strconv"

	"./cache"
	"./config"
	"./database"
	"./messagingqueue"
	"github.com/gorilla/mux"
	logging "github.com/op/go-logging"
	"github.com/rs/cors"
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
	}

	//Load messaging queue manager
	var mqManager messagingqueue.Manager
	err = mqManager.Init(configuration.Rabbit)
	if err != nil {
		log.Error("Error while loading message queue manager: ", err)
	}

	//handler for API
	var h Handler
	h.Init(dbManager, cacheManager, mqManager)

	// Declare a new router
	r := LoadRouter(h)
	handler := cors.Default().Handler(r)

	//Listens for the trigger from messaging queue
	go func() {
		msgs := mqManager.Consume()
		for {
			select {
			case msg := <-msgs:
				log.Info("Recieved event: ", string(msg.Body))
				h.ReloadCache()
			}
		}
	}()

	log.Info("Server runnning on port: ", strconv.Itoa(configuration.Server.Port))
	http.ListenAndServe(":"+strconv.Itoa(configuration.Server.Port), handler)

}

//LoadRouter returns a router instance
func LoadRouter(h Handler) *mux.Router {
	router := mux.NewRouter()
	//declare routes
	router.HandleFunc("/employee", h.AddEmployee).Methods("POST")
	router.HandleFunc("/employee/{id}", h.GetEmployeeByID).Methods("GET")
	router.HandleFunc("/reload", h.TriggerReloadEvent).Methods("GET")
	router.Path("/employee").Queries("sex", "{sex}").Queries("lastid", "{lastid}").Queries("limit", "{limit}").HandlerFunc(h.GetEmployeeBySex).Methods("GET")
	sh := http.StripPrefix("/swaggerui/", http.FileServer(http.Dir("./dist")))
	router.PathPrefix("/swaggerui/").Handler(sh)
	return router
}
