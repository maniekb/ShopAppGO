package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductDBResponse struct {
	ID              primitive.ObjectID `json:"id" bson:"_id"`
	Price           int            	   `json:"price" bson:"price"`
	Manufacturer    int            	   `json:"manufacturer" bson:"manufacturer"`
	Views        	int            	   `json:"views" bson:"views"`
	Writer			primitive.ObjectID `json:"writer" bson:"writer"`
	Title        	string             `json:"title" bson:"title"`
	Description     string             `json:"description" bson:"description"`
	Sold        	bool               `json:"sold" bson:"sold"`
	CreatedAt       time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at" bson:"updated_at"`
}

type ProductUploadInput struct {
	Price           int             `json:"price" bson:"price"`
	Manufacturer    int             `json:"manufacturer" bson:"manufacturer"`
	Writer			primitive.ObjectID `json:"writer" bson:"writer"`
	Title        	string             `json:"title" bson:"title"`
	Description     string             `json:"description" bson:"description"`
}

type ProductCartItem struct {
	Price			int             `json:"price" bson:"price"`
	Quantity		int             `json:"quantity" bson:"quantity"`
}