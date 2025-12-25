// housing application
package app

import (
	"log"
	"os"
)

type Application struct {
	Logger *log.Logger
}

/**
* 无法共享状态（最致命的问题）
* 这是最核心的区别。Application 通常被用作“依赖注入”的容器。
* 返回指针 (&)：所有的模块（比如你的 API 处理器、中间件）拿到的都是同一个内存地址。如果你在某个地方修改了 app.Logger，全局都会生效。
* 返回对象 ({})：每个模块拿到的都是一个副本。如果你在 A 模块修改了 app 的某个状态，B 模块是感知不到的，因为它们操作的是完全不同的内存块。
 */

func NewApplication() (*Application, error) {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	app := &Application{
		Logger: logger,
	}
	return app, nil
}
