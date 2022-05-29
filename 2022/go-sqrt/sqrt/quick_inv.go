package sqrt

import "math"

// QuickSqrt32 倒排平方根
func QuickInvSqrt32(x float32) float32 {
	const MAGIC_NUMBER = 0x5f375a86
	xhalf := 0.5 * x
	i := math.Float32bits(x)
	i = MAGIC_NUMBER - (i >> 1)
	x = math.Float32frombits(i)
	x = x * (1.5 - xhalf*x*x)
	x = x * (1.5 - xhalf*x*x)
	return x
}

// QuickSqrt32 基于快速平方根倒数的开方算法
func QuickSqrt32(x float32) float32 {
	return 1.0 / QuickInvSqrt32(x)
}

// QuickInvSqrt64
func QuickInvSqrt64(x float64) float64 {
	const MAGIC_NUMBER = 0x5fe6ec85e7de30da
	xhalf := 0.5 * x
	i := math.Float64bits(x)
	i = MAGIC_NUMBER - (i >> 1)
	x = math.Float64frombits(i)
	x = x * (1.5 - xhalf*x*x)
	x = x * (1.5 - xhalf*x*x)
	return x
}

// QuickSqrt64 基于快速平方根倒数的开方算法
func QuickSqrt64(x float64) float64 {
	return 1.0 / QuickInvSqrt64(x)
}
