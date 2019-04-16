package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"./cache"
	"./config"
	"./database"
	"./model"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

var router *mux.Router

var employees = [4]model.Employee{
	model.Employee{ID: "782976", FirstName: "Mike", LastName: "Yale", Age: 25, Sex: "M"},
	model.Employee{ID: "782977", FirstName: "Tim", LastName: "Kane", Age: 29, Sex: "F"},
	model.Employee{ID: "782978", FirstName: "Alice", LastName: "Jane", Age: 28, Sex: "M"},
	model.Employee{ID: "782979", FirstName: "Bob", LastName: "Smith", Age: 29, Sex: "M"},
}

var sex string = "M"
var limit string = "2"
var lastid string = "-1"

func Test1(t *testing.T) {

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
	router = LoadRouter(h)
}

func Test2GetUnknownEmployee(t *testing.T) {
	unknownID := "0372097"
	request, _ := http.NewRequest("GET", "/employee/"+unknownID, nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	log.Info("Recieved from sever: ", response.Body)
	require.Equal(t, http.StatusNotFound, response.Code, "Employee not found")
}

func Test3GetEmptyEmployeesWithPagination(t *testing.T) {

	request, _ := http.NewRequest("GET", "/employee/sex"+sex+"?lastid="+lastid+"&limit="+limit, nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	log.Info("Recieved from sever: ", response.Body)
	require.Equal(t, http.StatusNotFound, response.Code, "No employees found")
}

func Test4AddEmployee(t *testing.T) {
	for _, emp := range employees {
		jsonEmployee, _ := json.Marshal(emp)
		request, _ := http.NewRequest("POST", "/employee", bytes.NewBuffer(jsonEmployee))
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)
		require.Equal(t, http.StatusCreated, response.Code, "Employee saved successfully")
	}
}

func Test5GetEmployee(t *testing.T) {
	var test model.Employee
	for _, emp := range employees {
		request, _ := http.NewRequest("GET", "/employee/"+emp.ID, nil)
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)
		log.Info("Recieved from sever: ", response.Body)
		json.Unmarshal(response.Body.Bytes(), &test)
		require.Equal(t, emp, test)
	}
}

func Test6GetEmployeesWithPagination(t *testing.T) {

	request, _ := http.NewRequest("GET", "/employee/sex"+sex+"?lastid="+lastid+"&limit="+limit, nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	log.Info("Recieved from sever: ", response.Body)
	require.Equal(t, http.StatusOK, response.Code)

}
