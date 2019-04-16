package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"./cache"
	"./database"
	"./model"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Handler handles all the API functionalities
type Handler struct {
	dbManager    database.Manager
	cacheManager cache.Manager
}

//Init initialized database manager and cache manager
func (h *Handler) Init(dbManager database.Manager, cacheManager cache.Manager) {
	h.dbManager = dbManager
	h.cacheManager = cacheManager
}

//addEmployee adds the employee in database and also to the cache
func (h *Handler) addEmployee(w http.ResponseWriter, r *http.Request) {

	var employee model.Employee
	//read request body
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&employee); err != nil {
		log.Error("Error decoding employee: ", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	//save to database
	err := h.dbManager.Save(employee)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//save to cache
	h.cacheManager.AddItem(employee)

	//return sucess
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Employee saved successfully")
}

//getEmployeeByID returns the employee details. First it checks in the cache and the in the dataabase if cache miss
func (h *Handler) getEmployeeByID(w http.ResponseWriter, r *http.Request) {
	var employee *model.Employee

	vars := mux.Vars(r)
	employeeID := vars["id"]

	log.Info("EmployeeID: ", employeeID)

	employee, err := h.cacheManager.GetItem(employeeID)
	if err != nil {
		//case: data not found in cache or some internal error
		log.Info("Error getting data from cache: ", err)
		//get data from database
		employee, err := h.dbManager.Fetch(bson.M{"_id": employeeID})
		if err != nil {
			log.Error("Error fetching data from database: ", err)
			if strings.ContainsAny(err.Error(), "no documents") {
				http.Error(w, "Employee not found", http.StatusNotFound)
				return
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

		} else {
			//add data to cache
			h.cacheManager.AddItem(*employee)
		}
	}
	empBytes, err := json.Marshal(employee)
	if err != nil {
		log.Error("Error marshalling emp: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(empBytes)
}

//getEmployeeBySex returns employees details based on the sex. (i.e M/F)
func (h *Handler) getEmployeeBySex(w http.ResponseWriter, r *http.Request) {
	var response model.Response
	vars := mux.Vars(r)
	sex := vars["sex"]
	lastid, err := strconv.Atoi(vars["lastid"])
	count, err := strconv.Atoi(vars["limit"])
	limit := int64(count)
	if err != nil {
		log.Error("Error converting count : ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Info("Sex: ", sex)
	log.Info("Lastid: ", lastid)
	log.Info("Limit: ", limit)

	//prepare query
	query := bson.M{"Sex": sex, "_id": bson.M{"$gt": lastid}}
	options := options.FindOptions{}
	options.Sort = bson.M{"_id": 1}
	options.Projection = bson.M{"_id": 1}
	options.Limit = &limit

	//get only employees ID
	emps, err := h.dbManager.FindEmployeeIds(query, &options)
	if err != nil {
		log.Error("Error fetching data from database: ", err)
		if strings.ContainsAny(err.Error(), "no documents") {
			http.Error(w, "No employees found", http.StatusNotFound)
			return
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	//need to check if this is last page
	response.Employees = *emps
	response.Next = "/employee/sex/" + sex + "?lastid=" + response.Employees[len(response.Employees)-1].ID + "&limit=" + strconv.Itoa(int(limit))
	resBytes, err := json.Marshal(response)
	if err != nil {
		log.Error("Error marshalling emp: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(resBytes)
}

//FindEmployee finds employee with id in cache, if miss load from db and update cache
func (h *Handler) FindEmployee(empids []string) *model.Employee {
	return nil
}
