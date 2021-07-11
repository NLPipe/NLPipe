package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"mime/multipart"
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
		log.Panicf("Failed to marshal the environment variables %v", cfg)
	}
	return string(ret)
}

func main() {
	cfg = loadConfigFromEnv()
	log.Infof("Configuration from environment variables: %v\n", cfg.Dump())

	r := httprouter.New()

	r.GET("/", homePage)
	r.NotFound = http.FileServer(http.Dir("html"))

	r.POST("/api/upload", upload)
	r.GET("/api/result/:uuid", result)

	log.Info("API listening on :8001")
	log.Fatal(http.ListenAndServe(":8001", r))
}

func homePage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("content-type", "text/html")
	http.ServeFile(w, r, "./html/index.html")
}

func upload(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	uuidBytes, _ := ioutil.ReadFile("/proc/sys/kernel/random/uuid")
	uuid := strings.TrimRight(string(uuidBytes), "\n")

	// Max 10 MB
	err := r.ParseMultipartForm(10 << 20)
	file, handler, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Warningf("Error while reading the file: %v", err)
		return
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Errorf("Error while closing file. %v", err)
			return
		}
	}(file)

	if strings.Split(handler.Header.Get("Content-Type"), "/")[0] != "audio" {
		w.WriteHeader(http.StatusBadRequest)
		log.Warningf("Got a non audio file. Aborting.")
		return
	}

	if handler.Size == 0 {
		w.WriteHeader(http.StatusBadRequest)
		log.Warningf("Got a 0-size file. Aborting.")
		return
	}

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Warningf("Error while reading file bytes for %v: %v", uuid, err)
		return
	}

	_, err = uploadFile(uuid, fileBytes)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Errorf("Error during the upload of the file %v: %v", uuid, err)
		return
	}

	_, err = PutItem(uuid)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Errorf("Error while putting on DynamoDB %v: %v", uuid, err)
		return
	}

	http.Redirect(w, r, "/result.html?uuid="+uuid, http.StatusSeeOther)
}

func result(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	// Take the `uuid` request parameter
	uuid := p.ByName("uuid")
	if uuid == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Make the query
	item, err := GetResult(uuid)
	if err != nil {
		log.Errorf("Result query error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Check if the UUID doesn't exist in the DB
	if item.UUID == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Return result
	json, err := json.Marshal(item)
	if err != nil {
		w.WriteHeader(500)
		log.Errorf("Failed to marshal the result: %v", err)
		return
	}

	_, err = w.Write(json)

	if err != nil {
		w.WriteHeader(500)
		log.Errorf("Failed to write the JSON: %v", err)
		return
	}
}
