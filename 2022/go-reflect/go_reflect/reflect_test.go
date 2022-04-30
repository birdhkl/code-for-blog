package go_reflect_test

import (
	"go-reflect/go_reflect"
	"reflect"
	"testing"
)

type Tag struct {
	s string `test:"B"`
	n int64  `test:"A"`
}

func (t Tag) Test1() string {
	return "Test1" + t.s
}

func (t *Tag) Test2() string {
	return "Test2" + t.s
}

type DstBindTag struct {
	A string         `Bind:"A"`
	B int64          `Bind:"B"`
	C bool           `Bind:"C"`
	D [2]int         `Bind:"D"`
	E map[string]int `Bind:"E"`
	F []string       `Bind:"F"`
}

type SrcBindTag struct {
	A string
	B int64
	C bool
	D [2]int
	E map[string]int
	F []string
}

// TestAssertType 测试动态比较reflect.Type
func TestAssertType(t *testing.T) {
	testCases := []interface{}{
		false, "",
		int(1), int8(1), int16(1), int32(1), int64(1),
		uint(1), uint8(1), uint16(1), uint32(1), uint64(1),
		float32(1), float64(1),
		map[string]int{},
	}
	for _, testCase := range testCases {
		if err := go_reflect.AssertType(reflect.TypeOf(testCase)); err != nil {
			t.Errorf("%v not right err=%v", testCase, err)
		}
	}
}

// TestBindAttr  绑定属性
func TestBindAttr(t *testing.T) {
	const constA = "a"
	const constB = 1
	const constC = true
	constD := [2]int{1, 2}
	constE := map[string]int{constA: constB}
	constF := []string{constA}

	testCase := struct {
		src *SrcBindTag
		dst *DstBindTag
		tag string
	}{
		src: &SrcBindTag{A: constA, B: constB, C: constC, D: constD, E: constE, F: constF},
		dst: &DstBindTag{},
		tag: "Bind",
	}

	if err := go_reflect.BindAttr(testCase.dst, testCase.src, testCase.tag); err != nil {
		t.Errorf("%v not right err=%v", testCase, err)
	}
	if constA != testCase.dst.A {
		t.Errorf("dst.A = %s, not %s", testCase.dst.A, constA)
	}
	if constB != testCase.dst.B {
		t.Errorf("dst.B = %d, not %d", testCase.dst.B, constB)
	}
	if constC != testCase.dst.C {
		t.Errorf("dst.C = %v, not %v", testCase.dst.C, constC)
	}
	if !reflect.DeepEqual(constD, testCase.dst.D) {
		t.Errorf("dst.D = %v, not %v", testCase.dst.D, constD)
	}
	if !reflect.DeepEqual(constE, testCase.dst.E) {
		t.Errorf("dst.E = %v, not %v", testCase.dst.E, constE)
	}
	if !reflect.DeepEqual(constF, testCase.dst.F) {
		t.Errorf("dst.F = %v, not %v", testCase.dst.F, constF)
	}
	testCase.dst.D[0] = 11111
	testCase.dst.E[constA] = 11111
	testCase.dst.F[0] = "not right"
	if !reflect.DeepEqual(constD, testCase.src.D) {
		t.Errorf("dst.D = %v, not %v", testCase.dst.D, constD)
	}
	if !reflect.DeepEqual(constE, testCase.src.E) {
		t.Errorf("dst.E = %v, not %v", testCase.dst.E, constE)
	}
	if !reflect.DeepEqual(constF, testCase.src.F) {
		t.Errorf("dst.F = %v, not %v", testCase.dst.F, constF)
	}
	// 绑定错误边界
	type DstBindErrorNotAttrTag struct {
		G int `Bind:"G"`
	}
	if err := go_reflect.BindAttr(&DstBindErrorNotAttrTag{}, testCase.src, testCase.tag); err == nil {
		t.Errorf("not G, want err")
	}
	// 绑定错误边界
	type DstBindErrorTypeTag struct {
		F int `Bind:"F"`
	}
	if err := go_reflect.BindAttr(&DstBindErrorTypeTag{}, testCase.src, testCase.tag); err == nil {
		t.Errorf("F want []string, want err")
	}
}

// TestStruct 测试反射结构体
func TestStruct(t *testing.T) {
	tag := Tag{s: "1", n: 2}
	tagType := reflect.TypeOf(tag)
	// 反射获取field数量
	if tagType.NumField() != 2 {
		t.Errorf("Tag NumFiled Not 2")
	}
	// 反射获取结构体种类
	if tagType.Kind() != reflect.Struct {
		t.Errorf("Tag Kind Not Struct")
	}
	// 反射获取指针
	if reflect.TypeOf(&tag).Kind() != reflect.Ptr {
		t.Errorf("*Tag Kind Not Ptr")
	}
	// 获取属性的tag
	field := tagType.Field(0)
	if field.Tag.Get("test") != "B" {
		t.Errorf("Tag[0] not s,tag not B")
	}
	// 通过名称获取属性
	wantField, ok := tagType.FieldByName("s")
	if !ok || !reflect.DeepEqual(field, wantField) {
		t.Errorf("Tag.s %v != %v", field, wantField)
	}
	// 校验属性Type和Kind是否runtime唯一
	if field.Type != reflect.TypeOf("") {
		t.Errorf("string Type Not Unique, %v", field.Type)
	}
	if field.Type.Kind() != reflect.String {
		t.Errorf("string %s kind not string", field.Type)
	}
	// 获取属性的tag
	field = tagType.Field(1)
	if field.Tag.Get("test") != "A" {
		t.Errorf("Tag[0] not s,tag not A")
	}
	// 获取值(type, unsafe.Pointer)
	tagField := reflect.ValueOf(tag)
	if tagField.NumField() != 2 {
		t.Errorf("Tag NumFiled Not 2")
	}
	if tagField.NumMethod() != 1 {
		t.Errorf("Tag NumMethod %d Not 1", tagField.NumMethod())
	}
	// 入口参数数量
	if tagField.Method(0).Type().NumIn() != 0 {
		t.Errorf("Tag Method Pointer")
	}
	// 非指针结构体调用，参数nil
	res := tagField.Method(0).Call(nil)
	if len(res) != 1 {
		t.Errorf("method not Test1")
	}
	if res, ok := res[0].Interface().(string); !ok || res != tag.Test1() {
		t.Errorf("method not Test1")
	}
	// 指针结构体调用，参数空
	tagField = reflect.ValueOf(&tag)
	res = tagField.Method(0).Call([]reflect.Value{})
	if len(res) != 1 {
		t.Errorf("method not Test1")
	}
	if res, ok := res[0].Interface().(string); !ok || res != tag.Test1() {
		t.Errorf("method not Test1")
	}
}
