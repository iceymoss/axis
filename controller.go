package main

import (
	"context"
	"fmt"
	"github.com/iceymoss/axis/framework"
	"log"
	"time"
)

func FooControllerHandler(c *framework.Context) error {
	finish := make(chan struct{}, 1)
	panicChan := make(chan interface{}, 1)

	durationCtx, cancel := context.WithTimeout(c.BaseContext(), time.Duration(2*time.Second))
	defer cancel()

	// mu := sync.Mutex{}
	go func() {
		defer func() {
			if p := recover(); p != nil {
				panicChan <- p
			}
		}()
		// Do real action
		time.Sleep(10 * time.Second)
		c.Json("ok")

		finish <- struct{}{}
	}()
	select {
	case p := <-panicChan:
		c.WriterMux().Lock()
		defer c.WriterMux().Unlock()
		log.Println(p)
		c.Json("panic")
	case <-finish:
		fmt.Println("finish")
	case <-durationCtx.Done():
		c.WriterMux().Lock()
		defer c.WriterMux().Unlock()
		c.Json("time out")
		c.SetHasTimeout()
	}
	return nil
}

func GetUserListController(ctx *framework.Context) error {
	ctx.Json(map[string]interface{}{
		"code":  200,
		"error": nil,
		"msg":   "成功",
		"data": []map[string]string{
			{"name": "iceymoss", "age": "18"},
			{"name": "kos", "age": "20"},
		},
	})
	return nil
}

func SubjectDelController(ctx *framework.Context) error {
	ctx.Json(map[string]interface{}{
		"code":  200,
		"error": nil,
		"msg":   "成功",
		"data":  "hello",
	})
	return nil
}

func UserLoginController(c *framework.Context) error {
	foo, _ := c.QueryString("foo", "def")
	// 等待10s才结束执行
	time.Sleep(10 * time.Second)
	// 输出结果
	c.SetOkStatus().Json("ok, UserLoginController: " + foo)
	return nil
}
