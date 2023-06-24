package gee

import (
    "net/http"
    "strings"
)

//为每种请求方式都创建一个独立的 Trie 树用于路由匹配
//用 roots 来存储每种请求方式（如 GET、POST、PUT 等）的 Trie 树的根节点
type router struct {
    roots map[string]*node
    handlers map[string]HandlerFunc
}

func newRouter() *router {
    return &router{
        roots: make(map[string]*node),
        handlers: make(map[string]HandlerFunc),
    }
}

//路由模式中只允许有一个通配符
func parsePattern(pattern string) []string {
    vs := strings.Split(pattern, "/")

    parts := make([]string, 0)
    for _, item := range vs {
        if item != "" {
            parts = append(parts, item)
            if item[0] == '*' {
                break
            }
        }
    }
    return parts
}

//向router添加路由规则
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
    parts := parsePattern(pattern)

    key := method + "-" + pattern
    //先检查请求方法是否存在
    _, ok := r.roots[method]
    if !ok {
        //创建根节点
        r.roots[method] = &node{}
    }
    //调用根节点的 insert 方法
    r.roots[method].insert(pattern, parts, 0)
    r.handlers[key] = handler
}

//根据请求方法和路径查找匹配的路由规则
//返回匹配到的节点（*node）和解析得到的参数（map[string]string）
//如果没有通配符，返回的 参数映射 就为空
func (r *router) getRoute(method string, path string) (*node, map[string]string) {
    searchParts := parsePattern(path)
    //创建存储解析得到的参数的切片
    params := make(map[string]string)
    root, ok := r.roots[method]

    if !ok {
        return nil, nil
    }

    n := root.search(searchParts, 0)

    if n != nil {
        parts := parsePattern(n.pattern)
        for index, part := range parts {
            if part[0] == ':' {
                //将通配符参数作为键，路径中的实际值作为值
                params[part[1:]] = searchParts[index]
            }
            /*
            如果遇到了以 * 开头的部分且长度大于 1，就将该部分的名称（去除了 *）作为键，将路径部分切片 searchParts 中当前位置到结尾的所有值用斜杠 / 连接起来作为值，存储到参数映射 params 中。

            例如，对于路由模式 /static/*filepath 和路径 /static/css/geektutu.css，当遇到 *filepath 时，params["filepath"] 将被赋值为 "css/geektutu.css"。这样，处理函数就可以通过访问 params["filepath"] 来获取匹配到的路径部分。
            */
            if part[0] == '*' && len(part) > 1 {
                params[part[1:]] = strings.Join(searchParts[index:], "/")
                break
            }
        }
        return n, params
    }

    return nil, nil
}

//用于获取指定请求方法下的所有路由节点
func (r *router) getRoutes(method string) []*node {
    root, ok := r.roots[method]
    if !ok {
        return nil
    }

    nodes := make([]*node, 0)
    root.travel(&nodes)
    return nodes
}

//处理路由请求
func (r *router) handle(c *Context) {
    n, params := r.getRoute(c.Method, c.Path)

    if n != nil {
        c.Params = params
        //构建处理函数的键 key，用于在 handlers 映射中查找对应的处理函数
        key := c.Method + "-" + n.pattern
        //r.handlers[key](c)

        //将与路由节点匹配的处理函数 r.handlers[key] 添加到上下文对象 c 的处理函数链 handlers 中
        c.handlers = append(c.handlers, r.handlers[key])
    } else {
        /*
        如果路由节点 n 为空，表示未找到匹配的路由，那么将添加一个默认的处理函数到上下文对象 c 的处理函数链中。
        这个默认的处理函数返回一个包含 "404 NOT FOUND" 错误信息的字符串响应。*/
        c.handlers = append(c.handlers, func(c *Context) {
            c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
        })
    }
    c.Next()  //这里调用Next函数，相当于是整个中间件链的入口
    /*
    与之前直接调用路由处理函数不同的是
    在添加中间件功能后，
    收到路由请求后，ServeHTTP函数先根据前缀添加 会用到的中间件 到上下文里，
    然后调用handle路由处理函数，handle函数在映射中查找对应的处理函数
    并将处理函数添加到中间件链中
    最后调用Next函数，开始执行上下文的中间件链
    */
}