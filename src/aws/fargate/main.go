package main

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var cfg *Config

type Config struct {
	Region           string
	DynamoDbEndpoint string
	DynamoDbTable    string
	S3Endpoint       string
	S3Bucket         string
	S3ForcePathStyle bool
}

func loadConfigFromEnv() *Config {
	fps, err := strconv.ParseBool(os.Getenv("S3_FORCE_PATH_STYLE"))
	if err != nil {
		fps = false
	}

	cfg := &Config{
		Region:           os.Getenv("REGION"),
		DynamoDbEndpoint: os.Getenv("DYNAMODB_ENDPOINT"),
		DynamoDbTable:    os.Getenv("DYNAMODB_TABLE"),
		S3Endpoint:       os.Getenv("S3_ENDPOINT"),
		S3Bucket:         os.Getenv("S3_BUCKET"),
		S3ForcePathStyle: fps,
	}

	return cfg
}

func (cfg Config) Dump() string {
	ret, err := json.Marshal(cfg)
	if err != nil {
		panic(fmt.Sprintf("Failed to marshal the environment variables %v", cfg))
	}
	return string(ret)
}

func main() {
	cfg = loadConfigFromEnv()
	log.Printf("Configuration from environment variables: %v\n", cfg.Dump())

	r := httprouter.New()

	r.GET("/", homePage)
	r.NotFound = http.FileServer(http.Dir("html"))

	r.POST("/api/upload", upload)
	r.GET("/api/result/:uuid", result)

	log.Println("API listening on :8001")
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
		w.WriteHeader(400)
		log.Fatalf("Empty body. Cannot process %v", uuid)
	}
	if err != nil {
		w.WriteHeader(400)
		log.Fatalf("Error while reading file %v %v", uuid, err)
	}

	result, err := uploadFile(uuid, body)

	if err != nil {
		w.WriteHeader(500)
		log.Fatalf("Error during the upload of the file %v %v", uuid, err)
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
