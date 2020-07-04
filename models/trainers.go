package models

type Trainer struct {
	Name string `json:"name,omitempty" bson:"name,omitempty"`
	Age  int    `json:"age,omitempty" bson:"age,omitempty"`
	City string `json:"city,omitempty" bson:"city,omitempty"`
}
