//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"bubble/internal/biz"
	"bubble/internal/conf"
	"bubble/internal/data"
	"bubble/internal/server"
	"bubble/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.

/*
wire的使⽤用：
1. ⼀一般⼤大型项⽬目会⽤用到。
2. 编写代码时注意使⽤用依赖注⼊入，把⽤用到的依赖项使⽤用参数传⼊入，⽽而不不是⾃自⼰己直接写死。
3. 使⽤用wire把构造函数连接起来，编写⼀一个注⼊入器器。
4. 命令⾏行行⼯工具wire⽣生成Go代码到wire_gen.go⽂文件。
5. 调⽤用 wire_gen.go 中⽣生成的函数。
注意的地⽅方：
1. wire.go 最上⾯面要加 //go:build wireinject
2. wire.go 需要和最终产出对象在同⼀一个包内。在哪⾥里里⽤用就在哪⾥里里创建wire.go⽂文件
*/
func wireApp(*conf.Server, *conf.Data, log.Logger) (*kratos.App, func(), error) {
	//利用wire生成项目调用层级代码，
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
