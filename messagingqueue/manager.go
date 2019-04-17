package messagingqueue

import (
	"strconv"

	"../config"
	"github.com/op/go-logging"
	"github.com/streadway/amqp"
)

var log = logging.MustGetLogger("messagingqueue_manager")

var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

//Manager handles messaging queue operations
type Manager struct {
	messagequeue config.MessageQueueConfiguration
	conn         *amqp.Connection
	channel      *amqp.Channel
}

//Init establishes conenction with message queue instances and create channel/queue
func (m *Manager) Init(messagequeue config.MessageQueueConfiguration) error {
	var connector Connector
	var conn *amqp.Connection
	var err error

	m.messagequeue = messagequeue
	if len(m.messagequeue.Username) > 0 {
		conn, err = connector.Connect("amqp://" + m.messagequeue.Username + ":" + m.messagequeue.Password + "@" + m.messagequeue.Host + ":" + strconv.Itoa(m.messagequeue.Port))
	} else {
		conn, err = connector.Connect("amqp://" + m.messagequeue.Host + ":" + strconv.Itoa(m.messagequeue.Port))
	}
	if err != nil {
		log.Info("Error intializing database manager: ", err)
		return err
	}
	m.conn = conn
	ch, err := conn.Channel()
	if err != nil {
		log.Info("Error creating channel: ", err)
		return err
	}
	m.channel = ch
	_, err = ch.QueueDeclare(
		messagequeue.Queue, // name
		false,              // durable
		false,              // delete when unused
		false,              // exclusive
		false,              // no-wait
		nil,                // arguments
	)
	if err != nil {
		log.Info("Error creating queue: ", err)
		return err
	}
	defer ch.Close()
	return nil

}

//Publish publishes messages to messgaing queue
func (m *Manager) Publish(message string) error {

	err := m.channel.Publish(
		"",                   // exchange
		m.messagequeue.Queue, // routing key
		false,                // mandatory
		false,                // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	if err != nil {
		log.Error("Error publishing message to queue: ", err)
		return err
	}
	return nil
}
