package sqrt

import (
	"fmt"
	"math"
)

// NewtonSqrt use Newton-Raphson Method. Core Algorithm is next = (last + m / last) / 2
func NewtonSqrt[T SqrtValue](value T) T {
	if value < 0 {
		panic(fmt.Sprintf("Sqrt A Negative Value %f", value))
	}
	if value == 0 {
		return value
	}
	v := value
	for math.Abs(float64(v*v-value)) > PRECISE {
		v = (v + value/v) / 2.
	}
	return v
}
