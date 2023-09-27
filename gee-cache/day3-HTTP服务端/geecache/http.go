 // 提供被其他节点访问的能力(基于http)

package geecache

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

const defaultBasePath = "/_geecache/"

// HTTPPool 为 HTTP 对等点池实现了 PeerPicker。
type HTTPPool struct {
	// this peer's base URL, e.g. "https://example.net:8000"
	self     string // 用来记录自己的地址，包括主机名/IP 和端口
	basePath string // 作为节点间通讯地址的前缀，默认是 /_geecache/，那么 http://example.com/_geecache/ 开头的请求，就用于节点间的访问
	// 因为一个主机上还可能承载其他的服务，加一段 Path 是一个好习惯。比如，大部分网站的 API 接口，一般以 /api 作为前缀
}

// NewHTTPPool 初始化 HTTP 对等点池
func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

// 带有服务器名称的日志信息
func (p *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}

// ServeHTTP 处理所有 http 请求
func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 首先判断访问路径的前缀是否是 basePath，不是返回错误
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		panic("HTTPPool serving unexpected path: " + r.URL.Path)
	}
	p.Log("%s %s", r.Method, r.URL.Path)
	// /<basepath>/<groupname>/<key> required
	// 这部分代码将请求的路径分割成两部分，使用 / 分割
	// 从 r.URL.Path 中去掉 p.basePath 部分，然后通过 / 分割得到两个部分。
	// 如果分割后的部分不等于2，表示请求路径不符合预期，会返回 HTTP 400 Bad Request 响应
	parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	groupName := parts[0]
	key := parts[1]

	// 获取缓存分组:
	group := GetGroup(groupName)
	if group == nil {
		http.Error(w, "no such group: "+groupName, http.StatusNotFound)
		return
	}

	// 获取缓存数据并返回
	view, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 设置响应头并写入响应数据
	w.Header().Set("Content-Type", "application/octet-stream") // 二进制流
	w.Write(view.ByteSlice())
}