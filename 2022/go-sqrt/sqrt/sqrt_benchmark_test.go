package sqrt_test

import (
	"go-sqrt/sqrt"
	"math"
	"testing"
)

func prepareFloat32Case() []float32 {
	return []float32{2.0, 1223.0, 121.0, 999.0}
}

func prepareFloat64Case() []float64 {
	return []float64{2.0, 1223.0, 121.0, 999.0}
}

func BenchmarkBinaryFloat32(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, value := range prepareFloat32Case() {
			sqrt.BinarySqrt(value)
		}
	}

}

func BenchmarkBinaryFloat64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, value := range prepareFloat64Case() {
			sqrt.BinarySqrt(value)
		}
	}
}

func BenchmarkNewtonFloat32(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, value := range prepareFloat32Case() {
			sqrt.NewtonSqrt(value)
		}
	}
}

func BenchmarkNewtonFloat64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, value := range prepareFloat64Case() {
			sqrt.NewtonSqrt(value)
		}
	}
}

func BenchmarkQuickSqrt32(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, value := range prepareFloat32Case() {
			sqrt.QuickSqrt32(value)
		}
	}
}

func BenchmarkQuickSqrt64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, value := range prepareFloat64Case() {
			sqrt.QuickSqrt64(value)
		}
	}
}

func BenchmarkMath(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, value := range prepareFloat64Case() {
			math.Sqrt(value)
		}
	}
}
