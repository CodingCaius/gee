package gee

import (
    "encoding/json"
    "fmt"
    "net/http"
)

type H map[string]interface{}

type Context struct {
    //origin objects
    Writer http.ResponseWriter
    Req *http.Request
    //request info
    Path string
    Method string
    Params map[string]string //参数映射
    //reponse info
    StatusCode int

    //middleware
    handlers []HandlerFunc
    index int
}


func newContext(w http.ResponseWriter, req *http.Request) *Context {
    return &Context{
        Writer: w,
        Req: req,
        Path: req.URL.Path,
        Method: req.Method,
        index: -1,
    }
}

//该 Next() 方法的目的是在中间件链中继续执行下一个中间件
func (c *Context) Next() {
    c.index++
    //获取中间件的数量
    s := len(c.handlers)
    for ; c.index < s; c.index++ {
        c.handlers[c.index](c)
    }

}

//该 Fail() 方法的作用是终止当前请求的处理过程，并返回一个带有指定状态码和错误信息的 JSON 响应
func (c *Context) Fail(code int, err string) {
    c.index = len(c.handlers)
    c.JSON(code, H{"message": err})
}

//获取路由中的参数值
func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

//该方法用于获取HTTP POST请求中指定键（key）对应的表单数据值（value）
func (c *Context) PostForm(key string) string {
    return c.Req.FormValue(key)
}

//该方法用于获取HTTP请求URL中指定键（key）对应的查询参数值（value）。
func (c *Context) Query(key string) string {
    return c.Req.URL.Query().Get(key)
}

func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

//设置HTTP响应头部的键（key）和值（value）
func (c *Context) SetHeader(key string, value string) {
    c.Writer.Header().Set(key, value)
}

//返回文本格式的HTTP响应
func (c *Context) String(code int, format string, values ...interface{}) {
    c.SetHeader("Content-Type", "text/plain")
    c.Status(code)
    //将格式化后的字符串作为HTTP响应的内容写入到响应写入器中，以便将其发送给客户端
    c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

func (c *Context) JSON(code int, obj interface{}) {
    c.SetHeader("Context-Type", "application/json")
    c.Status(code)
    //创建一个json编码器，用于将json数据写入c.Writer中
    encoder := json.NewEncoder(c.Writer)
    //使用 JSON 编码器将给定的对象（obj）进行编码，并将编码后的 JSON 数据写入到响应写入器中。
    //如果在编码过程中发生错误，将会使用 http.Error 函数将错误信息返回给客户端，并设置状态码为 500。
    if err := encoder.Encode(obj); err != nil {
        http.Error(c.Writer, err.Error(), 500)
    }

}


//将给定的字节数据（data）作为原始数据（binary data）发送给客户端作为 HTTP 响应。
func (c *Context) Data(code int, data []byte) {
    c.Status(code)
    c.Writer.Write(data)
}


func (c *Context) HTML(code int, html string) {
    c.SetHeader("Context-Type", "text/html")
    c.Status(code)
    c.Writer.Write([]byte(html))
}
