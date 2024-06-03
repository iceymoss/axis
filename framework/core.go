package framework

import (
	"log"
	"net/http"
	"strings"
)

// Core 框架核心结构
// 框架的核心，负责路由匹配，server接口实现，初始路由，注册方法
type Core struct {
	//router map[string]ControllerHandler
	//二级哈希，一级存储对应Method，二级存储对应Method下URI
	//router map[string]map[string]ControllerHandler

	// 前缀树路由匹配
	router map[string]*Tree // all routers

	// 从core这边设置的中间件
	middlewares []ControllerHandler
}

// NewCore 初始化对象Core
func NewCore() *Core {
	methods := []string{"GET", "POST", "PUT", "DELETE"}
	router := make(map[string]*Tree)
	for _, method := range methods {
		router[method] = NewTree()
	}
	return &Core{router: router}
}

// Use 注册中间件
func (c *Core) Use(middlewares ...ControllerHandler) {
	c.middlewares = append(c.middlewares, middlewares...)
}

// Get GET方法路由注册
func (c *Core) Get(url string, handlers ...ControllerHandler) {
	//大小写不敏感，增加用户容错
	upperUrl := strings.ToUpper(url)
	if err := c.router["GET"].AddRouter(upperUrl, handlers); err != nil {
		log.Fatal("add router error: ", err)
	}
}

// Post POST方法路由注册
func (c *Core) Post(url string, handlers ...ControllerHandler) {
	upperUrl := strings.ToUpper(url)
	if err := c.router["POST"].AddRouter(upperUrl, handlers); err != nil {
		log.Fatal("add router error: ", err)
	}
}

// Put PUT方法路由注册
func (c *Core) Put(url string, handlers ...ControllerHandler) {
	upperUrl := strings.ToUpper(url)
	if err := c.router["PUT"].AddRouter(upperUrl, handlers); err != nil {
		log.Fatal("add router error: ", err)
	}
}

// Delete DELETE方法路由注册
func (c *Core) Delete(url string, handlers ...ControllerHandler) {
	upperUrl := strings.ToUpper(url)
	if err := c.router["DELETE"].AddRouter(upperUrl, handlers); err != nil {
		log.Fatal("add router error: ", err)
	}
}

func (c *Core) Group(prefix string) IGroup {
	return NewGroup(c, prefix)
}

// FindRouteByRequest 匹配路由，如果没有匹配到，返回nil
func (c *Core) FindRouteByRequest(request *http.Request) []ControllerHandler {
	// uri 和 method 全部转换为大写，保证大小写不敏感
	uri := request.URL.Path
	method := request.Method
	upperMethod := strings.ToUpper(method)

	// 查找第一层map
	if methodHandlers, ok := c.router[upperMethod]; ok {
		return methodHandlers.FindHandler(uri)
	}
	return nil
}

// ServeHTTP 框架核心结构实现了Handler接口
// 所有请求都进入这个函数, 这个函数负责路由分发
func (c *Core) ServeHTTP(response http.ResponseWriter, request *http.Request) {

	// 封装自定义context
	ctx := NewContext(request, response)

	// 寻找路由
	handlers := c.FindRouteByRequest(request)
	if handlers == nil {
		// 如果没有找到，这里打印日志
		ctx.Json(404, "not found")
		return
	}

	ctx.SetHandlers(handlers)

	// 调用路由函数，如果返回err 代表存在内部错误，返回500状态码
	if err := ctx.Next(); err != nil {
		ctx.Json(500, "inner error")
		return
	}
}
