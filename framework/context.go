package framework

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// Context 自定义 Context
type Context struct {
	request        *http.Request
	responseWriter http.ResponseWriter
	ctx            context.Context
	handler        ControllerHandler

	// 是否超时标记位
	hasTimeout bool

	// 写保护机制
	writerMux *sync.Mutex

	// 当前请求的handler链条
	handlers []ControllerHandler

	// 当前请求调用到调用链的哪个节点
	index int
}

func NewContext(r *http.Request, w http.ResponseWriter) *Context {
	return &Context{
		request:        r,
		responseWriter: w,
		ctx:            r.Context(),
		writerMux:      &sync.Mutex{},
		index:          -1,
	}
}

// #region base function

func (ctx *Context) WriterMux() *sync.Mutex {
	return ctx.writerMux
}

func (ctx *Context) GetRequest() *http.Request {
	return ctx.request
}

func (ctx *Context) SetHasTimeout() {
	ctx.hasTimeout = true
}

// SetHandlers 将涉及到的中间件和控制器加入ctx中
func (ctx *Context) SetHandlers(handlers []ControllerHandler) {
	ctx.handlers = handlers
}

func (ctx *Context) HasTimeout() bool {
	return ctx.hasTimeout
}

func (ctx *Context) GetResponse() http.ResponseWriter {
	return ctx.responseWriter
}

//实现自定义Context实现Context接口

func (ctx *Context) BaseContext() context.Context {
	return ctx.request.Context()
}

func (ctx *Context) Deadline() (deadline time.Time, ok bool) {
	return ctx.BaseContext().Deadline()
}

func (ctx *Context) Err() error {
	return ctx.BaseContext().Err()
}

func (ctx *Context) Value(key interface{}) interface{} {
	return ctx.BaseContext().Value(key)
}

func (ctx *Context) Done() <-chan struct{} {
	return ctx.BaseContext().Done()
}

// #endregion 获取参数

// QueryInt 获取query指定key的int值
func (ctx *Context) QueryIntTest(key string, def int) int {
	params := ctx.QueryAll()
	if values, ok := params[key]; ok {
		length := len(values)
		if length > 0 {
			intVal, err := strconv.Atoi(values[length-1])
			if err != nil {
				return def
			}
			return intVal
		}
	}
	return def
}

// QueryArray 获取query的指定key值
func (ctx *Context) QueryArray(key string, def []string) []string {
	params := ctx.QueryAll()
	if values, ok := params[key]; ok {
		return values
	}
	return def
}

// QueryString 获取query指定key的string值
func (ctx *Context) QueryStringTest(key string, def string) string {
	params := ctx.QueryAll()
	if values, ok := params[key]; ok {
		length := len(values)
		if length > 0 {
			return values[length-1]
		}
	}
	return def
}

// QueryAll 提取所有query参数
// map[string][]string: http://example.com/search?tag=science&tag=technology&tag=math
// result：map["tag"][]string{"science", "technology", "math"}
func (ctx *Context) QueryAllTest() map[string][]string {
	if ctx.request != nil {
		//使用Query()方法解析URL中的查询字符串
		return map[string][]string(ctx.request.URL.Query())
	}
	return map[string][]string{}
}

func (ctx *Context) FormInt(key string, def int) int {
	params := ctx.FormAll()
	if values, ok := params[key]; ok {
		length := len(values)
		if length > 0 {
			intVal, err := strconv.Atoi(values[length-1])
			if err != nil {
				return def
			}
			return intVal
		}
	}
	return def
}

func (ctx *Context) FormString(key string, def string) string {
	params := ctx.FormAll()
	if values, ok := params[key]; ok {
		length := len(values)
		if length > 0 {
			return values[length-1]
		}
	}
	return def
}

// FormAll 提取所有表单参数
func (ctx *Context) FormAll() map[string][]string {
	if ctx.request != nil {
		return map[string][]string(ctx.request.PostForm)
	}
	return map[string][]string{}
}

func (ctx *Context) BindJson(obj interface{}) error {
	if ctx.request != nil {
		body, err := ioutil.ReadAll(ctx.request.Body)
		if err != nil {
			return err
		}
		//原始的 io.ReadCloser 被耗尽(即读到EOF)，因此它不能被再次读取（如再次解析或由多个不同的处理器/中间件使用）。
		//为了能够让后续的处理逻辑也能读取请求体中的内容，需要重新设置 ctx.request.Body。
		ctx.request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		err = json.Unmarshal(body, obj)
		if err != nil {
			return err
		}
	} else {
		return errors.New("ctx.request empty")
	}
	return nil
}

// response

func (ctx *Context) Json(status int, obj interface{}) error {
	if ctx.HasTimeout() {
		return nil
	}
	ctx.responseWriter.Header().Set("Content-Type", "application/json")
	ctx.responseWriter.WriteHeader(status)
	byt, err := json.Marshal(obj)
	if err != nil {
		ctx.responseWriter.WriteHeader(500)
		return err
	}
	ctx.responseWriter.Write(byt)
	return nil
}

func (ctx *Context) HTML(status int, obj interface{}, template string) error {
	return nil
}

func (ctx *Context) Text(status int, obj string) error {
	return nil
}

// Next 执行handlers的下一个方法(中间件+控制器), 通过移动index控制请求调用链
func (ctx *Context) Next() error {
	ctx.index++
	if ctx.index < len(ctx.handlers) {
		if err := ctx.handlers[ctx.index](ctx); err != nil {
			return err
		}
	}
	return nil
}
