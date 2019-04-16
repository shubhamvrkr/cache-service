package database

import (
	"context"
	"strconv"
	"time"

	"../config"
	"../model"
	"github.com/op/go-logging"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

var log = logging.MustGetLogger("database_manager")

var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

//Manager handles database operations
type Manager struct {
	database   config.DatabaseConfiguration
	client     *mongo.Client
	collection *mongo.Collection
}

//Init establishes conenction with database instances and creates necessary database and collection
func (m *Manager) Init(database config.DatabaseConfiguration) error {
	var connector Connector
	var client *mongo.Client
	var err error

	m.database = database
	if len(m.database.Username) > 0 {
		credentials := options.Credential{Username: m.database.Username, Password: m.database.Password, PasswordSet: true}
		client, err = connector.Connect("mongodb://"+m.database.Host+":"+strconv.Itoa(m.database.Port), options.Client().SetAuth(credentials))
	} else {
		client, err = connector.Connect("mongodb://"+m.database.Host+":"+strconv.Itoa(m.database.Port), options.Client())
	}
	if err != nil {
		log.Info("Error intializing database manager: ", err)
		return err
	}
	m.collection = client.Database(m.database.Name).Collection("employees")
	//create index for colelction
	indexView := m.collection.Indexes()
	indexView.CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys: bsonx.Doc{{Key: "Sex", Value: bsonx.Int32(1)}},
		},
	)
	m.client = client
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
	employee = model.Employee{ID: e["_id"].(string), FirstName: e["FirstName"].(string), LastName: e["LastName"].(string), Age: int(e["Age"].(int32)), Sex: e["Sex"].(string)}
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
		employees = append(employees, model.Employee{ID: e["_id"].(string), FirstName: e["FirstName"].(string), LastName: e["LastName"].(string), Age: int(e["Age"].(int32)), Sex: e["Sex"].(string)})
	}
	if err := cur.Err(); err != nil {
		log.Error("Error database cursor : ", err)
		return nil, err
	}
	return &employees, nil
}
