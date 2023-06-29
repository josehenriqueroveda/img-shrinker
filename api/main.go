package main

import (
	"bytes"
	"fmt"
	"image"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func storeImages(c *gin.Context) {
	form, _ := c.MultipartForm()
	files := form.File["images"]
	filePaths := []string{}
	for _, file := range files {
		fileExt := filepath.Ext(file.Filename)
		originalFileName := strings.TrimSuffix(filepath.Base(file.Filename), filepath.Ext(file.Filename))
		now := time.Now()
		newFilename := strings.ReplaceAll(strings.ToLower(originalFileName), " ", "-") + "-" + fmt.Sprintf("%v", now.Unix()) + fileExt
		filePath := "http://localhost:8800/temp/bucket/" + newFilename

		filePaths = append(filePaths, filePath)
		out, err := os.Create("./temp/bucket/" + newFilename)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		defer out.Close()

		renderFile, _ := file.Open()
		_, err = io.Copy(out, renderFile)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	}
	c.JSON(http.StatusOK, gin.H{"filepaths": filePaths})
}

func shrinkImages(c *gin.Context) {
	form, _ := c.MultipartForm()
	files := form.File["images"]
	filePaths := []string{}
	for _, file := range files {
		fileExt := filepath.Ext(file.Filename)
		originalFileName := strings.TrimSuffix(filepath.Base(file.Filename), filepath.Ext(file.Filename))
		now := time.Now()
		newFilename := strings.ReplaceAll(strings.ToLower(originalFileName), " ", "-") + "-" + fmt.Sprintf("%v", now.Unix()) + fileExt
		filePath := "http://localhost:8800/temp/images/" + newFilename

		filePaths = append(filePaths, filePath)
		readerFile, _ := file.Open()
		imageFile, _, err := image.Decode(readerFile)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		if file.Size > 1024*1024 {
			for {
				src := imaging.Resize(imageFile, 720, 0, imaging.Lanczos)
				buf := new(bytes.Buffer)
				err = imaging.Encode(buf, src, imaging.JPEG)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				if int64(buf.Len()) < 1024*1024 {
					err = imaging.Save(src, fmt.Sprintf("./temp/images/%v", newFilename))
					if err != nil {
						c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
						return
					}
					break
				}
				imageFile = src
			}
		} else {
			err = imaging.Save(imageFile, fmt.Sprintf("./temp/images/%v", newFilename))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
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
