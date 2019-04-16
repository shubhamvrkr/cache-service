package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"./model"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

var employees = [4]model.Employee{
	model.Employee{ID: "782976", FirstName: "Mike", LastName: "Yale", Age: 25, Sex: "M"},
	model.Employee{ID: "782977", FirstName: "Tim", LastName: "Kane", Age: 29, Sex: "F"},
	model.Employee{ID: "782978", FirstName: "Alice", LastName: "Jane", Age: 28, Sex: "M"},
	model.Employee{ID: "782979", FirstName: "Bob", LastName: "Smith", Age: 29, Sex: "M"},
}

func Router() *mux.Router {

	var h Handler

	r := mux.NewRouter()

	//declare routes
	r.HandleFunc("/employee", h.addEmployee).Methods("POST")
	r.HandleFunc("/employee/proto", h.addEmployeeProto).Methods("POST")
	r.HandleFunc("/employee/{id}", h.getEmployeeByID).Methods("GET")
	r.HandleFunc("/employee/{sex}/sex", h.getEmployeeBySex).Methods("GET")
	return r

}

func TestAddEmployee(t *testing.T) {

	jsonEmployee, _ := json.Marshal(employees[0])
	request, _ := http.NewRequest("POST", "/employee", bytes.NewBuffer(jsonEmployee))
	response := httptest.NewRecorder()
	Router().ServeHTTP(response, request)
	assert.Equal(t, http.StatusCreated, response.Code, "Employee saved successfully")
}
