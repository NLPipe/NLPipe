package main

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func main() {
	r := httprouter.New()

	log.Println("API listening on :8001")

	r.GET("/", homePage)
	r.NotFound = http.FileServer(http.Dir("html"))

	r.POST("/api/upload", upload)
	r.GET("/api/result/:uuid", result)
	log.Fatal(http.ListenAndServe(":8001", r))
}

func homePage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("content-type", "text/html")
	http.ServeFile(w, r, "./html/index.html")
}

func upload(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	uuidBytes, _ := ioutil.ReadFile("/proc/sys/kernel/random/uuid")
	uuid := strings.TrimRight(string(uuidBytes), "\n")
	fmt.Println("UUID:", uuid)

	body, err := ioutil.ReadAll(r.Body)
	if len(body) == 0 {
		log.Fatalf("Empty body. Cannot process %v", uuid)
		w.WriteHeader(400)
	}
	if err != nil {
		log.Fatalf("Error while reading file %v %v", uuid, err)
		w.WriteHeader(400)
	}

	result, err := uploadFile(uuid, body)

	if err != nil {
		log.Fatalf("Error during the upload of the file %v %v", uuid, err)
		w.WriteHeader(500)
	}

	fmt.Printf("Result for %v: %v", uuid, result)

	_, err = PutItem(uuid)

	if err != nil {
		w.WriteHeader(500)
	}

	http.Redirect(w, r, "/result.html?uuid="+uuid, http.StatusSeeOther)
}

func result(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	// Take the `uuid` request parameter
	uuid := p.ByName("uuid")
	if uuid == "" {
		w.WriteHeader(400)
		return
	}

	// Make the query
	item, err := GetResult(uuid)
	if err != nil {
		panic(fmt.Sprintf("Query error, %v", err))
	}

	// Check if the UUID doesn't exist in the DB
	if item.UUID == "" {
		w.WriteHeader(404)
		return
	}

	// Return result
	json, err := json.Marshal(item)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	_, err = w.Write(json)

	if err != nil {
		w.WriteHeader(500)
		return
	}
}
