package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
)

const uploadDir = "./uploads"
const staticImageName = "uploaded_image.jpg"

func main() {
	// Создаем или открываем лог файл
	logFile, err := os.OpenFile("/var/log/myapp.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(logFile)
	defer logFile.Close()

	// Создаем директорию uploads, если её нет
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.Mkdir(uploadDir, os.ModePerm)
	}

	// Настройка маршрутизатора
	router := mux.NewRouter()
	router.HandleFunc("/", HomeHandler)
	router.HandleFunc("/upload", UploadHandler).Methods("POST")
	router.HandleFunc("/picture", GetPictureHandler).Methods("GET")
	router.HandleFunc("/delete_picture", DeletePictureHandler).Methods("DELETE")

	// Запуск сервера
	fmt.Println("Starting server at :8080")
	log.Println("Starting server at :8080")
	http.ListenAndServe(":8080", router)
}

// HomeHandler обслуживает HTML-форму
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
            <input type="file" name="picture" accept="image/*">
            <button type="submit">Upload</button>
        </form>
    </body>
    </html>`
	w.Write([]byte(html))
}

// UploadHandler обрабатывает загрузку изображений
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("picture")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	filePath := filepath.Join(uploadDir, staticImageName)
	f, err := os.Create(filePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()

	_, err = io.Copy(f, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "File uploaded successfully: %s\n", staticImageName)
}

// GetPictureHandler обрабатывает получение изображения
func GetPictureHandler(w http.ResponseWriter, r *http.Request) {
	filePath := filepath.Join(uploadDir, staticImageName)
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "File not found.", http.StatusNotFound)
		return
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.WriteHeader(http.StatusOK)
	w.Write(fileBytes)
}

// DeletePictureHandler обрабатывает удаление изображения
func DeletePictureHandler(w http.ResponseWriter, r *http.Request) {
	filePath := filepath.Join(uploadDir, staticImageName)
	err := os.Remove(filePath)
	if err != nil {
		http.Error(w, "File not found.", http.StatusOK) // Change here to avoid 404
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "File deleted successfully: %s\n", staticImageName)
}
