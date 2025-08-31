package controller

import (
	"context"
	"net/http"
	"time"

	config "challenge-project/config"
	lib "challenge-project/lib"
	productMdl "challenge-project/model/products"
	socketService "challenge-project/services/socket"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const productCollectionName = "products"

func CreateProduct(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	var p productMdl.Product
	if err := c.Bind(&p); err != nil {
		return lib.JSONError(c, http.StatusBadRequest, "Invalid Params!")
	}

	p.ID = primitive.NewObjectID()
	now := time.Now().Unix()
	p.CreatedAt = now
	p.UpdatedAt = now

	db := config.GetDatabase()
	_, err := db.Collection(productCollectionName).InsertOne(ctx, p)
	if err != nil {
		return lib.JSONError(c, http.StatusInternalServerError, "Internal DB Error!")
	}

	socketService.EmitToNamespace(config.GetConfig().Socket.Namespace, "Product Created!", p)

	return lib.JSONCreated(c, p)
}

func ListProducts(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var results []productMdl.Product

	db := config.GetDatabase()
	cursor, err := db.Collection(productCollectionName).Find(ctx, bson.M{})
	if err != nil {
		return lib.JSONError(c, http.StatusInternalServerError, "Internal DB Error!")
	}

	if err := cursor.All(ctx, &results); err != nil {
		return lib.JSONError(c, http.StatusInternalServerError, "Invalid Decode!")
	}

	return lib.JSONSuccess(c, results)
}

func GetProduct(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var p productMdl.Product
	id := c.Param("id")

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return lib.JSONError(c, http.StatusBadRequest, "invalid id")
	}

	db := config.GetDatabase()
	if err := db.Collection(productCollectionName).FindOne(ctx, bson.M{"_id": oid}).Decode(&p); err != nil {
		return lib.JSONError(c, http.StatusNotFound, "Product Not Found!")
	}

	return lib.JSONSuccess(c, p)
}

func UpdateProduct(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var p, updated productMdl.Product
	id := c.Param("id")

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return lib.JSONError(c, http.StatusBadRequest, "Invalid Param!")
	}

	if err := c.Bind(&p); err != nil {
		return lib.JSONError(c, http.StatusBadRequest, "Invalid Param!")
	}

	p.UpdatedAt = time.Now().Unix()
	update := bson.M{
		"$set": bson.M{
			"name":        p.Name,
			"category":    p.Category,
			"description": p.Description,
			"price":       p.Price,
			"amount":      p.Amount,
			"updated_at":  p.UpdatedAt,
		},
	}

	db := config.GetDatabase()
	res, err := db.Collection(productCollectionName).UpdateByID(ctx, oid, update)
	if err != nil || res.MatchedCount == 0 {
		return lib.JSONError(c, http.StatusInternalServerError, "Internal DB Error!")
	}

	_ = db.Collection(productCollectionName).FindOne(ctx, bson.M{"_id": oid}).Decode(&updated)

	socketService.EmitToNamespace(config.GetConfig().Socket.Namespace, "Product Updated", updated)

	return lib.JSONSuccess(c, updated)
}

func DeleteProduct(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id := c.Param("id")

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return lib.JSONError(c, http.StatusBadRequest, "Invalid Param!")
	}

	db := config.GetDatabase()
	res, err := db.Collection(productCollectionName).DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil || res.DeletedCount == 0 {
		return lib.JSONError(c, http.StatusInternalServerError, "Internal DB Error!")
	}

	socketService.EmitToNamespace(config.GetConfig().Socket.Namespace, "Product Deleted!", map[string]string{"id": id})

	return lib.JSONSuccess(c, map[string]string{"id": id})
}
