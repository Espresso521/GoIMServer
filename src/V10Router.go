package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type V10Router struct {
	router *gin.RouterGroup
}

func (r *V10Router) addV10Router() {
  r.router.GET("/route", func(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"message": "Custom route"})
  })

	// 设置一个get请求的路由，url为/ping, 处理函数（或者叫控制器函数）是一个闭包函数。
	r.router.GET("/ping", func(c *gin.Context) {
			// 通过请求上下文对象Context, 直接往客户端返回一个json
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.router.POST("/app/start/", func(c *gin.Context) {
		c.Request.ParseMultipartForm(1024*1024)
		for key, values := range c.Request.PostForm {
			log.Printf("Form field %q, Values %q\n", key, values)
	  }
	  c.JSON(http.StatusOK, getAnswer())
  })

	r.router.POST("/run/start/", func(c *gin.Context) {
		c.Request.ParseMultipartForm(1024*1024)
	  c.JSON(http.StatusOK, ReceiveFile(c.Writer, c.Request))
  })

	r.router.POST("/run/finish/", func(c *gin.Context) {
		c.Request.ParseMultipartForm(1024*1024)
	  c.JSON(http.StatusOK, ReceiveFile(c.Writer, c.Request))
  })

	r.router.POST("/run/status/", func(c *gin.Context) {
		c.Request.ParseMultipartForm(1024*1024)
		for key, values := range c.Request.PostForm {
			log.Printf("Form field %q, Values %q\n", key, values)
	  }
	  c.JSON(http.StatusOK, getAnswer())
  })
}

func ReceiveFile(w http.ResponseWriter, r *http.Request) map[string]any {
	file, header, err := r.FormFile("image")
	if err != nil {
			panic(err)
	}
	defer file.Close()
	fmt.Printf("File name %s\n", header.Filename)

	for key, values := range r.Form {
		log.Printf("Form field %q, Values %q\n", key, values)
 }

	// create a new file in the local file system
	f, err := os.OpenFile("./uploads/"+header.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return gin.H{
			"auto_token":"a48396e4f5bec65ddd415cb802cd37be7a5784cae",
			"time":time.Now(),
			"error":fmt.Sprintf("File uploaded failed: Error is %v ", err.Error()),
		}
	}
	defer f.Close()

	// copy the uploaded file data to the new file
	_, err = io.Copy(f, file)
	if err != nil {
		return gin.H{
			"auto_token":"a48396e4f5bec65ddd415cb802cd37be7a5784cae",
			"time":time.Now(),
			"error":fmt.Sprintf("File uploaded failed: Error is %v ", err.Error()),
		}
	}

	fStat, _ := f.Stat()

	return gin.H{
		"auto_token":"a48396e4f5bec65ddd415cb802cd37be7a5784cae",
		"time":time.Now(),
		"extra":fmt.Sprintf("File uploaded successfully: file name is %v, file size is %v bytes", f.Name(), fStat.Size()),
	}
}