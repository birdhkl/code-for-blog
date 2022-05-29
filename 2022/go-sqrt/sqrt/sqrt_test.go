package sqrt_test

import (
	"fmt"
	"go-sqrt/sqrt"
	"math"
	"testing"
)

type sqrtCase[T sqrt.SqrtValue] struct {
	value T
}

func (scase *sqrtCase[T]) isRight(value T) bool {
	wantResult := math.Sqrt(float64(scase.value))
	result := math.Abs(wantResult-float64(value)) <= sqrt.ResultPrecise
	if !result {
		fmt.Printf("not right wantResult:%f, got: %f, dleta: %f\n", wantResult, value, math.Abs(wantResult-float64(value)))
	}
	return result
}

func TestBinary(t *testing.T) {
	cases32 := []sqrtCase[float32]{
		{
			value: 121,
		},
		{
			value: 2,
		},
		{
			value: 18123,
		},
		{
			value: 9912,
		},
	}

	for _, c := range cases32 {
		if !c.isRight(sqrt.BinarySqrt(c.value)) {
			t.Errorf("sqrt %f, result %f", c.value, sqrt.BinarySqrt(c.value))
			return
		}
	}

	cases64 := []sqrtCase[float64]{
		{
			value: 121,
		},
		{
			value: 2,
		},
		{
			value: 18123,
		},
		{
			value: 9912,
		},
	}
	for _, c := range cases64 {
		if !c.isRight(sqrt.BinarySqrt(c.value)) {
			t.Errorf("sqrt %f, result %f", c.value, sqrt.BinarySqrt(c.value))
			return
		}
	}
}

func TestNewton(t *testing.T) {
	cases32 := []sqrtCase[float32]{
		{
			value: 121,
		},
		{
			value: 2,
		},
		{
			value: 18123,
		},
		{
			value: 9912,
		},
	}

	for _, c := range cases32 {
		if !c.isRight(sqrt.NewtonSqrt(c.value)) {
			t.Errorf("sqrt %f, result %f", c.value, sqrt.BinarySqrt(c.value))
			return
		}
	}

	cases64 := []sqrtCase[float64]{
		{
			value: 121,
		},
		{
			value: 2,
		},
		{
			value: 18123,
		},
		{
			value: 9912,
		},
	}
	for _, c := range cases64 {
		if !c.isRight(sqrt.NewtonSqrt(c.value)) {
			t.Errorf("sqrt %f, result %f", c.value, sqrt.BinarySqrt(c.value))
			return
		}
	}
}

func TestQuick(t *testing.T) {
	cases32 := []sqrtCase[float32]{
		{
			value: 121,
		},
		{
			value: 2,
		},
		{
			value: 18123,
		},
		{
			value: 9912,
		},
	}

	for _, c := range cases32 {
		if !c.isRight(sqrt.QuickSqrt32(c.value)) {
			t.Errorf("QuickSqrt32 %f, result %f", c.value, sqrt.QuickSqrt32(c.value))
			return
		}
	}

	cases64 := []sqrtCase[float64]{
		{
			value: 121,
		},
		{
			value: 2,
		},
		{
			value: 18123,
		},
		{
			value: 9912,
		},
	}

	for _, c := range cases64 {
		if !c.isRight(sqrt.QuickSqrt64(c.value)) {
			t.Errorf("QuickSqrt64 %f, result %f", c.value, sqrt.QuickSqrt64(c.value))
			return
		}
	}
}
