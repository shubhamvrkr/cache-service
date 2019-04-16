package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"./database"
	"./model"
	"github.com/gorilla/mux"
)

//Handler handles all the API functionalities
type Handler struct {
	dbManager database.Manager
}

//Init initialized database manager and cache manager
func (h *Handler) Init(dbManager database.Manager) {
	h.dbManager = dbManager
}

//addEmployee adds the employee in database and also to the cache
func (h *Handler) addEmployee(w http.ResponseWriter, r *http.Request) {
	var employee model.Employee
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&employee); err != nil {
		log.Error("Error decoding employee: ", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	//save to database
	//update cache

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Employee saved successfully")
}

//getEmployeeByID returns the employee details. First it checks in the cache and the in the dataabase if cache miss
func (h *Handler) getEmployeeByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	employeeID := vars["id"]
	log.Info("EmployeeID: ", employeeID)

	//find in cache
	//get from database
	//update cache
	w.WriteHeader(http.StatusNotFound)
	w.WriteHeader(http.StatusOK)

}

//getEmployeeBySex returns employees details based on the sex. (i.e M/F)
func (h *Handler) getEmployeeBySex(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sex := vars["sex"]
	log.Info("Sex: ", sex)

	w.WriteHeader(http.StatusNotFound)
	w.WriteHeader(http.StatusOK)
}
