package messagingqueue

import (
	"context"
	"strconv"
	"time"

	"../config"
	"../model"
	"github.com/op/go-logging"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	m.messagequeue = database
	if len(m.messagequeue.Username) > 0 {
		conn, err = connector.Connect("amqp://" + m.messagequeue.Username + ":" + m.messagequeue.Password + "@" + m.messagequeue.Host + ":" + strconv.Itoa(m.database.Port))
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

//Save saves element to the database
func (m *Manager) Save(employee model.Employee) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := m.collection.InsertOne(ctx, bson.M{"_id": employee.ID, "FirstName": employee.FirstName, "LastName": employee.LastName, "Age": employee.Age, "Sex": employee.Sex})
	if err != nil {
		log.Error("Error inserting element: ", err)
		return err
	}
	return nil
}

//Fetch gets the object from database by ID
func (m *Manager) Fetch(filter map[string]interface{}) (*model.Employee, error) {
	var e bson.M
	var employee model.Employee
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := m.collection.FindOne(ctx, filter).Decode(&e)
	if err != nil {
		log.Error("Error fetching element: ", err)
		return nil, err
	}
	employee = model.Employee{ID: e["_id"].(int32), FirstName: e["FirstName"].(string), LastName: e["LastName"].(string), Age: int(e["Age"].(int32)), Sex: e["Sex"].(string)}
	return &employee, nil
}

//Find retrieves objects from database by query criteria
func (m *Manager) Find(filter map[string]interface{}, opt *options.FindOptions) (*[]model.Employee, error) {
	var employees []model.Employee
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cur, err := m.collection.Find(ctx, filter, opt)
	if err != nil {
		log.Error("Error finding elements: ", err)
		return nil, err
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {

		var e bson.M
		err := cur.Decode(&e)
		if err != nil {
			log.Error("Error decoding object: ", err)
		}
		employees = append(employees, model.Employee{ID: e["_id"].(int32), FirstName: e["FirstName"].(string), LastName: e["LastName"].(string), Age: int(e["Age"].(int32)), Sex: e["Sex"].(string)})
	}
	if err := cur.Err(); err != nil {
		log.Error("Error database cursor : ", err)
		return nil, err
	}
	return &employees, nil
}

func (m *Manager) FindEmployeeIds(filter map[string]interface{}, opt *options.FindOptions) (*[]int32, error) {
	var employeesids []int32
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cur, err := m.collection.Find(ctx, filter, opt)
	if err != nil {
		log.Error("Error finding elements: ", err)
		return nil, err
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {

		var e bson.M
		err := cur.Decode(&e)
		if err != nil {
			log.Error("Error decoding object: ", err)
		}
		employeesids = append(employeesids, e["_id"].(int32))
	}
	if err := cur.Err(); err != nil {
		log.Error("Error database cursor : ", err)
		return nil, err
	}
	return &employeesids, nil
}
