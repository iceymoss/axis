package framework

import (
	"log"
	"net/http"
)

// Core 框架核心结构
type Core struct {
	router map[string]ControllerHandler
}

// NewCore 初始化对象Core
func NewCore() *Core {
	return &Core{router: map[string]ControllerHandler{}}
}

func (c *Core) Get(url string, handler ControllerHandler) {
	c.router[url] = handler
}

// 框架核心结构实现Handler接口
func (c *Core) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	log.Println("ServeHTTP")
	ctx := NewContext(request, response)

	//TODO：路由匹配算法
	router := c.router["foo"]
	if router == nil {
		return
	}

	log.Println("core.router")
	router(ctx)
}
