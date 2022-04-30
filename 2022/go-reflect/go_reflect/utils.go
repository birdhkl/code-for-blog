package go_reflect

import (
	"fmt"
	"reflect"
	"unicode"
)

// AssertType 断言t的种类符合原生数据类型
func AssertType(t reflect.Type) error {
	switch t.Kind() {
	case reflect.Bool, reflect.String:
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
	case reflect.Float32, reflect.Float64:
	case reflect.Array, reflect.Slice:
		return AssertType(t.Elem())
	case reflect.Map:
		if err := AssertType(t.Key()); err != nil {
			return err
		}
		if err := AssertType(t.Elem()); err != nil {
			return err
		}
	default:
		return fmt.Errorf("type %v not right", t.Kind())
	}
	return nil
}

// BindAttr 根据Tag将源数据绑定到目标数据 根据dstData的tag标签对应的属性名，从srcData中获取并设置到dstData中
func BindAttr(dstData interface{}, srcData interface{}, tag string) error {
	dstType := reflect.TypeOf(dstData)
	dstVal := reflect.ValueOf(dstData)
	srcType := reflect.TypeOf(srcData)
	srcVal := reflect.ValueOf(srcData)

	if dstType.Kind() == reflect.Ptr {
		dstType = dstType.Elem()
		dstVal = dstVal.Elem()
	}

	if srcType.Kind() == reflect.Ptr {
		srcType = srcType.Elem()
		srcVal = srcVal.Elem()
	}

	if dstType.Kind() != reflect.Struct || srcType.Kind() != reflect.Struct {
		return fmt.Errorf("srcType %v, dstType %v", srcType, dstType)
	}

	dstLen := dstType.NumField()
	srcLen := srcType.NumField()

	dstTagAttrs := map[string]int{}
	for dstIndex := 0; dstIndex < dstLen; dstIndex++ {
		dstF := dstType.Field(dstIndex)
		bindAttr := dstF.Tag.Get(tag)
		if bindAttr == "" {
			continue
		}
		if unicode.IsLower([]rune(bindAttr)[0]) {
			return fmt.Errorf("tag %s not Export Attr", bindAttr)
		}

		dstTagAttrs[bindAttr] = dstIndex
	}

	for srcIndex := 0; srcIndex < srcLen; srcIndex++ {
		srcF := srcType.Field(srcIndex)
		dstIndex, ok := dstTagAttrs[srcF.Name]
		if !ok {
			continue
		}
		dstF := dstType.Field(dstIndex)
		if dstF.Type != srcF.Type {
			return fmt.Errorf("src[%s(%v)]!=dst[%s(%v)]", srcF.Name, srcF.Type, dstF.Name, dstF.Type)
		}
		if err := AssertType(srcF.Type); err != nil {
			return fmt.Errorf("FiledType[%s(%v),err=%v]", srcF.Name, srcF.Type, err)
		}
		srcFV := srcVal.Field(srcIndex)
		dstFV := dstVal.Field(dstIndex)

		switch dstF.Type.Kind() {
		case reflect.String:
			dstFV.SetString(srcFV.String())
		case reflect.Array:
			n := reflect.Copy(dstFV, srcFV)
			if srcFV.Len() != n {
				return fmt.Errorf("src %s not all copied", srcF.Name)
			}
		case reflect.Slice:
			if !dstFV.CanSet() {
				return fmt.Errorf("dst %s can't be set", dstF.Name)
			}
			dstFV.Set(reflect.MakeSlice(dstF.Type, srcFV.Len(), srcFV.Len()))
			n := reflect.Copy(dstFV, srcFV)
			if srcFV.Len() != n {
				return fmt.Errorf("src %s not all copied", srcF.Name)
			}
		case reflect.Map:
			if !dstFV.CanSet() {
				return fmt.Errorf("dst %s can't be set", dstF.Name)
			}
			if srcFV.IsNil() {
				return fmt.Errorf("src %s map nil", srcF.Name)
			}
			dstFV.Set(reflect.MakeMap(dstF.Type))
			iter := srcFV.MapRange()
			for iter.Next() {
				dstFV.SetMapIndex(iter.Key(), iter.Value())
			}
		case reflect.Float32, reflect.Float64:
			dstFV.SetFloat(srcFV.Float())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			dstFV.SetInt(srcFV.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			dstFV.SetUint(srcFV.Uint())
		case reflect.Bool:
			dstFV.SetBool(srcFV.Bool())
		}
		delete(dstTagAttrs, srcF.Name)
	}

	if len(dstTagAttrs) != 0 {
		return fmt.Errorf("not all tag binded %v", dstTagAttrs)
	}

	return nil
}
