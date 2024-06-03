# Image API

Base URL
http://dokalab.com

Endpoints
# 1. Upload Image
URL: /upload
Method: POST
Description: Uploads an image to the server.

Example Request:
bash
curl -F "picture=@path_to_your_image" http://dokalab.com/upload

Response:
Status Code: 200 OK
Body: Confirmation message with the filename.

# 2. Get Image
URL: /picture/{filename}
Method: GET
Description: Retrieves an image from the server by its filename.

Example Request:
bash
curl http://dokalab.com/picture

Response:
Status Code: 200 OK
Body: The image file.

# 3. Delete Image
URL: /picture/{filename}
Method: DELETE
Description: Deletes an image from the server by its filename.

Example Request:
bash
curl -X DELETE http://dokalab.com/delete_picture


Response:
Status Code: 200 OK
Body: Confirmation message indicating the file was deleted

File Storage
All uploaded images are stored in the ./uploads directory on the server under the name uploaded_image.jpg.

# Logs
Log files for the application and server are available for monitoring:

Application logs: /var/log/myapp.log
Nginx access logs: /var/log/nginx/access.log
Nginx error logs: /var/log/nginx/error.log
Example Log Viewing Commands
To view the logs, use the following commands:

bash
sudo tail -f /var/log/myapp.log
sudo tail -f /var/log/nginx/access.log
sudo tail -f /var/log/nginx/error.log

![Снимок экрана 2024-06-02 в 5 06 15 PM](https://github.com/Svetozara3363/ImageApi/assets/120119368/8c720ea8-ae6b-4d80-9718-d51ac0aefc1d)

![2024-06-02 17 21 29](https://github.com/Svetozara3363/ImageApi/assets/120119368/e6ea02b4-26d5-439a-a181-c6c8d831255f)

