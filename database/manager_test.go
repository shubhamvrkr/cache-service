package database

import (
	"testing"

	"../config"
	"../model"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

var dbManager Manager

var employees = [4]model.Employee{
	model.Employee{ID: "782976", FirstName: "Shubham", LastName: "Verekar", Age: 25, Sex: "M"},
	model.Employee{ID: "782977", FirstName: "Shalini", LastName: "Gupta", Age: 29, Sex: "F"},
	model.Employee{ID: "782978", FirstName: "Shubhabrata", LastName: "Mukherjee", Age: 28, Sex: "M"},
	model.Employee{ID: "782979", FirstName: "Niket", LastName: "Yende", Age: 29, Sex: "M"},
}

var localconfiguration = &config.Configuration{
	Database: config.DatabaseConfiguration{Host: "192.168.99.100", Port: 27017, Username: "", Password: "", Name: "mydatabase"},
	Server:   config.ServerConfiguration{Port: 8080},
}

func TestDatabaseInsert(t *testing.T) {
	dbManager.Init(localconfiguration.Database)
	for _, employee := range employees {
		err := dbManager.Save(bson.M{"_id": employee.ID, "FirstName": employee.FirstName, "LastName": employee.LastName, "Age": employee.Age, "Sex": employee.Sex})
		require.Equal(t, nil, err)
	}
}

func TestDatabaseReadAll(t *testing.T) {
	var M map[string]interface{}
	dbManager.Init(localconfiguration.Database)
	emps, err := dbManager.Find(M)
	if err != nil {
		t.Error("Error: ", err)
	}
	for i, emp := range *emps {
		require.Equal(t, employees[i], emp)
	}
}

func TestDatabaseRead(t *testing.T) {
	emp, err := dbManager.Fetch(bson.M{"_id": "782976"})
	if err != nil {
		t.Error("Error: ", err)
	}
	require.Equal(t, employees[0], *emp)
}

func TestDatabaseReadPagination(t *testing.T) {

	query := bson.M{
		"Sex": "F",
	}
	emps, err := dbManager.Find(query)
	if err != nil {
		t.Error("Error: ", err)
	}
	require.Equal(t, 1, len(*emps))
	for _, emp := range *emps {
		require.Equal(t, employees[1], emp)
	}
}
