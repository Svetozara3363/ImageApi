## API Endpoints

### Upload Image
- **URL:** curl -X POST https://dokalab.com/api/upload -F "file=@path/to/your/image.jpg"


- **Method:** `POST`
- **Response:**
  - `200 OK` on successful upload
  - `400 Bad Request` in case of error

### Get Image
- **URL:** curl -X GET https://dokalab.com/api/pictures --output image.jpg

- **Method:** `GET`
- **Response:**
  - `200 OK` with image data in binary format
  - `404 Not Found` if the image is not found

### Delete Image
- **URL:** curl -X DELETE https://dokalab.com/api/delete.  
- **Method:** `DELETE`
- **Response:**
  - `200 OK` on successful deletion
  - `500 Internal Server Error` in case of error
