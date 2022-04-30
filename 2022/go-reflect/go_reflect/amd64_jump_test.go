// +build !windows amd64

package go_reflect_test

import (
	"fmt"
	"go-reflect/go_reflect"
	"testing"
)

func MonkeyPatchTestFuncA() string {
	res := "A"
	for i := 0; i < 10; i++ {
		res = fmt.Sprintf("%s%d", res, i)
	}
	return res
}

func MonkeyPatchTestFuncB() string {
	res := "B"
	for i := 0; i < 10; i++ {
		res = fmt.Sprintf("%s%d", res, i)
	}
	return res
}

// TestDoMonkeyPatch 猴子补丁测试
func TestDoMonkeyPatch(t *testing.T) {
	const resA = "A0123456789"
	const resB = "B0123456789"
	if MonkeyPatchTestFuncA() != resA || MonkeyPatchTestFuncB() != resB {
		t.Errorf("TestFunc Not Right")
	}
	go_reflect.DoMonkeyPatch(MonkeyPatchTestFuncA, MonkeyPatchTestFuncB)
	if MonkeyPatchTestFuncA() != resB || MonkeyPatchTestFuncB() != resB {
		t.Errorf("Monkey Patch Not Succ")
	}
}
