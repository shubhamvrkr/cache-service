package messagingqueue

import (
	"testing"

	"../config"
	"../model"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var dbManager Manager

var employees = [4]model.Employee{
	model.Employee{ID: 782976, FirstName: "Mike", LastName: "Yale", Age: 25, Sex: "M"},
	model.Employee{ID: 782977, FirstName: "Tim", LastName: "Kane", Age: 29, Sex: "F"},
	model.Employee{ID: 782978, FirstName: "Alice", LastName: "Jane", Age: 28, Sex: "M"},
	model.Employee{ID: 782979, FirstName: "Bob", LastName: "Smith", Age: 29, Sex: "M"},
}

var localconfiguration = &config.Configuration{
	Database: config.DatabaseConfiguration{Host: "localhost", Port: 27017, Username: "", Password: "", Name: "mydatabase"},
	Cache:    config.CacheConfiguration{Memory: 256},
	Server:   config.ServerConfiguration{Port: 8080},
}

func TestDatabaseInsert(t *testing.T) {
	dbManager.Init(localconfiguration.Database)
	for _, employee := range employees {
		err := dbManager.Save(employee)
		require.Equal(t, nil, err)
	}
}

func TestDatabaseReadAll(t *testing.T) {
	var M map[string]interface{}
	dbManager.Init(localconfiguration.Database)
	options := options.FindOptions{}
	emps, err := dbManager.Find(M, &options)
	if err != nil {
		t.Error("Error: ", err)
	}
	for i, emp := range *emps {
		require.Equal(t, employees[i], emp)
	}
}

func TestDatabaseRead(t *testing.T) {
	emp, err := dbManager.Fetch(bson.M{"_id": 782976})
	if err != nil {
		t.Error("Error: ", err)
	}
	require.Equal(t, employees[0], *emp)
}

func TestDatabaseReadPagination(t *testing.T) {

	limit := int64(2)
	query := bson.M{"Sex": "M", "_id": bson.M{"$gt": -1}}
	options := options.FindOptions{}
	options.Sort = bson.M{"_id": 1}
	options.Projection = bson.M{"_id": 1}
	options.Limit = &limit

	emps, err := dbManager.FindEmployeeIds(query, &options)
	if err != nil {
		t.Error("Error: ", err)
	}
	require.Equal(t, 2, len(*emps))
}
