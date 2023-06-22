package gee

import (
    "log"
    "net/http"
)

type HandlerFunc func(*Context)

//Engine implement the interface of ServerHTTP
type Engine struct {
    router *router
}

//New is the constructor of gee.Engine
func New() *Engine {
    return &Engine{router: newRouter()} //创建一个新的 Engine 实例，并通过 newRouter() 函数创建一个新的路由器实例，并将其赋值给 router 字段
}

func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
    log.Printf("Route %4s - %s", method, pattern)
    engine.router.addRoute(method, pattern, handler)
}

// GET defines the method to add GET request
func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}


func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}

func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}



//engine结构体定义了serveHTTP方法，因此实现了http.Handler接口
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    c := newContext(w, req)  //创建上下文处理请求
    engine.router.handle(c)
}
