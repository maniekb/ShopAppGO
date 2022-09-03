package services

import "example/web-service-gin/models"

type ProductsService interface {
	GetProducts(string, int64, int64, []int, int, int) ([]*models.ProductDBResponse, int64, error)
	GetProduct(string) (*models.ProductDBResponse, error)
	UploadProduct(*models.ProductUploadInput) (*models.ProductDBResponse, error)
}