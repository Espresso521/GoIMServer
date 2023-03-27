// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

func ReceiveFile(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20) 

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
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
	}
	defer f.Close()

	// copy the uploaded file data to the new file
	_, err = io.Copy(f, file)
	if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
	}

	fStat, _ := f.Stat()

	// return a success message to the client
	fmt.Fprintf(w, "{'auth_token':'a48396e4f5bec65ddd415cb802cd37be7a5784cae', 'time':'%s', 'extra':'File uploaded successfully: file name is %v, file size is %v bytes'}", time.Now(), f.Name(), fStat.Size())
}

type Person struct {
	Name string
	Age  int
}

func DisplayPersonHandler(w http.ResponseWriter, r *http.Request) {
	var p Person

	// 将请求体中的 JSON 数据解析到结构体中
	// 发生错误，返回400 错误码
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
	}

	fmt.Fprintf(w, "Person: %+v", p)
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func serveTest(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/test" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	http.ServeFile(w, r, "test.html")
}

func DisplayFormDataHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(1024*1024) 
	if err := r.ParseForm(); err != nil {
			panic(err)
	}

	for key, values := range r.Form {
		log.Printf("Form field %q, Values %q\n", key, values)
 }

 fmt.Fprintf(w, "{'auth_token':'a48396e4f5bec65ddd415cb802cd37be7a5784cae', 'time':'%s'}", time.Now())
}

func appStartTime(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(1024*1024) 
	if r.URL.Path != "/api/v/1.0/app/start/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	// for key,value:= range r.Header{
	// 		log.Printf("%s=>%s\n",key,value)
	// }

	// for k, v := range r.URL.Query() {
	// 	log.Printf("ParamName %q, Value %q\n", k, v)
	// 	log.Printf("ParamName %q, Get Value %q\n", k, r.URL.Query().Get(k))
	// }

	for key, values := range r.Form {
		log.Printf("Form field %q, Values %q\n", key, values)
 }

	fmt.Fprintf(w, "{'auth_token':'a48396e4f5bec65ddd415cb802cd37be7a5784cae', 'time':'%s'}", time.Now())
}

func authLogin(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(1024*1024) 
	if r.URL.Path != "/auth/token/login/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	for key, values := range r.Form {
		log.Printf("Form field %q, Values %q\n", key, values)
 }

	fmt.Fprintf(w, "{'auth_token':'a48396e4f5bec65ddd415cb802cd37be7a5784cae', 'time':'%s'}", time.Now())
}

func handleHello(w http.ResponseWriter, r *http.Request) {
	res, err := json.Marshal(map[string]string{
			"message": "hello from the server",
	})
	if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func handleSPA(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w,r,"public/index.html")
}

func main() {
	hub := newHub()
	go hub.run()

	// regist router
	myrouter := MyRouter{mux.NewRouter()}

	// Serve the static files for the Svelte app
	myrouter.HandleFunc("/global.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w,r,"public/global.css")
	}, "GET")
	myrouter.HandleFunc("/build/bundle.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w,r,"public/build/bundle.js")
	}, "GET")
	myrouter.HandleFunc("/build/bundle.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w,r,"public/build/bundle.css")
	}, "GET")
	myrouter.HandleFunc("/favicon.png", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w,r,"public/favicon.png")
	}, "GET")
	
	myrouter.HandleFunc("/", serveHome, "GET")
	myrouter.HandleFunc("/test", serveTest, "GET")
	myrouter.HandleFunc("/api/v/1.0/app/start/", appStartTime, "POST")
	myrouter.HandleFunc("/display_form_data", DisplayFormDataHandler, "POST")
	myrouter.HandleFunc("/parse_json_request", DisplayPersonHandler, "POST")
	myrouter.HandleFunc("/upload_file", ReceiveFile, "POST")
	myrouter.HandleFunc("/auth/token/login/", authLogin, "POST")
	myrouter.HandleFunc("/api/v/1.0/run/start/", ReceiveFile, "POST")
	myrouter.HandleFunc("/api/v/1.0/run/status/", DisplayFormDataHandler, "POST")
	myrouter.HandleFunc("/api/v/1.0/run/finish/", ReceiveFile, "POST")
	myrouter.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
			serveWs(hub, w, r)
		}, "GET")

	myrouter.HandleFunc("/hello.json", handleHello, "GET")
	myrouter.HandleFunc("/svelte", handleSPA, "GET")

	log.Println("Server start")
	err := http.ListenAndServe(":5211", myrouter.router)
	if err != nil {
		log.Fatal("ListenAndServeTLS: ", err)
	}

    // Serve static files from the "public" directory
    // fs := http.FileServer(http.Dir("public"))

    // // Handle requests to the root URL ("/") by serving the "index.html" file
    // http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    //     http.ServeFile(w, r, "public/index.html")
    // })
	
		// 	// Start the server on port 8080
		// 	log.Fatal(http.ListenAndServe(":5211", fs))
}