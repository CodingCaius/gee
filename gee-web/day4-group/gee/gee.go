package gee

import (
    "log"
    "net/http"
)

type HandlerFunc func(*Context)


/*
在整个应用程序中，只创建一个 Engine 实例，并且该实例被所有的路由组共享使用。
通过将所有的路由组都关联到同一个 Engine 实例上，可以实现不同路由组之间的交互和共享资源

通过建立双向关联，进行信息互通和资源共享
*/


//Engine implement the interface of ServerHTTP
type (
    RouterGroup struct {
        prefix string  //路由组的前缀，用于添加到组内每个路由的路径前面
        middlewares []HandlerFunc  //应用于该路由组的中间件函数
        parent *RouterGroup  //指向父路由组的指针，支持嵌套路由组的创建
        engine *Engine   //所有的路由组共享一个引擎实例
    }

    Engine struct {
        *RouterGroup  //匿名字段
        router *router  //实际的路由器，用于处理和匹配路由请求
        groups []*RouterGroup //store all groups
    }
)

//New is the constructor of gee.Engine
func New() *Engine {
    //创建一个新的 Engine 实例，并通过 newRouter() 函数创建一个新的路由器实例，
    //并将其赋值给 router 字段
    engine := &Engine{router: newRouter()}
    //双向连接
    engine.RouterGroup = &RouterGroup{engine: engine}  //根路由组，默认为空
    //将根路由组加入groups
    engine.groups = []*RouterGroup{engine.RouterGroup}
    return engine
}

// Group is defined to create a new RouterGroup
// remember all groups share the same Engine instance
//用于创建并返回一个新的子路由组，prefix指定路由组的前缀
func (group *RouterGroup) Group(prefix string) *RouterGroup {
    //从当前路由组 group 中获取引擎实例，并将其赋值给 engine 变量,供之后赋值
    engine := group.engine
    newGroup := &RouterGroup{
        prefix: group.prefix + prefix,
        parent: group,
        engine: engine,
    }
    engine.groups = append(engine.groups, newGroup)
    return newGroup
}

//将和路由有关的函数，都交给RouterGroup实现
//因为 (*Engine).engine 是指向自己的。
//这样实现，我们既可以像原来一样添加路由，也可以通过分组添加路由

//comp为传入的路由组件
func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
    pattern := group.prefix + comp
    log.Printf("Route %4s - %s", method, pattern)
    group.engine.router.addRoute(method, pattern, handler)
}

// GET defines the method to add GET request
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}


func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}



//engine结构体定义了serveHTTP方法，因此实现了http.Handler接口
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    c := newContext(w, req)  //创建上下文处理请求
    engine.router.handle(c)
}
