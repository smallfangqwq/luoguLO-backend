package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"strconv"
)

func GetData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // Get params
	var result DBDiscussTemplate
	val, err := strconv.Atoi(params["id"])
	if err != nil {
		result.Title = "404"
		json.NewEncoder(w).Encode(result)
		return
	}
	err = Session.DB("luogulo").C("discuss").Find(bson.M{"postid": val}).One(&result)
	if err != nil {
		result.Title = "404"
	}
	json.NewEncoder(w).Encode(result)
}

func webMain() {

	r := mux.NewRouter()
	r.HandleFunc("/api/discuss/data/{id}", GetData).Methods("GET")
	err := http.ListenAndServe(":"+strconv.Itoa(Config.Http.Port), r)
	if err != nil {
		panic(err)
	}
}
