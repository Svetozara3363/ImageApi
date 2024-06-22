package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

const uploadDir = "./uploads"

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

	router := mux.NewRouter()
	router.HandleFunc("/", HomeHandler)
	router.HandleFunc("/upload", UploadHandler).Methods("POST")
	router.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir(uploadDir))))
	router.HandleFunc("/api/", APIRootHandler)
	router.HandleFunc("/api/upload", UploadHandler).Methods("POST")

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
        <form action="/api/upload" method="post" enctype="multipart/form-data">
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

	log.Printf("File uploaded successfully")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "File uploaded successfully", "path": "/uploads/` + handler.Filename + `"}`))
}

func APIRootHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("API is working"))
}
