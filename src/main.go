// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

// 记录每个URL请求的执行时长
func Logging() Middleware {

    // 创建中间件
    return func(f http.HandlerFunc) http.HandlerFunc {

        // 创建一个新的handler包装http.HandlerFunc
        return func(w http.ResponseWriter, r *http.Request) {

            // 中间件的处理逻辑
						log.Println("Logging: " + r.URL.Path) 
 
            // 调用下一个中间件或者最终的handler处理程序
            f(w, r)
        }
    }
}

// 验证请求用的是否是指定的HTTP Method，不是则返回 400 Bad Request
func Method(m string) Middleware {

	return func(f http.HandlerFunc) http.HandlerFunc {

			return func(w http.ResponseWriter, r *http.Request) {

					if r.Method != m {
							http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
							return
					}

					f(w, r)
			}
	}
}

// 把应用到http.HandlerFunc处理器的中间件
// 按照先后顺序和处理器本身链起来供http.HandleFunc调用
func Chain(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, m := range middlewares {
			f = m(f)
	}
	return f
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
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
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "test.html")
	log.Println(r.URL.Path + " END") 
}

func appStartTime(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/api/v/1.0/app/start" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "test.html")
	log.Println(r.URL.Path + " END") 
}

func main() {
	hub := newHub()
	go hub.run()

	router := mux.NewRouter()

	// regist router
	router.HandleFunc("/", serveHome)
	router.HandleFunc("/test", Chain(serveTest, Method("GET"), Logging())).Methods("GET")
	// /api/v/1.0/app/start
	router.HandleFunc("/api/v/1.0/app/start", Chain(appStartTime, Method("POST"), Logging())).Methods("POST")
	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	log.Println("Server start")
	err := http.ListenAndServe(":5211", router)
	if err != nil {
		log.Fatal("ListenAndServeTLS: ", err)
	}
}