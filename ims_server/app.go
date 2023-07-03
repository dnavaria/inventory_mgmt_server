package ims_server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type App struct {
	Router      *mux.Router
	MongoClient *mongo.Client
	DB          *mongo.Database
	Collection  *mongo.Collection
	CTX         *context.Context
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Println(err)
		log.Fatal(msg)
	}
}

func (app *App) Initialize(mongo_uri string) error {
	// Create a new mongodb client
	connectionString := mongo_uri

	// Create a new context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	app.CTX = &ctx

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
	checkErr(err, "Error while connecting to database")

	// Assign the client to the app
	app.MongoClient = client

	// Assign the database to the app
	database := app.MongoClient.Database("ims_database")
	app.DB = database

	// Assign the collection to the app
	collection := app.DB.Collection("inventory")
	app.Collection = collection

	// Create a new router
	app.Router = mux.NewRouter().StrictSlash(true)
	app.handleRoutes()
	return nil
}

func (app *App) getProducts(w http.ResponseWriter, r *http.Request) {
	products, err := getProducts(app.Collection)
	if err != nil {
		log.Println("getProducts :: ", err, products)
		sendError(w, http.StatusInternalServerError, err.Error())
	}
	if products == nil {
		log.Println("getProducts :: ", err, products)
		sendError(w, http.StatusNotFound, "No products found")
		return
	}
	sendResponse(w, http.StatusOK, products)
}

func (app *App) getProduct(w http.ResponseWriter, r *http.Request) {
	// Get the product id from the url
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Println("getProduct :: ", err, id)
		sendError(w, http.StatusBadRequest, "Invalid product id")
		return
	}

	product, err := getProduct(app.Collection, id)
	if err != nil {
		log.Println("getProduct :: ", err, product)
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if product == nil {
		log.Println("getProduct :: ", err, product)
		sendError(w, http.StatusNotFound, "Product not found")
		return
	}
	sendResponse(w, http.StatusOK, product)
}

func (app *App) createProduct(w http.ResponseWriter, r *http.Request) {
	var product Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		msg := fmt.Sprintf("createProduct :: %v", err)
		log.Println(msg)
		sendError(w, http.StatusBadRequest, false)
		return
	}
	result, err := createProduct(app.Collection, product)
	if err != nil || !result {
		msg := fmt.Sprintf("createProduct :: %v :: %v", err, result)
		log.Println(msg)
		sendError(w, http.StatusInternalServerError, false)
		return
	}

	sendResponse(w, http.StatusCreated, result)
}

func (app *App) deleteProduct(w http.ResponseWriter, r *http.Request) {
	// Get the product id from the url
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Println("deleteProduct :: ", err, id)
		sendError(w, http.StatusBadRequest, "Invalid product id")
		return
	}
	result, err := deleteProduct(app.Collection, id)
	if err != nil || !result {
		log.Println("deleteProduct :: ", err, result, id)
		sendError(w, http.StatusInternalServerError, false)
		return
	}
	sendResponse(w, http.StatusOK, result)
}

func (app *App) updateProduct(w http.ResponseWriter, r *http.Request) {
	// Get the product id from the url
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Println("updateProduct :: ", err, id)
		sendError(w, http.StatusBadRequest, "Invalid product id")
		return
	}

	// get which field to update from headers
	field := r.Header.Get("field")
	if field == "" {
		log.Println("updateProduct :: ", err, field)
		sendError(w, http.StatusBadRequest, "Invalid field")
		return
	}

	// prepare the update query
	var product Product
	err = json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		log.Println("updateProduct :: error occured while decoding json body", err, product)
		sendError(w, http.StatusBadRequest, false)
		return
	}
	result, err := false, nil
	switch field {
	case "name":
		result, err = updateProductName(app.Collection, id, product.Name)
	case "quantity":
		result, err = updateProductQuantity(app.Collection, id, product.Quantity)
	case "price":
		result, err = updateProductPrice(app.Collection, id, product.Price)
	case "description":
		result, err = updateProductDescription(app.Collection, id, product.Description)
	default:
		result, err = updateProduct(app.Collection, id, product)
	}

	if err != nil || !result {
		log.Println("updateProduct :: ", err, result)
		sendError(w, http.StatusInternalServerError, false)
		return
	}
	sendResponse(w, http.StatusOK, result)
}

func (app *App) index(w http.ResponseWriter, r *http.Request) {
	sendResponse(w, http.StatusOK, "IMS API is running...")
}

func (app *App) handleRoutes() {
	app.Router.HandleFunc("/", app.index).Methods("GET")
	app.Router.HandleFunc("/inventory", app.getProducts).Methods("GET")
	app.Router.HandleFunc("/inventory/{id}", app.getProduct).Methods("GET")
	app.Router.HandleFunc("/inventory", app.createProduct).Methods("POST")
	app.Router.HandleFunc("/inventory/{id}", app.updateProduct).Methods("PUT")
	app.Router.HandleFunc("/inventory/{id}", app.deleteProduct).Methods("DELETE")
}

func (app *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, app.Router))
}

func sendResponse(w http.ResponseWriter, statusCode int, payload interface{}) {
	response, err := json.Marshal(payload)
	checkErr(err, "Error while marshalling response")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(response)
}

func sendError(w http.ResponseWriter, statusCode int, err interface{}) {
	errorString := make(map[string]interface{})
	if err != nil {
		errorString["error"] = err
	}
	sendResponse(w, statusCode, errorString)
}
