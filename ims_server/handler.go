package ims_server

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func getProducts(inventory *mongo.Collection) ([]Product, error) {
	// Create a new context
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	// Find all products
	cursor, err := inventory.Find(ctx, bson.M{})
	checkErr(err, "Error while getting products")

	defer cursor.Close(ctx)

	// Iterate through the cursor and decode each document
	var products []Product
	for cursor.Next(ctx) {
		var product Product
		err := cursor.Decode(&product)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

func getProduct(inventory *mongo.Collection, id int) (*Product, error) {
	// Create a new context
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	// Find the product
	var product Product
	err := inventory.FindOne(ctx, bson.M{"id": id}).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // No product found, return nil
		}
		return nil, err // Other error occurred
	}

	return &product, nil
}

func createProduct(inventory *mongo.Collection, product Product) (bool, error) {
	// Create a new context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Insert the product
	result, err := inventory.InsertOne(ctx, product)
	if err != nil {
		return false, err
	}
	if result.InsertedID == nil {
		return false, nil
	}

	return true, nil
}

func deleteProduct(inventory *mongo.Collection, product_id int) (bool, error) {
	// Create a new context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Delete the product
	result, err := inventory.DeleteOne(ctx, bson.M{"id": product_id})
	if err != nil {
		return false, err
	}
	if result.DeletedCount == 0 {
		return false, nil
	}
	return true, nil
}

func updateProduct(inventory *mongo.Collection, product_id int, product Product) (bool, error) {
	// Create a new context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Update the product
	result, err := inventory.UpdateOne(ctx, bson.M{"id": product_id}, bson.M{"$set": product})
	if err != nil {
		return false, err
	}
	if result.ModifiedCount == 0 {
		return false, nil
	}
	return true, nil
}

func updateProductQuantity(inventory *mongo.Collection, product_id int, quantity uint64) (bool, error) {
	// Create a new context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Update the product
	result, err := inventory.UpdateOne(ctx, bson.M{"id": product_id}, bson.M{"$set": bson.M{"quantity": quantity}})
	if err != nil {
		return false, err
	}
	if result.ModifiedCount == 0 {
		return false, nil
	}
	return true, nil
}

func updateProductPrice(inventory *mongo.Collection, product_id int, price float64) (bool, error) {
	// Create a new context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Update the product
	result, err := inventory.UpdateOne(ctx, bson.M{"id": product_id}, bson.M{"$set": bson.M{"price": price}})
	if err != nil {
		log.Println("updateProductPrice error :: error while updating product price :: ", err)
		return false, err
	}
	if result.ModifiedCount == 0 {
		log.Println("updateProductPrice error :: no product found with given id :: ", product_id)
		return false, nil
	}
	return true, nil
}

func updateProductName(inventory *mongo.Collection, product_id int, name string) (bool, error) {
	// Create a new context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Update the product
	result, err := inventory.UpdateOne(ctx, bson.M{"id": product_id}, bson.M{"$set": bson.M{"name": name}})
	if err != nil {
		return false, err
	}
	if result.ModifiedCount == 0 {
		return false, nil
	}
	return true, nil
}

func updateProductDescription(inventory *mongo.Collection, product_id int, description string) (bool, error) {
	// Create a new context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Update the product
	result, err := inventory.UpdateOne(ctx, bson.M{"id": product_id}, bson.M{"$set": bson.M{"description": description}})
	if err != nil {
		return false, err
	}
	if result.ModifiedCount == 0 {
		return false, nil
	}
	return true, nil
}
