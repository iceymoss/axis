package framework

import (
	"context"
	"fmt"
	"log"
	"time"
)

func TimeoutHandler(function ControllerHandler, d time.Duration) ControllerHandler {
	return func(c *Context) error {
		finish := make(chan struct{}, 1)
		panicChan := make(chan interface{}, 1)

		//将ctx添加过期时间
		durationCtx, cancel := context.WithTimeout(c, d)
		defer cancel()

		c.request.WithContext(durationCtx)

		go func() {
			defer func() {
				if p := recover(); p != nil {
					panicChan <- p
				}
			}()
			//执行具体业务逻辑
			function(c)
			//业务执行完毕，通知主协程
			finish <- struct{}{}
		}()
		//主协程阻塞等待
		select {
		case p := <-panicChan: //异常情况，也会在业务出现异常的时候，通过 panicChan 来传递异常信号。
			c.WriterMux().Lock()
			defer c.WriterMux().Unlock()
			log.Println(p)
		case <-finish: //业务结束
			fmt.Println("finish")
		case <-durationCtx.Done(): //超时
			c.WriterMux().Lock()
			defer c.WriterMux().Unlock()
			c.Json(500, "time out")
			c.SetHasTimeout()
		}
		return nil
	}
}
