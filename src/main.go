// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func main() {
	hub := newHub()
	go hub.run()

	defaultRouter := gin.Default()
	addDefaultRouter(defaultRouter, hub)
  // Optionally, set up additional routes using a custom router group
  v10Router := V10Router{defaultRouter.Group("/api/v/1.0")}
	v10Router.addV10Router()

  // Start the server
	log.Println("Server start")
  defaultRouter.Run(":5211")
}

func addDefaultRouter(defaultRouter *gin.Engine, hub *Hub) {
  // Serve static files from the Svelte build directory
  defaultRouter.Use(static.Serve("/", static.LocalFile("./public", true)))
	defaultRouter.GET("/ws", func(c *gin.Context) {
		serveWs(hub, c.Writer, c.Request)
	})
	defaultRouter.GET("/hello.json", func(ctx *gin.Context) {
		handleHello(ctx.Writer, ctx.Request)
	})
	defaultRouter.POST("/auth/token/login/", func(c *gin.Context) {
		c.Request.ParseMultipartForm(1024*1024)
		for key, values := range c.Request.PostForm {
			log.Printf("Form field %q, Values %q\n", key, values)
	  }
	  c.JSON(http.StatusOK, getAnswer())
	})
}

func handleHello(w http.ResponseWriter, r *http.Request) {

	v := fmt.Sprintf("This message come from GO backend server. Server Time Now is : %s", time.Now())

	res, err := json.Marshal(map[string]string{
			"message": v,
	})
	if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func getAnswer() string {
	ret := fmt.Sprintf("{'auth_token':'a48396e4f5bec65ddd415cb802cd37be7a5784cae', 'time':'%s'}", time.Now())
	return ret
}