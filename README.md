## API Endpoints

### Upload Image
- **URL:** [https://dokalab.com/upload](https://dokalab.com/api/upload)
- **Method:** `POST`
- **Response:**
  - `200 OK` on successful upload
  - `400 Bad Request` in case of error

### Get Image
- **URL:** [https://dokalab.com/pictures](https://dokalab.com/api/pictures)
- **Method:** `GET`
- **Response:**
  - `200 OK` with image data in binary format
  - `404 Not Found` if the image is not found

### Delete Image
- **URL:** [https://dokalab.com/delete](https://dokalab.com/api/delete)
- **Method:** `DELETE`
- **Response:**
  - `200 OK` on successful deletion
  - `500 Internal Server Error` in case of error
