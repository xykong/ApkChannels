package services

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/xykong/ApkChannels/sign"
	"log"
	"net/http"
	"path/filepath"
)

type HTTPGenericResponse struct {
	Code    int
	Message string
}

func Start() {

	addr := viper.GetString("address")

	r := gin.Default()

	r.GET("/ping", ping())
	r.GET("/download/*apk", download())

	_ = r.Run(addr) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func download() func(c *gin.Context) {

	root := viper.GetString("root")
	root, err := filepath.Abs(root)
	if err != nil {
		log.Fatalf("root path Abs failed: %v", err)
	}
	fmt.Println("Absolute root path:", root)

	return func(c *gin.Context) {

		apk := c.Param("apk")
		fullPath := filepath.Join(root, apk)

		v1 := c.Query("v1") // shortcut for c.Request.URL.Query().Get("v1")

		log.Printf("request apk file %s with channel: %s", fullPath, v1)

		srcFile, offset, err := sign.CreateSrcReader(fullPath)
		if err != nil {
			c.JSON(http.StatusNotFound, HTTPGenericResponse{
				Code:    http.StatusNotFound,
				Message: fmt.Sprintf("open file failed: %v", fullPath),
			})
			return
		}
		defer srcFile.Close()

		c.Writer.Header().Add("Content-type", "application/octet-stream")
		err = sign.V2WriteStream(c.Writer, srcFile, offset, v1)

		if err != nil {
			c.JSON(http.StatusServiceUnavailable, HTTPGenericResponse{
				Code:    http.StatusServiceUnavailable,
				Message: "v1 sign add channel failed: " + err.Error(),
			})
			return
		}
	}
}

func ping() func(c *gin.Context) {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	}
}
