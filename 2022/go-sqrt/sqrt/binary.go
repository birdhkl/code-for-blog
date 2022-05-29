package sqrt

import (
	"fmt"
)

const ResultPrecise = 0.001

type SqrtValue interface {
	~float32 | ~float64
}

// BinarySqrt
// float可以精确到有效数字8位， double到了17位
// 精度最好低于数据类型(float)最大精度一个数量级，否则会由于 (lo+hi)/2 一直等于 lo 而进入死循环
func BinarySqrt[T SqrtValue](value T) T {
	if value < 0 {
		panic(fmt.Sprintf("BinarySqrt Negative Value %f", value))
	}
	if value == 0 {
		return 0
	}
	min := T(0.0)
	max := value
	const SquareRootPrecise = 10e-6
	for (max - min) > SquareRootPrecise {
		mid := min + (max-min)/2.
		delta := mid*mid - value
		if -ResultPrecise <= delta && delta <= ResultPrecise {
			min = mid
			break
		}
		if delta > 0 {
			max = mid
		} else {
			min = mid
		}
	}
	return min
}
