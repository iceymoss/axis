package main

import (
	"github.com/iceymoss/axis/framework"
	"github.com/iceymoss/axis/framework/middleware"
	"time"
)

//func registerRouter(core *framework.Core) {
//	// core.Get("foo", framework.TimeoutHandler(FooControllerHandler, time.Second*1))
//	core.Get("foo", FooControllerHandler)
//	core.Get("hello")
//}

// RegisterRouter 注册路由规则
func RegisterRouter(core *framework.Core) {
	// 需求1+2:HTTP方法+静态路由匹配
	core.Get("/user/list", framework.TimeoutHandler(GetUserListController, 1*time.Second))
	core.Get("/user/test", UserLoginController)
	// 需求3:批量通用前缀
	subjectApi := core.Group("/subject")
	{
		// 需求4:动态路由
		subjectApi.Delete("/:id", SubjectDelController)
		subjectApi.Put("/:id", SubjectDelController)
		subjectApi.Get("/:id", SubjectDelController)
		subjectApi.Get("/list/all", SubjectDelController)
	}
}

// 注册路由规则
func registerRouter(core *framework.Core) {
	// 在core中使用middleware.Test3() 为单个路由增加中间件
	core.Get("/user/list", middleware.Test3(), GetUserListController)
	core.Get("/user/test", UserLoginController)
	// 批量通用前缀
	subjectApi := core.Group("/subject")
	{
		// 在group中使用middleware.Test3() 为单个路由增加中间件
		subjectApi.Get("/:id", middleware.Test3(), SubjectDelController)
	}
}
