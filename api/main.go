package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func uploadMultipleImages(c *gin.Context) {
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
		out, err := os.Create("./temp/images/" + newFilename)
		if err != nil {
			log.Fatal(err)
		}
		defer out.Close()

		renderFile, _ := file.Open()
		_, err = io.Copy(out, renderFile)
		if err != nil {
			log.Fatal(err)
		}
	}
	c.JSON(http.StatusOK, gin.H{"filepath": filePaths})
}

func main() {
	r := gin.Default()

	r.GET("/api/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.POST("/api/images/shrink", uploadMultipleImages)

	r.Run(":8800")
}
