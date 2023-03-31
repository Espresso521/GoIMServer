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

type v20router struct {
	router *gin.RouterGroup
}

func (r *v20router) addV20Router() {
	r.router.POST("/auth/login/", func(c *gin.Context) {
		c.Request.ParseMultipartForm(1024*1024)
		for key, values := range c.Request.PostForm {
			log.Printf("Form field %q, Values %q\n", key, values)
		}

		name := c.PostForm("username")
		password := c.PostForm("password")

		/**
    @SerializedName("errorCode") var errorCode: Int = -1,
    @SerializedName("errorMsg") var errorMsg: String? = "",
    @SerializedName("data") var data: T? = null,
    @SerializedName("time") var time: String? = "",
		*/
		if (name != "kotaku" || password != "kotaku-blog.link") {
			c.JSON(http.StatusBadRequest, gin.H{
				"errorCode":http.StatusBadRequest,
				"errorMsg":"Invalide username or password",
				"time":time.Now(),
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"errorCode":http.StatusOK,
				"errorMsg":"Success",
				"data":gin.H{
					"token":"a48396e4f5bec65ddd415cb802cd37be7a5784cae",
				},
				"time":time.Now(),
			})
		}
  })
}

func v20ReceiveFile(w http.ResponseWriter, r *http.Request) map[string]any {
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

func v20GetAnswer(c *gin.Context) map[string]any {
	c.Request.ParseMultipartForm(1024*1024)
	for key, values := range c.Request.PostForm {
		log.Printf("Form field %q, Values %q\n", key, values)
	}

	return gin.H{
		"auto_token":"a48396e4f5bec65ddd415cb802cd37be7a5784cae",
		"time":time.Now(),
	}
}