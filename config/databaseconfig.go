package config

//DatabaseConfiguration holds the configuration for the database
type DatabaseConfiguration struct {
	//Hostname of the database to establish a connection
	Host string
	//Port of the database to establish a connection
	Port int
	//Database name
	Name string
	//Credentials set for the database (opttional)
	Username string
	Password string
}
