package main

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

const (
	uploadDir  = "./uploads"
	dbUser     = "myuser"
	dbPassword = "mypassword"
	dbName     = "mydb"
	dbHost     = "localhost"
	dbPort     = "5432"
)

var db *sql.DB

func main() {
	log.Println("Starting application")

	var err error
	psqlInfo := os.Getenv("DATABASE_URL")
	if psqlInfo == "" {
		psqlInfo = "host=" + dbHost + " port=" + dbPort + " user=" + dbUser + " password=" + dbPassword + " dbname=" + dbName + " sslmode=disable"
	}

	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/api/upload", UploadHandler).Methods("POST")
	router.HandleFunc("/api/pictures", GetPictureHandler).Methods("GET")
	router.HandleFunc("/api/pictures", DeletePictureHandler).Methods("DELETE")
	router.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir(uploadDir))))

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"POST", "GET", "DELETE"},
		AllowedHeaders: []string{"Content-Type"},
	})

	handler := c.Handler(router)

	log.Println("Checking upload directory")
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		err := os.Mkdir(uploadDir, os.ModePerm)
		if err != nil {
			log.Fatalf("Error creating upload directory: %v", err)
		}
	}

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

	log.Printf("Uploading file: %s", handler.Filename)
	filePath := filepath.Join(uploadDir, "uploaded_image.jpg")
	f, err := os.Create(filePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error creating file: %v", err)
		return
	}
	defer f.Close()

	_, err = io.Copy(f, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error copying file: %v", err)
		return
	}

	// Save to the database
	imageData, err := os.ReadFile(filePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error reading file: %v", err)
		return
	}

	_, err = db.Exec("DELETE FROM pictures")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error deleting previous image from database: %v", err)
		return
	}

	_, err = db.Exec("INSERT INTO pictures (name, data) VALUES ($1, $2)", handler.Filename, imageData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error inserting image into database: %v", err)
		return
	}

	imageUrl := "/uploads/uploaded_image.jpg"
	response := map[string]string{"imageUrl": imageUrl}

	log.Printf("File uploaded successfully, URL: %s", imageUrl)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func GetPictureHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request for picture")
	var name string
	var data []byte

	err := db.QueryRow("SELECT name, data FROM pictures ORDER BY name DESC LIMIT 1").Scan(&name, &data)
	if err != nil {
		http.Error(w, "File not found.", http.StatusNotFound)
		log.Printf("Error fetching image from database: %v", err)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Write(data)
}

func DeletePictureHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to delete picture")
	_, err := db.Exec("DELETE FROM pictures")
	if err != nil {
		http.Error(w, "Error deleting picture.", http.StatusInternalServerError)
		log.Printf("Error deleting image from database: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Image deleted successfully"))
}
