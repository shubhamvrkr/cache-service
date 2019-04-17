//Package config provides the necessary configuration values
package config

import (
	"fmt"
	"os"
	"strconv"

	logging "github.com/op/go-logging"
	"github.com/spf13/viper"
)

var log = logging.MustGetLogger("config_manager")

var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

//env variable
var (
	databaseHost     = "DATABASE_HOST"
	databasePort     = "DATABASE_PORT"
	databaseUsername = "DATABASE_USER"
	databasePassword = "DATABASE_PASS"
	serverPort       = "SERVER_PORT"
	databaseName     = "DATABASE_NAME"
	cacheMem         = "CACHE_MEM"
	mqHost           = "MQ_HOST"
	mqPort           = "MQ_PORT"
	mqUsername       = "MQ_USERNAME"
	mqPassword       = "MQ_PASSWORD"
	mqQueue          = "MQ_QUEUE"
)

//Manager loads the configuration from specfied yaml configfile
//Incase environment flags are set for the configuration, environment variables overrides file configuration
type Manager struct {
	//File path of configuration yaml
	configfile string

	configuration Configuration
}

//Load loads the configuation from config.yml file or from enviroment variables and return Configuration instance
func (m *Manager) Load(filepath string, configtype string) (*Configuration, error) {

	var configuration Configuration
	var val string
	var ok bool

	if len(filepath) == 0 {
		//default location of configuration file
		m.configfile = "../config.yml"
		viper.AddConfigPath(".")

	} else {
		//set the custom file path
		m.configfile = filepath
	}
	//set the configuration properties
	viper.SetConfigType(configtype)
	viper.SetConfigFile(m.configfile)

	//read config file
	if err := viper.ReadInConfig(); err != nil {
		log.Error("Error reading config file: ", err)
		return nil, err
	}
	//marshal bytes into configuration object
	err := viper.Unmarshal(&configuration)
	if err != nil {
		log.Error("Unable to decode into configuration struct, %v", err)
		return nil, err
	}
	//override environment values
	//check database host
	val, ok = os.LookupEnv(databaseHost)
	if ok {
		configuration.Database.Host = val
	}
	//check database port
	val, ok = os.LookupEnv(databasePort)
	if ok {

		port, err := strconv.Atoi(val)
		if err != nil {
			log.Error("Unable to parse port<string> into port<uint8>, %v", err)
			return nil, err
		}
		configuration.Database.Port = port
	}
	//check database username
	val, ok = os.LookupEnv(databaseUsername)
	if ok {
		configuration.Database.Username = val
	}
	//check database password
	val, ok = os.LookupEnv(databasePassword)
	if ok {
		configuration.Database.Password = val
	}
	//check database name
	val, ok = os.LookupEnv(databaseName)
	if ok {
		configuration.Database.Name = val
	}
	//check server port
	val, ok = os.LookupEnv(serverPort)
	if ok {

		port, err := strconv.Atoi(val)
		if err != nil {
			log.Error("Unable to parse port<string> into port<uint8>, %v", err)
			return nil, err
		}
		configuration.Server.Port = port
	}

	//check cache memory
	val, ok = os.LookupEnv(cacheMem)
	if ok {

		mem, err := strconv.Atoi(val)
		if err != nil {
			log.Error("Unable to parse port<string> into port<uint8>, %v", err)
			return nil, err
		}
		configuration.Cache.Memory = mem
	}
	//Check mq host
	val, ok = os.LookupEnv(mqHost)
	if ok {
		configuration.Rabbit.Host = val
	}
	//Check mq port
	val, ok = os.LookupEnv(mqPort)
	if ok {

		port, err := strconv.Atoi(val)
		if err != nil {
			log.Error("Unable to parse port<string> into port<uint8>, %v", err)
			return nil, err
		}
		configuration.Rabbit.Port = port
	}
	//Check mq username
	val, ok = os.LookupEnv(mqUsername)
	if ok {
		configuration.Rabbit.Username = val
	}
	//check mq password
	val, ok = os.LookupEnv(mqPassword)
	if ok {
		configuration.Rabbit.Password = val
	}
	//Check mq queue
	val, ok = os.LookupEnv(mqQueue)
	if ok {
		configuration.Rabbit.Queue = val
	}

	log.Info("Database Configuration")
	log.Info("Host: ", configuration.Database.Host)
	log.Info("Port: ", configuration.Database.Port)
	log.Info("Username: ", configuration.Database.Username)
	log.Info("Password: ", configuration.Database.Password)
	log.Info("Cache Configuration")
	log.Info("Memory: ", configuration.Cache.Memory)
	log.Info("Message Queue Configuration")
	log.Info("Host: ", configuration.Rabbit.Host)
	log.Info("Port: ", configuration.Rabbit.Port)
	log.Info("Username: ", configuration.Rabbit.Username)
	log.Info("Password: ", configuration.Rabbit.Password)
	log.Info("Queue: ", configuration.Rabbit.Queue)
	log.Info("Server Configuration")
	log.Info("Port: ", configuration.Server.Port)
	fmt.Println()

	m.configuration = configuration
	return &configuration, nil
}

//Returns configuration
func (m *Manager) getConfiguration() *Configuration {
	return &m.configuration
}
