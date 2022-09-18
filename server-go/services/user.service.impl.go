package services

import (
	"context"
	"errors"
	"strings"
	"time"

	"example/web-service-gin/models"
	"example/web-service-gin/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserServiceImpl struct {
	collection *mongo.Collection
	cartCollection *mongo.Collection
	productCollection *mongo.Collection
	paymentCollection *mongo.Collection
	ctx        context.Context
}

func NewUserServiceImpl(collection *mongo.Collection, cartCollection *mongo.Collection, 
	productCollection *mongo.Collection, paymentCollection *mongo.Collection, ctx context.Context) UserService {
	return &UserServiceImpl{collection, cartCollection, productCollection, paymentCollection, ctx}
}

// FindUserByID
func (us *UserServiceImpl) FindUserById(id string) (*models.DBResponse, error) {
	oid, _ := primitive.ObjectIDFromHex(id)

	var user *models.DBResponse

	query := bson.M{"_id": oid}
	err := us.collection.FindOne(us.ctx, query).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &models.DBResponse{}, err
		}
		return nil, err
	}

	return user, nil
}

// FindUserByEmail
func (us *UserServiceImpl) FindUserByEmail(email string) (*models.DBResponse, error) {
	var user *models.DBResponse

	query := bson.M{"email": strings.ToLower(email)}
	err := us.collection.FindOne(us.ctx, query).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &models.DBResponse{}, err
		}
		return nil, err
	}

	return user, nil
}

// UpsertUser
func (uc *UserServiceImpl) UpsertUser(email string, data *models.UpdateDBUser) (*models.DBResponse, error) {
	doc, err := utils.ToDoc(data)
	if err != nil {
		return nil, err
	}

	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(1)
	query := bson.D{{Key: "email", Value: email}}
	update := bson.D{{Key: "$set", Value: doc}}
	res := uc.collection.FindOneAndUpdate(uc.ctx, query, update, opts)

	var updatedUser *models.DBResponse

	if err := res.Decode(&updatedUser); err != nil {
		return nil, errors.New("no post with that Id exists")
	}

	var cart models.CartDBResponse
	cart.ID = primitive.NewObjectID()
	cart.UserID = updatedUser.ID
	cart.Items = []models.CartItemDBResponse {}
	uc.cartCollection.InsertOne(uc.ctx, &cart)

	return updatedUser, nil
}

func (us *UserServiceImpl) AddToCart(userOid primitive.ObjectID, productId string) (*models.CartDBResponse, error) {

	productOid, _ := primitive.ObjectIDFromHex(productId)

	var product *models.ProductDBResponse
	query := bson.M{"_id": productOid}
	us.productCollection.FindOne(us.ctx, query).Decode(&product)

	var cart *models.CartDBResponse
	query = bson.M{"userId": userOid}
	us.cartCollection.FindOne(us.ctx, query).Decode(&cart)


	var cartItem *models.CartItemDBResponse = getItemFromCart(cart, productOid)

	if cartItem != nil {
		var quantity = cartItem.Quantity + 1
		var price = cartItem.Price + product.Price

		arrayFilters := options.ArrayFilters{Filters: bson.A{bson.M{"x.productId": cartItem.ProductID}}}
		upsert := true
		opts := options.UpdateOptions{
			ArrayFilters: &arrayFilters,
			Upsert:       &upsert,
		}
		update := bson.M{
			"$set": bson.M{
				"items.$[x].price": price, "items.$[x].quantity": quantity,
			},
		}
		us.cartCollection.UpdateOne(us.ctx, query, update, &opts)

		var updatedCart *models.CartDBResponse

		return updatedCart, nil

	} else {
		var cartItem models.CartItemDBResponse
		cartItem.ProductID = productOid
		cartItem.ProductName = product.Title
		cartItem.Price = product.Price
		cartItem.Quantity = 1

		opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(1)
		query := bson.D{{Key: "userId", Value: userOid}}
		update := bson.M{"$push":bson.M{"items": cartItem}}
		us.cartCollection.FindOneAndUpdate(us.ctx, query, update, opts)

		var updatedCart *models.CartDBResponse

		return updatedCart, nil
	}
}

func (us *UserServiceImpl) RemoveFromCart(userOid primitive.ObjectID, productId string) (*models.CartDBResponse, error) {

	productOid, _ := primitive.ObjectIDFromHex(productId)

	var product *models.ProductDBResponse
	query := bson.M{"_id": productOid}
	us.productCollection.FindOne(us.ctx, query).Decode(&product)

	var cart *models.CartDBResponse
	query = bson.M{"userId": userOid}
	us.cartCollection.FindOne(us.ctx, query).Decode(&cart)

	var cartItem *models.CartItemDBResponse = getItemFromCart(cart, productOid)

	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(1)
	update := bson.M{"$pull":bson.M{"items": cartItem}}
	res := us.cartCollection.FindOneAndUpdate(us.ctx, query, update, opts)

	var updatedCart *models.CartDBResponse

	if err := res.Decode(&updatedCart); err != nil {
		return nil, errors.New("cannot get cart from DB")
	}

	return updatedCart, nil
}

func (us *UserServiceImpl) GetCart(userOid primitive.ObjectID) (*models.CartDBResponse, error) {

	var cart *models.CartDBResponse

	query := bson.M{"userId": userOid}
	err := us.cartCollection.FindOne(us.ctx, query).Decode(&cart)

	if err != nil {
		return nil, err
	}

	return cart, nil
}

func (us *UserServiceImpl) CreatePaymentHistory(userOid primitive.ObjectID, paymentData *models.SuccessPaymentInput) (bool, error) {

	var cart *models.CartDBResponse

	query := bson.M{"userId": userOid}
	err := us.cartCollection.FindOne(us.ctx, query).Decode(&cart)

	for _, item := range cart.Items {
        var paymentHistory models.PaymentDBResponse
		paymentHistory.ID = primitive.NewObjectID()
		paymentHistory.PaymentID = paymentData.PaymentData.ID
		paymentHistory.UserID = userOid
		paymentHistory.DateOfPurchase = time.Now()
		paymentHistory.Price = item.Price
		paymentHistory.Quantity = item.Quantity
		us.paymentCollection.InsertOne(us.ctx, paymentHistory)
    }

	if err != nil {
		return false, err
	}

	return true, nil
}

func (us *UserServiceImpl) ClearCart(userOid primitive.ObjectID) (*models.CartDBResponse, error) {

	var cart *models.CartDBResponse

	query := bson.M{"userId": userOid}
	err := us.cartCollection.FindOne(us.ctx, query).Decode(&cart)

	for _, item := range cart.Items {
		var product *models.ProductDBResponse

		query := bson.M{"_id": item.ProductID}
		err := us.productCollection.FindOne(us.ctx, query).Decode(&product)

		if err == nil {
			opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(1)
			update := bson.M{"$set":bson.M{"sold": product.Sold + 1}}	
			us.productCollection.FindOneAndUpdate(us.ctx, query, update, opts)
		}
    }

	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(1)
	update := bson.M{"$set":bson.M{"items": []*models.CartItemDBResponse {}}}
	us.cartCollection.FindOneAndUpdate(us.ctx, query, update, opts)

	if err != nil {
		return nil, err
	}

	return cart, nil
}

func (us *UserServiceImpl) GetHistory(userOid primitive.ObjectID) ([]models.PaymentDBResponse, error) {

	var paymentHistory []models.PaymentDBResponse

	filter := bson.M{"userId": userOid}
	findOptions := options.Find()

	cursor, _ := us.paymentCollection.Find(us.ctx, filter, findOptions)

	defer cursor.Close(us.ctx)

	for cursor.Next(us.ctx) {
		var payment models.PaymentDBResponse
		cursor.Decode(&payment)
		paymentHistory = append(paymentHistory, payment)
	}

	return paymentHistory, nil
}

func getItemFromCart(cart *models.CartDBResponse, productId primitive.ObjectID) (cartItem *models.CartItemDBResponse) {
    for _, item := range cart.Items {
        if item.ProductID == productId {
            return &item
        }
    }
    return nil
}