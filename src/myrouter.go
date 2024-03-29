package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type MyRouter struct {
	router *mux.Router
}

type Middleware func(http.HandlerFunc) http.HandlerFunc

// 把应用到http.HandlerFunc处理器的中间件
// 按照先后顺序和处理器本身链起来供http.HandleFunc调用
func Chain(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, m := range middlewares {
			f = m(f)
	}
	return f
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

// 记录每个URL请求的执行时长
func Logging() Middleware {

	// 创建中间件
	return func(f http.HandlerFunc) http.HandlerFunc {

			// 创建一个新的handler包装http.HandlerFunc
			return func(w http.ResponseWriter, r *http.Request) {

					// 中间件的处理逻辑
					log.Println("Http Start ===>>> " + r.Method + " : " + r.URL.Path) 

					// 调用下一个中间件或者最终的handler处理程序
					f(w, r)

					log.Println("Http End <<<=== " + r.URL.Path) 
			}
	}
}

func (r *MyRouter) HandleFunc(path string, f func(http.ResponseWriter,
	*http.Request), method string) {
		r.router.HandleFunc(path, Chain(f, Method(method), Logging()))
}


