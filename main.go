package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

const uploadDir = "./uploads"
const staticImageName = "uploaded_image.jpg"

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
	filePath := filepath.Join(uploadDir, staticImageName)
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

	imageUrl := "/uploads/" + staticImageName
	response := map[string]string{"imageUrl": imageUrl}

	log.Printf("File uploaded successfully, URL: %s", imageUrl)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func GetPictureHandler(w http.ResponseWriter, r *http.Request) {
	filePath := filepath.Join(uploadDir, staticImageName)
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
	w.Write(fileBytes)
}

func DeletePictureHandler(w http.ResponseWriter, r *http.Request) {
	filePath := filepath.Join(uploadDir, staticImageName)
	err := os.Remove(filePath)
	if err != nil {
		http.Error(w, "File not found.", http.StatusNotFound)
		log.Printf("Error deleting file: %v", err)
		return
	}

	response := map[string]string{"message": "File deleted successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
