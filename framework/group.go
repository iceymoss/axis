package framework

// IGroup 代表前缀分组
type IGroup interface {
	Get(string, ...ControllerHandler)
	Post(string, ...ControllerHandler)
	Put(string, ...ControllerHandler)
	Delete(string, ...ControllerHandler)
	Use(middlewares ...ControllerHandler)
	Group(uri string) IGroup
}

// Group 前缀匹配的具体实现者
type Group struct {
	core        *Core               // 指向core结构
	parent      *Group              //指向上一个Group，如果有的话
	prefix      string              // 这个group的通用前缀
	middlewares []ControllerHandler // 存放中间件
}

func NewGroup(core *Core, prefix string) *Group {
	return &Group{
		core:   core,
		parent: nil,
		prefix: prefix,
	}
}

// Use 注册中间件
func (g *Group) Use(middlewares ...ControllerHandler) {
	g.middlewares = append(g.middlewares, middlewares...)
}

// 获取某个group的middleware
// 这里就是获取除了Get/Post/Put/Delete之外设置的middleware
func (g *Group) getMiddlewares() []ControllerHandler {
	if g.parent == nil {
		return g.middlewares
	}

	return append(g.parent.getMiddlewares(), g.middlewares...)
}

func (g *Group) Get(uri string, handlers ...ControllerHandler) {
	uri = g.prefix + uri
	allHandlers := append(g.getMiddlewares(), handlers...)
	g.core.Get(uri, allHandlers...)
}

func (g *Group) Post(uri string, handlers ...ControllerHandler) {
	uri = g.prefix + uri
	g.core.Post(uri, handlers...)
}

func (g *Group) Put(uri string, handlers ...ControllerHandler) {
	uri = g.prefix + uri
	g.core.Put(uri, handlers...)
}

func (g *Group) Delete(uri string, handler ...ControllerHandler) {
	uri = g.prefix + uri
	g.core.Delete(uri, handler...)
}

// Group 实现 Group 方法
func (g *Group) Group(uri string) IGroup {
	cgroup := NewGroup(g.core, uri)
	cgroup.parent = g
	return cgroup
}
