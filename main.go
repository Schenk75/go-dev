package main

import (
	"fmt"

	"go-dev/mongodb"
)

func main() {
	fmt.Printf("go-dev\n")
	conf := mongodb.MongoConfig{URI: "mongodb://127.0.0.1:27017"}
	client, err := mongodb.NewMongoClient(conf)
	if err != nil {
		panic(fmt.Sprintf("create mongo client fail: %v\n", err))
	}
	fmt.Printf("create mongo client success\n")

	id, err := client.GetID(mongodb.GoDev, mongodb.Test2)
	fmt.Printf("id:%d err:%v\n", id, err)
}
