package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/mediocregopher/radix.v2/redis"
)

type RpushRequest struct {
	Queue string `json:"queue"`
	Data  string `json:"data"`
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/LRANGE/{queue}/{indexStart}/{indexEnd}", Lrange)
	router.HandleFunc("/RPUSH", Rpush)

	http.ListenAndServe(":8080", router)
}

func Lrange(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	queue := vars["queue"]
	starting, _ := strconv.Atoi(vars["indexStart"])
	ending, _ := strconv.Atoi(vars["indexEnd"])

    client, err := redis.Dial("tcp", os.Getenv("REDIS_HOST")+":"+os.Getenv("REDIS_PORT"))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	res, err := client.Cmd("LRANGE", queue, starting, ending).List()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
    
	jsonResponse, _ := json.Marshal(map[string][]string{"list_values": res})

    w.Write(jsonResponse)

}

func Rpush(w http.ResponseWriter, r *http.Request) {
	dec, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	var body RpushRequest
	err = json.Unmarshal(dec, &body)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

    client, err := redis.Dial("tcp", os.Getenv("REDIS_HOST")+":"+os.Getenv("REDIS_PORT"))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	res := client.Cmd("RPUSH", body.Queue, body.Data).Err
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	log.Println(res)
}

