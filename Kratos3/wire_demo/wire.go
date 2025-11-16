//go:build wireinject
// +build wireinject

package main

import (
	"main/demo"

	"github.com/google/wire"
)

func initZ() (demo.Z, error) {
	//传统写法
	// x := demo.NewX()
	// y := demo.NewY(x)
	// z, err := demo.NewZ(y)
	// fmt.Println(z, err)
	// return z, err

	//应用程序中是用一个注入器来连接提供者，注入器就是一个按照依赖顺序调用提供者。
	// wire.Build(demo.NewX, demo.NewY, demo.NewZ)
	wire.Build(demo.ProviderSet) //有提供者集合可以这样写
	return demo.Z{}, nil

}
