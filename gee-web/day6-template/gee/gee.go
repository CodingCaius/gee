package gee

import (
    "log"
    "net/http"
    "strings"
    "html/template"
    "path"
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

        htmlTemplates *template.Template
        funcMap template.FuncMap
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

// Use is defined to add middleware to the group
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
    group.middlewares = append(group.middlewares, middlewares...)
}


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




/*
通过调用 createStaticHandler 方法可以创建一个处理静态文件的路由处理器，并将其注册到路由器中，以便在收到对应路径的请求时提供静态文件的访问。
*/

//创建静态文件处理器
//relativePath 服务器中存放静态资源的相对路径
func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
    //得到静态文件的绝对路径
    absolutePath := path.Join(group.prefix, relativePath)
    //创建一个文件服务器 fileServer，它会剥离请求路径中的 absolutePath 部分，以便在文件系统中查找和提供静态文件
    fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
    return func(c *Context) {
        file := c.Param("filepath")
        // Check if file exists and/or if we have permission to access it
        if _, err := fs.Open(file); err != nil {
            c.Status(http.StatusNotFound)
            return
        }

        //将请求交给文件服务器处理，将文件的内容作为响应返回给客户端
        fileServer.ServeHTTP(c.Writer, c.Req)
    }
}

//用于在路由组中注册静态文件服务的方法
func (group *RouterGroup) Static(relativePath string, root string) {

	handler := group.createStaticHandler(relativePath, http.Dir(root))
    //构建了一个完整的 URL 模式，表示匹配 relativePath 前缀下的所有路径，并将通配符部分作为参数 "filepath" 提供给处理函数
	urlPattern := path.Join(relativePath, "/*filepath")
	// 将 URL 模式和处理函数注册为 GET 请求的处理器，即当匹配到该 URL 模式的 GET 请求时，将会执行对应的处理函数来提供静态文件服务
	group.GET(urlPattern, handler)
}

// 设置模板引擎的自定义函数映射
func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}

//用于加载指定模式下的HTML模板文件
func (engine *Engine) LoadHTMLGlob(pattern string) {
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
}


func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

//engine结构体定义了serveHTTP方法，因此实现了http.Handler接口

//ServeHTTP 函数也有变化，当我们接收到一个具体请求时，要判断该请求适用于哪些中间件，在这里我们简单通过 URL 的前缀来判断。得到中间件列表后，赋值给 c.handlers
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    var middlewares []HandlerFunc //用于存储匹配到的中间件函数
    for _, group := range engine.groups {  //通过遍历路由组来查找
        if strings.HasPrefix(req.URL.Path, group.prefix) {
            middlewares = append(middlewares, group.middlewares...)
        }
    }

    c := newContext(w, req)  //创建上下文处理请求
    //将处理函数从路由组中添加到上下文中
    c.handlers = middlewares
    c.engine = engine
    engine.router.handle(c)
}
