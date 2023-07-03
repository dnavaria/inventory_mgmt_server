package main

import (
	"ims/inventory_mgmt_server/ims_server"
)

func main() {
	mongo_uri := "mongodb://admin:ajUEKK7lBK087Lln8Uj6j@localhost:27017/"
	app := ims_server.App{}
	app.Initialize(mongo_uri)
	app.Run(":9000")
	defer app.MongoClient.Disconnect(*app.CTX)
}
