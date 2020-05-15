package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"local.org/complexlog"
)

const (
	DOMAIN = "api.test.local"
)

var (
	COMMAND = map[string]func(map[string]string) string{
		"help":    help,
		"list":    list,
		"read":    read,
		"comment": comment,
		"login":   login,
		"echo":    echo,
	}
	mdb MongoDb
)

func main() {

	mdb = MongoDb{}
	mdb.connect()

	//api router setting
	r := mux.NewRouter()
	r.Host(DOMAIN)

	r.HandleFunc("/cmd/{cmd}", withCORS(sendMethod)).Methods("OPTIONS")
	r.HandleFunc("/cmd/{cmd}", withCORS(runCmdApi)).Methods("POST")
	r.HandleFunc("/echo/{msg}", echoApi).Methods("GET")

	complexlog.Servlog("api init")

	//server web lintening
	err := http.ListenAndServe(":8000", r)
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
	return fmt.Sprintf(helpText, "User ")
}
func list(args map[string]string) string {
	mdb.getPost()
	return "read"
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
