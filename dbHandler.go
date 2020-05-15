package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"local.org/complexlog"
)

const (
	DBHOST = "127.0.0.1"
	DBPORT = "27017"
)

type MongoDb struct {
	ctx        context.Context
	cancel     func()
	client     *mongo.Client
	db         *mongo.Database
	collection map[string]*mongo.Collection
	col        *mongo.Collection
	err        error
}

func (mdb *MongoDb) connect() {
	mdb.ctx, mdb.cancel = context.WithTimeout(context.Background(), 10*time.Second)
	mdb.client, mdb.err = mongo.Connect(mdb.ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	mdb.db = mdb.client.Database("blog")

	mdb.collection = make(map[string]*mongo.Collection)
	mdb.collection["post"] = mdb.db.Collection("post")

	if mdb.err != nil {
		log.Fatal(mdb.err)
	}
	//defer mdb.client.Disconnect(mdb.ctx)
	//defer mdb.cancel()

	complexlog.Dblog("init")
}

func (mdb *MongoDb) createData() {

	// insert mock data in db
	post := Post{
		Date: "1/1/2222",
		View: 1,
		Text: "text 1",
	}

	_, err := mdb.collection["post"].InsertOne(mdb.ctx, post)
	if err != nil {
		panic(err)
	}

}

func (mdb *MongoDb) getPost() {
	cursor, err := mdb.collection["post"].Find(mdb.ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	var post []postJs
	if err = cursor.All(mdb.ctx, &post); err != nil {
		log.Fatal("err")
	}
	fmt.Println(post)
}
