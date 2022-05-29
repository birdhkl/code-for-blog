package sqrt

import (
	"fmt"
)

const ResultPrecise = 0.001

type SqrtValue interface {
	~float32 | ~float64
}

func BinarySqrt[T SqrtValue](value T) T {
	if value < 0 {
		panic(fmt.Sprintf("BinarySqrt Negative Value %f", value))
	}
	if value == 0 {
		return 0
	}
	min := T(0.0)
	max := value
	const SquareRootPrecise = 0.00000001
	for min+SquareRootPrecise < max {
		mid := min + (max-min)/2.0
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
