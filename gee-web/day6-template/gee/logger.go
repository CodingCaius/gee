package gee

import (
    "log"
    "time"
)

/*
该中间件的主要作用是在处理HTTP请求时记录请求的开始时间，并在请求完成后计算并输出请求的处理时间和其他相关信息。
它可以用作一个全局的日志中间件，用于记录每个请求的处理情况，方便调试、性能监控和故障排查。
*/

//用于记录请求日志的中间件函数
func Logger() HandlerFunc {
    return func(c *Context) {
        //Start time
        t := time.Now()
        //调用下一个中间件或请求处理函数
        c.Next()
        //计算并记录请求的处理时间和其他相关信息
        log.Printf("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(t))
    }
}