package model

//Employee model
type Employee struct {
	ID        int32  `bson:"_id" json:"id"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Age       int    `json:"age"`
	Sex       string `json:"sex"`
}
