package sqrt

import (
	"fmt"
)

const PRECISE = 0.001

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
	for min < max {
		mid := (min + max) / 2.0
		delta := mid*mid - value
		if -PRECISE <= delta && delta <= PRECISE {
			return mid
		}
		if delta > 0 {
			max = mid
		} else {
			min = mid
		}
	}
	return min
}
