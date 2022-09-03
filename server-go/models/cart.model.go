package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CartItemDBResponse struct {
	ProductID          primitive.ObjectID  	`json:"productId" bson:"productId"`
	Price              int 					`json:"price" bson:"price"`
	Quantity		   int 					`json:"quantity" bson:"quantity"`
}

type CartDBResponse struct {
	ID              primitive.ObjectID  	`json:"id" bson:"_id"`
	UserID			primitive.ObjectID 		`json:"userId" bson:"userId"`
	Items       	[]CartItemDBResponse 		`json:"items" bson:"items"`
}