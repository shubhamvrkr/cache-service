package messagingqueue

import (
	"testing"

	"cacheservice/config"

	"github.com/stretchr/testify/require"
)

//default local configuration
var localconfiguration = &config.Configuration{
	Database: config.DatabaseConfiguration{Host: "localhost", Port: 27017, Username: "", Password: "", Name: "mydatabase"},
	Cache:    config.CacheConfiguration{Memory: 256},
	Rabbit:   config.MessageQueueConfiguration{Host: "localhost", Port: 5672, Username: "", Password: "", Queue: "events"},
	Server:   config.ServerConfiguration{Port: 8080},
}

var mqManager Manager

func Test1(t *testing.T) {
	message := "reload"
	mqManager.Init(localconfiguration.Rabbit)
	msgs := mqManager.Consume()
	err := mqManager.Publish(message)
	msg := <-msgs
	log.Info("Recieved: ", string(msg.Body))
	require.Equal(t, nil, err)
	require.Equal(t, message, string(msg.Body))
}
