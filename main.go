package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
)

const uploadDir = "./uploads"

func main() {
	// Создаем директорию uploads, если она не существует
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.Mkdir(uploadDir, os.ModePerm)
	}

	// Настройка маршрутизатора
	router := mux.NewRouter()
	router.HandleFunc("/", HomeHandler)
	router.HandleFunc("/upload", UploadHandler).Methods("POST")
	router.HandleFunc("/picture/{filename}", GetPictureHandler).Methods("GET")
	router.HandleFunc("/picture/{filename}", DeletePictureHandler).Methods("DELETE")

	// Запуск сервера
	fmt.Println("Starting server at :8080")
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
	file, handler, err := r.FormFile("picture")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	filePath := filepath.Join(uploadDir, handler.Filename)
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
	fmt.Fprintf(w, "Файл успешно загружен: %s\n", handler.Filename)
}

// GetPictureHandler обрабатывает получение изображений
func GetPictureHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filename := vars["filename"]

	filePath := filepath.Join(uploadDir, filename)
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "Файл не найден.", http.StatusNotFound)
		return
	}
	defer file.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg") // Укажите соответствующий тип изображения
	w.WriteHeader(http.StatusOK)
	w.Write(fileBytes)
}

// DeletePictureHandler обрабатывает удаление изображений
func DeletePictureHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filename := vars["filename"]

	filePath := filepath.Join(uploadDir, filename)
	err := os.Remove(filePath)
	if err != nil {
		http.Error(w, "Файл не найден.", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Файл успешно удален: %s\n", filename)
}
