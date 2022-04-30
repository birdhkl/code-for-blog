package go_reflect

import (
	"fmt"
	"reflect"
)

// 源码参考 http://r12f.com/posts/learning-golang-object-model-inbox-data-type/#more
var tableHeader = fmt.Sprintf("%-12s%-30s%-20s%-10s %-11s %-4s %-10s\n", "Var", "Type", "Address", "RootOffset", "LocalOffset", "Size", "Value")

func DumpObjectWithTableHeader(name string, p reflect.Value) (string, error) {
	obj, err := DumpObject(name, p)
	if err != nil {
		return "", err
	}
	return tableHeader + obj, nil
}

func DumpObject(name string, p reflect.Value) (string, error) {
	if p.Kind() == reflect.Interface || p.Kind() == reflect.Ptr {
		p = p.Elem()
		return dumpObject(name, p, p.UnsafeAddr(), p.UnsafeAddr()), nil
	}

	return "", fmt.Errorf("%v not support UnsafeAddr", p.Type().Name())
}

func dumpObject(path string, v reflect.Value, rootBaseAddr uintptr, localBaseAddr uintptr) string {
	detail := dumpObjectDetail(path, v, rootBaseAddr, localBaseAddr)

	switch v.Kind() {
	case reflect.Struct:
		childLocalBaseAddr := v.UnsafeAddr()
		for i := 0; i < v.NumField(); i++ {
			fieldPath := fmt.Sprintf("%s.%s", path, v.Type().Field(i).Name)
			detail += dumpObject(fieldPath, v.Field(i), rootBaseAddr, childLocalBaseAddr)
		}
	}
	return detail
}

func dumpObjectDetail(path string, v reflect.Value, rootBaseAddr uintptr, localBaseAddr uintptr) string {
	addr := v.UnsafeAddr()
	var val interface{} = "#unexported#"
	if v.CanInterface() {
		val = v.Interface()
	}
	return fmt.Sprintf("%-12s%-30s0x%018x%10v %11v %4v %10v\n", path, v.Type().String(), addr,
		addr-rootBaseAddr, addr-localBaseAddr, v.Type().Size(), val)

}
