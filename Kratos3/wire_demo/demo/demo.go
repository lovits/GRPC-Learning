package demo

import (
	"errors"

	"github.com/google/wire"
)

// 提供者函数必须是可导出的（首字母大写）以便被其他包导入。
// Wire中的提供者就是一个可以产生值的普通函数
type X struct {
	Value int
}

func NewX() X {
	return X{Value: 7}
}

// 提供者函数可以使用参数指定依赖项：
type Y struct {
	Value int
}

func NewY(x X) Y {
	return Y{Value: x.Value + 1}
}

// 提供者函数也是可以返回错误的。
type Z struct {
	Value int
}

func NewZ(y Y) (Z, error) {
	if y.Value == 0 {
		return Z{}, errors.New("bad y")
	}
	return Z{Value: y.Value + 2}, nil
}

// 提供者集合
var ProviderSet = wire.NewSet(NewX, NewY, NewZ)
