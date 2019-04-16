package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var localconfiguration = &Configuration{
	Database: DatabaseConfiguration{Host: "192.168.99.100", Port: 27017, Username: "", Password: "", Name: "mydatabase"},
	Cache:    CacheConfiguration{Memory: 256, CRS: "LRU"},
	Server:   ServerConfiguration{Port: 8080},
}

var envconfiguration = &Configuration{
	Database: DatabaseConfiguration{Host: "10.55.22.196", Port: 37017, Username: "shubham", Password: "verekar", Name: "mydatabase"},
	Cache:    CacheConfiguration{Memory: 512, CRS: "LRU"},
	Server:   ServerConfiguration{Port: 9090},
}

func TestLocalConfiguration(t *testing.T) {
	var configManager Manager
	configuration, _ := configManager.Load("", "yaml")
	require.Equal(t, configuration, localconfiguration, "configuration read from yaml should be equal to local configuration object")
}

func TestLocalCustomPathConfiguration(t *testing.T) {
	var path = "./config_test/config.yml"
	var configManager Manager
	configuration, _ := configManager.Load(path, "yaml")
	require.Equal(t, configuration, localconfiguration, "configuration read from yaml should be equal to local configuration object")
}

func TestEnvConfiguration(t *testing.T) {
	//set envrionment properties
	os.Setenv(databaseHost, "10.55.22.196")
	os.Setenv(databasePort, "37017")
	os.Setenv(databaseUsername, "shubham")
	os.Setenv(databasePassword, "verekar")
	os.Setenv(databaseName, "mydatabase")
	os.Setenv(serverPort, "9090")
	os.Setenv(cacheMem, "512")
	os.Setenv(cacheReplacementStatergy, "LRU")

	var configManager Manager
	configuration, _ := configManager.Load("", "yaml")
	require.Equal(t, configuration, envconfiguration, "configuration should be equal to env configuration object")
	require.Equal(t, configManager.getConfiguration(), envconfiguration, "configuration should be equal to env configuration object")
}
