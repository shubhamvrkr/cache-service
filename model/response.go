package model

//Response for pagination query
type Response struct {
	Employees []Employee `json:"employees"`
	Next      string     `json:"next"`
}
