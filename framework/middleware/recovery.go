package middleware

import "github.com/iceymoss/axis/framework"

func Recovery() framework.ControllerHandler {
	return func(c *framework.Context) error {
		defer func() {
			// 核心在增加这个recover机制，捕获c.Next()出现的panic
			if err := recover(); err != nil {
				c.Json(500, err)
			}
		}()
		// 使用next执行具体的业务逻辑
		c.Next()
		return nil
	}
}
