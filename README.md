# Image Shrinking API

This is a simple API that allows you to upload multiple images and returns their URLs after shrinking them.

## Installation

1. Clone the repository
2. Install the dependencies by running `go mod download`
3. Run the application by running `go run main.go`

## Usage

### Endpoints

#### GET /api/ping

Returns a JSON response with a message "pong".

#### POST /api/images/shrink

Accepts multiple images and returns their URLs after shrinking them.

##### Request

- Method: POST
- Headers:
  - Content-Type: multipart/form-data
- Body:
  - images: multiple files

##### Response

- Status Code: 200 OK
- Body:
  - filepath: an array of URLs of the uploaded images after shrinking

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.