package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

const uploadDir = "./uploads"

var db *sql.DB

func main() {
	log.Println("Starting application")

	logFile, err := os.OpenFile("./myapp.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	log.SetOutput(logFile)
	defer logFile.Close()

	log.Println("Checking upload directory")
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		err := os.Mkdir(uploadDir, os.ModePerm)
		if err != nil {
			log.Fatalf("Error creating upload directory: %v", err)
		}
	}

	log.Println("Connecting to PostgreSQL")
	connStr := "user=myuser password=mypassword dbname=mydb sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error connecting to PostgreSQL: %v", err)
	}
	defer db.Close()

	log.Println("Pinging PostgreSQL")
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error pinging PostgreSQL: %v", err)
	}
	log.Println("Successfully connected to PostgreSQL")

	router := mux.NewRouter()
	router.HandleFunc("/", HomeHandler)
	router.HandleFunc("/upload", UploadHandler).Methods("POST")
	router.HandleFunc("/picture/{id:[0-9]+}", GetPictureHandler).Methods("GET")
	router.HandleFunc("/pictures", GetAllPicturesHandler).Methods("GET")
	router.HandleFunc("/delete_picture/{id:[0-9]+}", DeletePictureHandler).Methods("DELETE")
	router.HandleFunc("/api/", APIRootHandler)
	router.HandleFunc("/api/upload", UploadHandler).Methods("POST")
	router.HandleFunc("/api/pictures", GetAllPicturesHandler).Methods("GET")

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

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	html := `
    <!DOCTYPE html>
    <html lang="en">
    <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <title>Image Upload</title>
    </head>
    <body>
        <h1>Upload an Image</h1>
        <form action="/upload" method="post" enctype="multipart/form-data">
            <input type="file" name="picture" accept="image/*" required>
            <button type="submit">Upload</button>
        </form>
    </body>
    </html>`
	w.Write([]byte(html))
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
	filePath := filepath.Join(uploadDir, handler.Filename)
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

	sessionID := uuid.New().String()

	var id int
	err = db.QueryRow("INSERT INTO images (filename, session_id) VALUES ($1, $2) RETURNING id", handler.Filename, sessionID).Scan(&id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error inserting file into database: %v", err)
		return
	}

	log.Printf("File uploaded successfully with ID: %d and Session ID: %s", id, sessionID)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message": "File uploaded successfully with ID: %d and Session ID: %s"}`, id, sessionID)
}

func GetPictureHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var filename string
	err := db.QueryRow("SELECT filename FROM images WHERE id = $1", id).Scan(&filename)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "File not found.", http.StatusNotFound)
			log.Printf("File not found for ID: %s", id)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("Error querying file from database: %v", err)
		}
		return
	}

	filePath := filepath.Join(uploadDir, filename)
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "File not found.", http.StatusNotFound)
		log.Printf("Error opening file: %v", err)
		return
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error reading file: %v", err)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.WriteHeader(http.StatusOK)
	w.Write(fileBytes)
}

func GetAllPicturesHandler(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session_id")
	if sessionID == "" {
		http.Error(w, "session_id is required", http.StatusBadRequest)
		return
	}

	rows, err := db.Query("SELECT id, filename FROM images WHERE session_id = $1", sessionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error querying images from database: %v", err)
		return
	}
	defer rows.Close()

	var pictures []struct {
		ID       int    `json:"id"`
		Filename string `json:"filename"`
	}

	for rows.Next() {
		var picture struct {
			ID       int    `json:"id"`
			Filename string `json:"filename"`
		}
		if err := rows.Scan(&picture.ID, &picture.Filename); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("Error scanning row: %v", err)
			return
		}
		pictures = append(pictures, picture)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error with rows: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(pictures); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error encoding response: %v", err)
		return
	}
}

func DeletePictureHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var filename string
	err := db.QueryRow("DELETE FROM images WHERE id = $1 RETURNING filename", id).Scan(&filename)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "File not found.", http.StatusNotFound)
			log.Printf("File not found for ID: %s", id)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("Error deleting file from database: %v", err)
		}
		return
	}

	filePath := filepath.Join(uploadDir, filename)
	err = os.Remove(filePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error deleting file: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "File deleted successfully: %s\n", filename)
}

func APIRootHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("API is working"))
}
