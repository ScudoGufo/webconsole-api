package main

import (
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// general type

type appHandler func(w http.ResponseWriter, r *http.Request)

// data struct

type Comment struct {
	ID     string `bson:"_id,omitempty"`
	Author string `bson:"author,omitempty"`
	Date   string `bson:"date,omitempty"`
	Text   string `bson:"text,omitempty"`
}

type Author struct {
	ID   string `bson:"_id,omitempty"`
	Name string `bson:"name,omitempty"`
	Nick string `bson:"nick,omitempty"`
}

type Post struct {
	ID      primitive.ObjectID   `bson:"_id,omitempty"`
	Author  primitive.ObjectID   `bson:"author,omitempty"`
	Date    string               `bson:"date,omitempty"`
	View    int                  `bson:"view,omitempty"`
	Text    string               `bson:"text,omitempty"`
	Comment []primitive.ObjectID `bson:"comment,omitempty"`
}

//web response

type postJs struct {
	id     string
	author string
	views  int
}

type listJs struct {
	result int
	post   []postJs
}
