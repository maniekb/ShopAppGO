package services

import (
	"context"
	"fmt"

	"example/web-service-gin/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProductsServiceImpl struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewProductsServiceImpl(collection *mongo.Collection, ctx context.Context) ProductsService {
	return &ProductsServiceImpl{collection, ctx}
}

// GetProducts
func (ps *ProductsServiceImpl) GetProducts(searchTerm string, skip int64, limit int64, 
	manufacturers []int, priceFrom int, priceTo int) ([]*models.ProductDBResponse, int64, error) {

	var products []*models.ProductDBResponse

	filter := bson.M{}
	findOptions := options.Find()

	if searchTerm != "" {
		filter = bson.M{
			"$or": []bson.M{
				{
					"title": bson.M{
						"$regex": primitive.Regex{
							Pattern: searchTerm,
							Options: "i",
						},
					},
				},
				{
					"description": bson.M{
						"$regex": primitive.Regex{
							Pattern: searchTerm,
							Options: "i",
						},
					},
				},
			},
		}
	}

	if len(manufacturers) > 0 {
		filter["manufacturer"] = bson.M{"$in": manufacturers}
	}

	if priceFrom > 0 || priceTo > 0 {
		filter["price"] = bson.M{
			"$gte": priceFrom,
			"$lt":  priceTo,
		}
	}

	total, _ := ps.collection.CountDocuments(ps.ctx, filter)

	findOptions.SetSkip(skip)
	findOptions.SetLimit(limit)

	cursor, _ := ps.collection.Find(ps.ctx, filter, findOptions)

	fmt.Println(cursor)
	defer cursor.Close(ps.ctx)

	for cursor.Next(ps.ctx) {
		var product *models.ProductDBResponse
		cursor.Decode(&product)
		products = append(products, product)
	}

	return products, total, nil
}

func (ps *ProductsServiceImpl) GetProduct(id string) (*models.ProductDBResponse, error) {
	oid, _ := primitive.ObjectIDFromHex(id)

	var product *models.ProductDBResponse

	query := bson.M{"_id": oid}
	err := ps.collection.FindOne(ps.ctx, query).Decode(&product)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &models.ProductDBResponse{}, err
		}
		return nil, err
	}

	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(1)
	update := bson.M{"$set":bson.M{"views": product.Views + 1}}	
	ps.collection.FindOneAndUpdate(ps.ctx, query, update, opts)

	return product, nil
}

func (ps *ProductsServiceImpl) UploadProduct(product *models.ProductUploadInput) (*models.ProductDBResponse, error) {
	res, err := ps.collection.InsertOne(ps.ctx, &product)

	if err != nil {
		return nil, err
	}

	var newProduct *models.ProductDBResponse
	query := bson.M{"_id": res.InsertedID}

	err = ps.collection.FindOne(ps.ctx, query).Decode(&newProduct)
	if err != nil {
		return nil, err
	}

	return newProduct, nil
}
