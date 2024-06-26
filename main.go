package main

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

const (
	dbConnString = "user=myuser password=yourpassword dbname=mydb sslmode=disable"
)

var db *sql.DB

func main() {
	log.Println("Starting application")

	logFile, err := os.OpenFile("./myapp.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	log.SetOutput(logFile)
	defer logFile.Close()

	db, err = sql.Open("postgres", dbConnString)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	router := mux.NewRouter()
	router.HandleFunc("/api/upload", UploadHandler).Methods("POST")
	router.HandleFunc("/api/pictures", GetPictureHandler).Methods("GET")
	router.HandleFunc("/api/pictures", DeletePictureHandler).Methods("DELETE")

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"POST", "GET", "DELETE"},
		AllowedHeaders: []string{"Content-Type"},
	})

	handler := c.Handler(router)

	log.Println("Starting server at :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received upload request")
	file, handler, err := r.FormFile("picture")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("Error getting file from form: %v", err)
		return
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error reading file: %v", err)
		return
	}

	encodedImage := base64.StdEncoding.EncodeToString(fileBytes)

	_, err = db.Exec("INSERT INTO pictures (name, data) VALUES ($1, $2)", handler.Filename, encodedImage)
	if err != nil {
		http.Error(w, "Error saving image to database.", http.StatusInternalServerError)
		log.Printf("Error saving image to database: %v", err)
		return
	}

	log.Println("File uploaded successfully")
	response := map[string]string{"message": "File uploaded successfully"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func GetPictureHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request for picture")
	var name, data string

	err := db.QueryRow("SELECT name, data FROM pictures ORDER BY name DESC LIMIT 1").Scan(&name, &data)
	if err != nil {
		http.Error(w, "File not found.", http.StatusNotFound)
		log.Printf("Error fetching image from database: %v", err)
		return
	}

	imageData, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error decoding image data: %v", err)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Write(imageData)
}

func DeletePictureHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to delete picture")
	var name string
	err := db.QueryRow("SELECT name FROM pictures ORDER BY name DESC LIMIT 1").Scan(&name)
	if err != nil {
		http.Error(w, "File not found.", http.StatusNotFound)
		log.Printf("Error fetching image name from database: %v", err)
		return
	}

	_, err = db.Exec("DELETE FROM pictures WHERE name = $1", name)
	if err != nil {
		http.Error(w, "Error deleting image from database.", http.StatusInternalServerError)
		log.Printf("Error deleting image from database: %v", err)
		return
	}

	response := map[string]string{"message": "File deleted successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
