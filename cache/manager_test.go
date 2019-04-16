package cache

import (
	"testing"

	"../config"
	"../model"
	"github.com/stretchr/testify/require"
)

//sample local configuration
var localconfiguration = &config.Configuration{
	Database: config.DatabaseConfiguration{Host: "192.168.99.100", Port: 27017, Username: "", Password: "", Name: "mydatabase"},
	Cache:    config.CacheConfiguration{Memory: 256},
	Server:   config.ServerConfiguration{Port: 8080},
}

var cacheManager Manager

var employees = [4]model.Employee{
	model.Employee{ID: 782976, FirstName: "Mike", LastName: "Yale", Age: 25, Sex: "M"},
	model.Employee{ID: 782977, FirstName: "Tim", LastName: "Kane", Age: 29, Sex: "F"},
	model.Employee{ID: 782978, FirstName: "Alice", LastName: "Jane", Age: 28, Sex: "M"},
	model.Employee{ID: 782979, FirstName: "Bob", LastName: "Smith", Age: 29, Sex: "M"},
}

//TestItemSet sets items in cache Memory
func TestItemSet(t *testing.T) {
	cacheManager.Init(localconfiguration.Cache)
	for _, emp := range employees {
		err := cacheManager.AddItem(emp)
		require.Equal(t, nil, err)
	}
}

//TestItemGet gets object from cache memory
func TestItemGet(t *testing.T) {

	for _, emp := range employees {
		employee, err := cacheManager.GetItem(emp.ID)
		require.Equal(t, nil, err)
		require.Equal(t, emp, *employee)
	}
}
