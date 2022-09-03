package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PaymentDBResponse struct {
	ID              primitive.ObjectID  	`json:"id" bson:"_id"`
	PaymentID       string				  	`json:"paymentId" bson:"paymentId"`
	UserID			primitive.ObjectID 		`json:"userId" bson:"userId"`
	Quantity       	int 					`json:"quantity" bson:"quantity"`
	Price			int 					`json:"price" bson:"price"`
	DateOfPurchase  time.Time          		`json:"dateOfPurchase,omitempty" bson:"dateOfPurchase,omitempty"`
}

type PaymentData struct {
	ID              string  	`json:"id" bson:"id"`
}

type SuccessPaymentInput struct {
	PaymentData PaymentData `json:"paymentData" bson:"paymentData"`
}