package main

import "github.com/iceymoss/axis/framework"

func FooControllerHandler(ctx *framework.Context) error {
	return ctx.Json(200, map[string]interface{}{"code": 0})
}
