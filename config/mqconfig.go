package config

//MessageQueueConfiguration holds the configuration for the messaging queue
type MessageQueueConfiguration struct {
	//Hostname of the MQ to establish a connection
	Host string
	//Port of the MQ to establish a connection
	Port int
	//Credentials set for the MQ (opttional)
	Username string
	Password string
	//Name of the channel
	Queue string
}
