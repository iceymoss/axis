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
	router map[string]map[string]ControllerHandler
}

// NewCore 初始化对象Core
func NewCore() *Core {
	methods := []string{"GET", "POST", "PUT", "DELETE"}
	router := make(map[string]map[string]ControllerHandler)
	for _, method := range methods {
		router[method] = map[string]ControllerHandler{}
	}
	return &Core{router: router}
}

// Get GET方法路由注册
func (c *Core) Get(url string, handler ControllerHandler) {
	//大小写不敏感，增加用户容错
	upperUrl := strings.ToUpper(url)
	c.router["GET"][upperUrl] = handler
}

// Post POST方法路由注册
func (c *Core) Post(url string, handler ControllerHandler) {
	upperUrl := strings.ToUpper(url)
	c.router["POST"][upperUrl] = handler
}

// Put PUT方法路由注册
func (c *Core) Put(url string, handler ControllerHandler) {
	upperUrl := strings.ToUpper(url)
	c.router["PUT"][upperUrl] = handler
}

// Delete DELETE方法路由注册
func (c *Core) Delete(url string, handler ControllerHandler) {
	upperUrl := strings.ToUpper(url)
	c.router["DELETE"][upperUrl] = handler
}

// FindRouteByRequest 路由匹配算法，没有匹配返回nil
func (c *Core) FindRouteByRequest(request *http.Request) ControllerHandler {
	//获取Method和URI
	uri := request.URL.Path
	method := request.Method
	//统一转为大写
	upperUri := strings.ToUpper(uri)
	upperMethod := strings.ToUpper(method)
	//匹配Method
	if methodHandlers, ok := c.router[upperMethod]; ok {
		//匹配URI并获取对应的控制器
		if handler, ok := methodHandlers[upperUri]; ok {
			return handler
		}
	}
	return nil
}

// ServeHTTP 框架核心结构实现了Handler接口
func (c *Core) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	log.Println("ServeHTTP")
	ctx := NewContext(request, response)
	router := c.FindRouteByRequest(ctx.request)
	if router == nil {
		ctx.Json(404, "not found")
		return
	}
	err := router(ctx)
	if err != nil {
		ctx.Json(500, "inner error")
		return
	}
}
