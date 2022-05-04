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
	w.Header().Set("Access-Control-Allow-Origin", Config.Http.AccessOrigin)
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

func UpdateData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", Config.Http.AccessOrigin)
	params := mux.Vars(r)
	val, err := strconv.Atoi(params["id"])
	if err != nil {
		json.NewEncoder(w).Encode(bson.M{"status": "error! No string !"})
		return
	}
	fromPage, err := strconv.Atoi(params["fromPage"])
	if err != nil {
		json.NewEncoder(w).Encode(bson.M{"status": "error! No string!"})
		return
	}
	endPage, err := strconv.Atoi(params["endPage"])
	if err != nil {
		json.NewEncoder(w).Encode(bson.M{"status": "error! No string! "})
		return
	}
	SaveAlone(Session, val, fromPage, endPage)
	json.NewEncoder(w).Encode(bson.M{"status": "ok."})
}

func webMain() {
	r := mux.NewRouter()
	r.HandleFunc("/api/discuss/data/{id}", GetData).Methods("GET")
	r.HandleFunc("/api/discuss/update/{id}/{fromPage}:{endPage}", UpdateData).Methods("GET")
	err := http.ListenAndServe(":"+strconv.Itoa(Config.Http.Port), r)
	if err != nil {
		panic(err)
	}
}
