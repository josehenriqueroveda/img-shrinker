# Image Shrinker API
[![made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](https://go.dev/)
[![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/josehenriqueroveda/img-shrinker.svg)](https://github.com/josehenriqueroveda/img-shrinker)

This is an API written in Go that allows to resize and store images. This is a common use case for image hosting websites, where users upload images and the website reduce its size, stores them and returns a URL to the image. It can be used as a microservice in a larger application and reduce costs by reducing the size of the images stored.

## Installation

1. Clone the repository:
```bash
git clone https://github.com/josehenriqueroveda/img-shrinker.git
```
2. Install the dependencies by running:
```bash
go mod download
```

3. Run the application by running:
```bash
go run main.go
```

## Usage

The application runs on port 8800 by default. You can change it on `r.Run(":8800")` on `main.go`.

### Endpoints

#### API Health Check
```http
GET /api/ping
```

Returns a JSON response if the service is up and running.
```json
{
  "message": "pong"
}
```

#### Image Storage
##### Request

- Method: POST
- Headers:
  - Content-Type: multipart/form-data
- Body:
  - images: multiple files

```http
POST /api/images/store
```

Accepts multiple images and returns their URLs after storing them on its original size.
##### Response

- Status Code: 200 OK
- Body:
  - filepath: an array of URLs of the uploaded images after shrinking

```json
{
  "filepath": [
    "http://localhost:8800/images/2021/08/01/1627820000_1.jpg",
    "http://localhost:8800/images/2021/08/01/1627820000_2.jpg"
  ]
}
```

#### Image Shrink
##### Request

- Method: POST
- Headers:
  - Content-Type: multipart/form-data
- Body:
  - images: multiple files

```http
POST /api/images/shrink
```

Accepts multiple images and returns their URLs after shrinking them.
##### Response

- Status Code: 200 OK
- Body:
  - filepath: an array of URLs of the uploaded images after shrinking

```json
{
  "filepath": [
    "http://localhost:8800/images/2021/08/01/1627820000_1.jpg",
    "http://localhost:8800/images/2021/08/01/1627820000_2.jpg"
  ]
}
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing
If you find a bug or have a feature request, please open an issue on the repository. If you would like to contribute code, please fork the repository and submit a pull request.
