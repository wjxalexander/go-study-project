// housing application
package app

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jingxinwangdev/go-prject/internal/api"
)

type Application struct {
	Logger         *log.Logger
	WorkoutHandler *api.WorkoutHandler
}

/**
* 无法共享状态（最致命的问题）
* 这是最核心的区别。Application 通常被用作“依赖注入”的容器。
* 返回指针 (&)：所有的模块（比如你的 API 处理器、中间件）拿到的都是同一个内存地址。如果你在某个地方修改了 app.Logger，全局都会生效。
* 返回对象 ({})：每个模块拿到的都是一个副本。如果你在 A 模块修改了 app 的某个状态，B 模块是感知不到的，因为它们操作的是完全不同的内存块。
 */

func NewApplication() (*Application, error) {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	workoutHandler := api.NewWorkoutHandler()
	app := &Application{
		Logger:         logger,
		WorkoutHandler: workoutHandler,
	}
	return app, nil
}

// http.ResponseWriter 实际上是一个接口（Interface），而不是一个结构体（Struct）。
// 接口变量：http.ResponseWriter 就像是一个盒子，里面已经装着那个指针了。
// 3. 调用方法会变得很麻烦 接口的设计初衷是让你直接调用它的方法。

func (app *Application) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Status is available\n")
}
