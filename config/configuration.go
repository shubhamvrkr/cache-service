package config

//Configuration holds the confirguration for the database, cache and server
type Configuration struct {
	//Holds database configuration
	Database DatabaseConfiguration
	//Holds cache configuration
	Cache CacheConfiguration
	//Rabbit configuration
	Rabbit MessageQueueConfiguration
	//Holds server configuration
	Server ServerConfiguration
}
