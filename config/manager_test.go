package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

//default local configuration
var localconfiguration = &Configuration{
	Database: DatabaseConfiguration{Host: "localhost", Port: 27017, Username: "", Password: "", Name: "mydatabase"},
	Cache:    CacheConfiguration{Memory: 256},
	Rabbit:   MessageQueueConfiguration{Host: "localhost", Port: 5672, Username: "", Password: "", Queue: "events"},
	Server:   ServerConfiguration{Port: 8080},
}

//sample env configuration
var envconfiguration = &Configuration{
	Database: DatabaseConfiguration{Host: "10.55.22.196", Port: 37017, Username: "shubham", Password: "verekar", Name: "mydatabase"},
	Cache:    CacheConfiguration{Memory: 512},
	Rabbit:   MessageQueueConfiguration{Host: "10.55.22.196", Port: 5672, Username: "shubham", Password: "verekar", Queue: "events"},
	Server:   ServerConfiguration{Port: 9090},
}

//TestLocalConfiguration tests reading config file from default loocation
func TestLocalConfiguration(t *testing.T) {
	var configManager Manager
	configuration, _ := configManager.Load("", "yaml")
	require.Equal(t, configuration, localconfiguration, "configuration read from yaml should be equal to local configuration object")
}

//TestLocalCustomPathConfiguration tests reading config file from specified location
func TestLocalCustomPathConfiguration(t *testing.T) {
	var path = "./config_test/config.yml"
	var configManager Manager
	configuration, _ := configManager.Load(path, "yaml")
	require.Equal(t, configuration, localconfiguration, "configuration read from yaml should be equal to local configuration object")
}

//TestEnvConfiguration tests overiding of environment variables
func TestEnvConfiguration(t *testing.T) {
	//set envrionment properties
	os.Setenv(databaseHost, "10.55.22.196")
	os.Setenv(databasePort, "37017")
	os.Setenv(databaseUsername, "shubham")
	os.Setenv(databasePassword, "verekar")
	os.Setenv(databaseName, "mydatabase")
	os.Setenv(serverPort, "9090")
	os.Setenv(cacheMem, "512")
	os.Setenv(mqHost, "10.55.22.196")
	os.Setenv(mqPort, "5672")
	os.Setenv(mqUsername, "shubham")
	os.Setenv(mqPassword, "verekar")
	os.Setenv(mqQueue, "events")

	var configManager Manager
	configuration, _ := configManager.Load("", "yaml")
	require.Equal(t, configuration, envconfiguration, "configuration should be equal to env configuration object")
	require.Equal(t, configManager.getConfiguration(), envconfiguration, "configuration should be equal to env configuration object")
}
