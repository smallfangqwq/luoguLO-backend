package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

func GetData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", Config.Http.AccessOrigin)
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("content-type", "application/json;charset=UTF-8") 
	params := mux.Vars(r) // Get params
	var result DBDiscussTemplate
	val, err := strconv.Atoi(params["id"])
	if err != nil {
		result.Title = "404"
		err := json.NewEncoder(w).Encode(result)
		if err != nil {
			return
		}
		return
	}
	err = Session.DB("luogulo").C("discuss").Find(bson.M{"postid": val}).One(&result)
	if err != nil {
		result.Title = "404"
	}
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		return
	}
}

func UpdateData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", Config.Http.AccessOrigin)
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("content-type", "application/json;charset=UTF-8")
	params := mux.Vars(r)
	val, err := strconv.Atoi(params["id"])
	if err != nil {
		err := json.NewEncoder(w).Encode(bson.M{"status": "error! No string !"})
		if err != nil {
			return
		}
		return
	}
	fromPage, err := strconv.Atoi(params["fromPage"])
	if err != nil {
		err := json.NewEncoder(w).Encode(bson.M{"status": "error! No string!"})
		if err != nil {
			return
		}
		return
	}
	endPage, err := strconv.Atoi(params["endPage"])
	if err != nil {
		err := json.NewEncoder(w).Encode(bson.M{"status": "error! No string! "})
		if err != nil {
			return
		}
		return
	}
	SaveAlone(Session, val, fromPage, endPage)
	err = json.NewEncoder(w).Encode(bson.M{"status": "ok."})
	if err != nil {
		return
	}
}

func UpdateDataAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", Config.Http.AccessOrigin)
	params := mux.Vars(r)
	val, err := strconv.Atoi(params["id"])
	if err != nil {
		//	err := json.NewEncoder(w).Encode(bson.M{"status": "error! No string! "})
		if err != nil {
			return
		}
		return
	}
	SaveNewDiscuss(Session, Config, val)
	//	err = json.NewEncoder(w).Encode(bson.M{"status": "ok."})
	if err != nil {
		return
	}
}

func webMain() {
	r := mux.NewRouter()
	r.HandleFunc("/api/discuss/data/{id}", GetData).Methods("GET")
	r.HandleFunc("/api/discuss/update/{id}", UpdateDataAll).Methods("GET")
	r.HandleFunc("/api/discuss/updatealone/{id}/{fromPage}:{endPage}", UpdateData).Methods("GET")
	err := http.ListenAndServe(":"+strconv.Itoa(Config.Http.Port), r)
	if err != nil {
		panic(err)
	}
}
