package services

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"example/web-service-gin/models"
)

type UserService interface {
	FindUserById(string) (*models.DBResponse, error)
	FindUserByEmail(string) (*models.DBResponse, error)
	UpsertUser(string, *models.UpdateDBUser) (*models.DBResponse, error)
	AddToCart(primitive.ObjectID, string) (*models.CartDBResponse, error)
	RemoveFromCart(primitive.ObjectID, string) (*models.CartDBResponse, error)
	GetCart(primitive.ObjectID) (*models.CartDBResponse, error)
	CreatePaymentHistory(primitive.ObjectID, *models.SuccessPaymentInput) (bool, error)
	ClearCart(primitive.ObjectID) (*models.CartDBResponse, error)
	GetHistory(primitive.ObjectID) ([]models.PaymentDBResponse, error)
}