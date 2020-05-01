package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"local.org/complexlog"
)

const (
	DOMAIN = "api.test.local"
	DBHOST = "127.0.0.1"
	DBPORT = "27017"
)

var (
	COMMAND = map[string]func(map[string]string) string{
		"help":    help,
		"read":    read,
		"comment": comment,
		"login":   login,
		"echo":    echo,
	}
)

type comment struct {
	id string
	author string 
	date string
	text string
}

type author struct{
	id string
	name string
	nick string
}

type post struct {
	id string
	author author
	date string 
	view string 
	text string
	comment []comment
}

type appHandler func(w http.ResponseWriter, r *http.Request)

func main() {

	// mongo db connection
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
        log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	//get collection
	db := client.Database("blog")
	postsCol := db.Collection("post")
	commentCol := db.Collection("comment")

	// insert mock data in db 
	lastInteraction , err  := commentCol.InsertOne(ctx, []interface{}{
		bson.D{
			{"text","comment 1"},
			{"author","io"},
			{"data","2/1/1992"},
		}
	})

	lastInteraction , err  := postsCol.InsertOne(ctx, []interface{}{
		bson.D{
			{"data","1/1/1992"},
			{"text","post 1"},
			{"author","io"},
			{"comment",[lastInteraction.InsertedID]},
		}
	})

	// done

	//api router setting
	r := mux.NewRouter()
	r.Host(DOMAIN)

	r.HandleFunc("/cmd/{cmd}", withCORS(sendMethod)).Methods("OPTIONS")
	r.HandleFunc("/cmd/{cmd}", withCORS(runCmdApi)).Methods("POST")
	r.HandleFunc("/echo/{msg}", echoApi).Methods("GET")

	complexlog.Servlog("api init")

	//server web lintening
	err = http.ListenAndServe(":8000", r)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// Simple wrapper to Allow CORS.
func withCORS(fn appHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		complexlog.Apilog("Cors call")
		fn(w, r)
	}
}

func sendMethod(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	complexlog.Apilog("Option call")
}

func echoApi(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	json.NewEncoder(w).Encode(map[string]string{"msg": params["msg"]})

	complexlog.Apilog("Echo call")

}

func runCmdApi(w http.ResponseWriter, r *http.Request) {

	cmd := mux.Vars(r)["cmd"]

	var params map[string]string
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var result string
	if checkCmdPresence(cmd) {
		result = COMMAND[cmd](params)
	} else {
		result = wip(params)
	}

	json.NewEncoder(w).Encode(map[string]string{"result": result})

	complexlog.Apilog("Cmd call")

}

func checkCmdPresence(cmd string) bool {

	if COMMAND[cmd] != nil {
		return true
	}
	return false
}

/*
######################### command functions
*/

func help(args map[string]string) string {
	return "help"
}
func read(args map[string]string) string {
	return "read"
}
func comment(args map[string]string) string {
	return "comment"
}
func login(args map[string]string) string {
	if args["u"] == "test" && args["p"] == "test" {
		return "nice"
	}
	return "fail"
}
func echo(args map[string]string) string {
	return args["e"]
}
func wip(args map[string]string) string {
	return "Not implemented"
}
