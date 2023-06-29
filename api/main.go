package main

import (
	"bytes"
	"fmt"
	"image"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/disintegration/imaging"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func storeImages(c *gin.Context) {
	form, _ := c.MultipartForm()
	files := form.File["images"]
	filePaths := make([]string, len(files))
	errChan := make(chan error, len(files))

	var wg sync.WaitGroup
	wg.Add(len(files))

	for i, file := range files {
		go func(i int, file *multipart.FileHeader) {
			defer wg.Done()

			fileExt := filepath.Ext(file.Filename)
			originalFileName := strings.TrimSuffix(filepath.Base(file.Filename), filepath.Ext(file.Filename))
			now := time.Now()
			newFilename := strings.ReplaceAll(strings.ToLower(originalFileName), " ", "-") + "-" + fmt.Sprintf("%v", now.Unix()) + fileExt
			filePath := "http://localhost:8800/temp/bucket/" + newFilename

			filePaths[i] = filePath

			out, err := os.Create("./temp/bucket/" + newFilename)
			if err != nil {
				errChan <- err
				return
			}
			defer out.Close()

			renderFile, err := file.Open()
			if err != nil {
				errChan <- err
				return
			}
			defer renderFile.Close()

			_, err = io.Copy(out, renderFile)
			if err != nil {
				errChan <- err
				return
			}
		}(i, file)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"filepaths": filePaths})
}

func shrinkImages(c *gin.Context) {
	form, _ := c.MultipartForm()
	files := form.File["images"]
	filePaths := make([]string, len(files))
	errChan := make(chan error, len(files))

	var wg sync.WaitGroup
	wg.Add(len(files))

	for i, file := range files {
		go func(i int, file *multipart.FileHeader) {
			defer wg.Done()

			fileExt := filepath.Ext(file.Filename)
			originalFileName := strings.TrimSuffix(filepath.Base(file.Filename), filepath.Ext(file.Filename))
			now := time.Now()
			newFilename := strings.ReplaceAll(strings.ToLower(originalFileName), " ", "-") + "-" + fmt.Sprintf("%v", now.Unix()) + fileExt
			filePath := "http://localhost:8800/temp/images/" + newFilename

			filePaths[i] = filePath

			readerFile, err := file.Open()
			if err != nil {
				errChan <- err
				return
			}
			defer readerFile.Close()

			imageFile, _, err := image.Decode(readerFile)
			if err != nil {
				errChan <- err
				return
			}

			if file.Size > 1024*1024 {
				for {
					src := imaging.Resize(imageFile, 720, 0, imaging.Lanczos)
					buf := new(bytes.Buffer)
					err = imaging.Encode(buf, src, imaging.JPEG)
					if err != nil {
						errChan <- err
						return
					}
					if int64(buf.Len()) < 1024*1024 {
						err = imaging.Save(src, fmt.Sprintf("./temp/images/%v", newFilename))
						if err != nil {
							errChan <- err
							return
						}
						break
					}
					imageFile = src
				}
			} else {
				err = imaging.Save(imageFile, fmt.Sprintf("./temp/images/%v", newFilename))
				if err != nil {
					errChan <- err
					return
				}
			}
		}(i, file)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"filepaths": filePaths})
}

func main() {
	r := gin.Default()

	r.LoadHTMLGlob("./templates/*")

	r.Static("/temp", "./temp")

	// Enable CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowHeaders = []string{"Origin", "Content-Lenght", "Content-Type", "Authorization"}
	r.Use(cors.New(config))

	r.GET("/api/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.POST("/api/images/store", storeImages)

	r.POST("/api/images/shrink", shrinkImages)

	r.Run(":8800")

}
