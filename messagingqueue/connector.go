package messagingqueue

import (
	"github.com/streadway/amqp"
)

//Connector handles connection operation to the mq
// TODO: handle retry mechanism efficiently
type Connector struct {
	uri string
}

//Connect connects to the messaging queue instance
func (c *Connector) Connect(uri string) (*amqp.Connection, error) {

	c.uri = uri
	conn, err := amqp.Dial(c.uri)
	if err != nil {
		log.Error("Error connecting messing queue: ", err)
		return nil, err
	}
	defer conn.Close()
	return conn, err
}

//Retry retries connection to mongoDB
func (c *Connector) Retry() (*amqp.Connectiont, error) {

	conn, err := amqp.Dial(c.uri)
	if err != nil {
		log.Error("Error connecting messing queue: ", err)
		return nil, err
	}
	defer conn.Close()
	return conn, err
}
