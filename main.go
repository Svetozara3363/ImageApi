package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
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
	file, header, err := r.FormFile("file")
	if err != nil {
		log.Printf("Error getting form file: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("Error reading file: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	imageName := "uploaded_image"

	_, err = db.Exec("DELETE FROM pictures WHERE name = $1", imageName)
	if err != nil {
		log.Printf("Error deleting old image: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = db.Exec("INSERT INTO pictures (name, data) VALUES ($1, $2)", imageName, fileBytes)
	if err != nil {
		log.Printf("Error inserting into database: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	savePath := filepath.Join("uploads", header.Filename)
	err = ioutil.WriteFile(savePath, fileBytes, 0644)
	if err != nil {
		log.Printf("Error saving file to server: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully uploaded image")
	w.WriteHeader(http.StatusOK)
}

func getPictureHandler(w http.ResponseWriter, r *http.Request) {
	var data []byte
	err := db.QueryRow("SELECT data FROM pictures WHERE name = $1", "uploaded_image").Scan(&data)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No image found in database")
			http.Error(w, "No image found", http.StatusNotFound)
		} else {
			log.Printf("Error querying database: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	log.Printf("Successfully retrieved image")
	w.Header().Set("Content-Type", "image/jpeg")
	w.Write(data)
}

func deletePictureHandler(w http.ResponseWriter, r *http.Request) {
	_, err := db.Exec("DELETE FROM pictures WHERE name = $1", "uploaded_image")
	if err != nil {
		log.Printf("Error deleting from database: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	savePath := filepath.Join("uploads", "uploaded_image.jpg")
	err = os.Remove(savePath)
	if err != nil {
		log.Printf("Error deleting file from server: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully deleted image")
	w.WriteHeader(http.StatusOK)
}

func main() {
	initDB()
	os.Mkdir("uploads", os.ModePerm)

	router := mux.NewRouter()
	router.HandleFunc("/api/upload", uploadHandler).Methods("POST")
	router.HandleFunc("/api/pictures", getPictureHandler).Methods("GET")
	router.HandleFunc("/api/delete", deletePictureHandler).Methods("DELETE")

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"https://dokalab.com"},
		AllowedMethods: []string{"GET", "POST", "DELETE"},
		AllowedHeaders: []string{"Content-Type"},
	})

	handler := c.Handler(router)
	fmt.Println("Starting server at :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
