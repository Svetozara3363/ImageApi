# Image API
Этот проект  API для загрузки, получения и удаления изображений с использованием Go и Gorilla Mux.
## Запуск сервера
go run main.go

1 . Откройте браузер и перейдите по адресу http://localhost:8080/.
## Метод: GET
Описание: Возвращает HTML-форму для загрузки изображения.
Пример запроса:
Загрузка изображения
URL: /upload
## Метод: POST
Описание: Загружает изображение на сервер.
Параметры: Принимает файл изображения с ключом picture в форме данных.
Пример запроса: curl -X POST -F 'picture=@/path/to/your/image.png' http://localhost:8080/upload

Пример ответа:
arduino
Файл успешно загружен: image.png
Получение изображения
URL: /picture/{filename}
## Метод: GET
Описание: Возвращает изображение с сервера.
Параметры:
filename - имя файла изображения.
Пример запроса:
curl http://localhost:8080/picture/image.png --output downloaded_image.png
Пример ответа: Файл изображения.
Удаление изображения
URL: /picture/{filename}
## Метод: DELETE
Описание: Удаляет изображение с сервера.
Параметры:
filename - имя файла изображения.
curl -X DELETE http://localhost:8080/picture/image.png
Пример ответа:
Файл успешно удален: image.png
<img width="634" alt="Снимок экрана 2024-05-30 в 2 31 02 PM" src="https://github.com/Svetozara3363/ImageApi/assets/120119368/a474c8cc-02a5-44e2-afe4-0205b81b1590">
