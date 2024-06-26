package main

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "myuser"
	password = "Sveta.2003"
	dbname   = "mydb"
)

var db *sql.DB

func initDB() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Error opening database: %q", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error connecting to the database: %q", err)
	}
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fileBase64 := base64.StdEncoding.EncodeToString(fileBytes)
	_, err = db.Exec("INSERT INTO pictures (name, data) VALUES ($1, $2)", "uploaded_image", fileBase64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func getPictureHandler(w http.ResponseWriter, r *http.Request) {
	var data string
	err := db.QueryRow("SELECT data FROM pictures WHERE name = $1", "uploaded_image").Scan(&data)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "No image found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	fileBytes, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Write(fileBytes)
}

func deletePictureHandler(w http.ResponseWriter, r *http.Request) {
	_, err := db.Exec("DELETE FROM pictures WHERE name = $1", "uploaded_image")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func main() {
	initDB()

	router := mux.NewRouter()
	router.HandleFunc("/upload", uploadHandler).Methods("POST")
	router.HandleFunc("/picture", getPictureHandler).Methods("GET")
	router.HandleFunc("/picture", deletePictureHandler).Methods("DELETE")

	fmt.Println("Starting server at :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
