# Image API

Base URL
http://dokalab.com

Endpoints
# 1. Upload Image
URL:/api/upload
Method: POST
Description: Uploads an image to the server.

Example Request:
bash
curl -F "picture=@path_to_your_image" https://dokalab.com/api/upload


Response:
Status Code: 200 OK
Body: Confirmation message with the filename and session ID in JSON format.

# 2. Get Image
URL: /api/picture/{id}
Method: GET
Description: Retrieves an image from the server by its ID.

Example Request:
bash
curl https://dokalab.com/api/picture/{id}


Response:
Status Code: 200 OK
Body: The image file.

# 3. Get All Images
URL: /api/pictures?session_id={session_id}
Method: GET
Description: Retrieves all images associated with a session ID.

Example Request:
bash
curl https://dokalab.com/api/pictures?session_id={session_id}

Response:
Status Code: 200 OK
Body: JSON array of image IDs and filenames.

# 4. Delete Image
URL: /api/delete_picture/{id}
Method: DELETE
Description: Deletes an image from the server by its ID.

Example Request:
bash
curl -X DELETE https://dokalab.com/api/delete_picture/{id}


Response:
Status Code: 200 OK
Body: Confirmation message indicating the file was deleted.

# File Storage

All uploaded images are stored in the ./uploads directory on the server.

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
<img width="893" alt="Снимок экрана 2024-06-18 в 12 25 25 PM" src="https://github.com/Svetozara3363/ImageApi/assets/120119368/53171572-e1bd-4b75-80e0-edd00b6b34f8">
<img width="889" alt="Снимок экрана 2024-06-18 в 12 25 10 PM" src="https://github.com/Svetozara3363/ImageApi/assets/120119368/1d8a58ba-11e6-4308-82aa-8c979173cc8d">
<img width="875" alt="Снимок экрана 2024-06-18 в 12 25 17 PM" src="https://github.com/Svetozara3363/ImageApi/assets/120119368/e5bc1f98-4198-4ab0-93c7-b08f30332e2e">


