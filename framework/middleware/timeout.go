package middleware

import (
	"context"
	"fmt"
	"github.com/iceymoss/axis/framework"
	"log"
	"time"
)

func Timeout(d time.Duration) framework.ControllerHandler {
	return func(c *framework.Context) error {
		finish := make(chan struct{}, 1)
		panicChan := make(chan interface{}, 1)
		// 执行业务逻辑前预操作：初始化超时context
		durationCtx, cancel := context.WithTimeout(c.BaseContext(), d)
		defer cancel()

		go func() {
			defer func() {
				if p := recover(); p != nil {
					panicChan <- p
				}
			}()
			// 使用next执行具体的业务逻辑
			c.Next()
			finish <- struct{}{}
		}()

		select {
		case p := <-panicChan:
			log.Println(p)
			c.Json(500, "")
		case <-finish:
			fmt.Println("finish")
		case <-durationCtx.Done():
			c.SetHasTimeout()
			c.Json(500, "time out")
		}
		return nil
	}
}
