package model

import "gopkg.in/mgo.v2/bson"

type (
    Deputy struct {
        Id bson.ObjectId `json:"id" bson:"_id"`
        Title string `json:"title" bson:"title"`
        FirstName string `json:"firstname" bson:"firstname"`
        LastName string `json:"lastname" bson:"lastname"`
        Province string `json:"province" bson:"province"`
    }
)
