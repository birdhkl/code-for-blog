// +build amd64

package go_reflect

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
	"syscall"
	"unsafe"
)

// 字节转换成整形
func LittleEndianBytesToXs(b []byte, x interface{}) {
	bytesBuffer := bytes.NewBuffer(b)
	if err := binary.Read(bytesBuffer, binary.LittleEndian, x); err != nil {
		panic(err)
	}
}

// DoMonkeyPatch 猴子补丁
func DoMonkeyPatch(dst interface{}, src interface{}) {
	if reflect.ValueOf(dst).Kind() != reflect.Func {
		panic(fmt.Sprintf("src %s not func", reflect.ValueOf(src).Type().Name()))
	}
	if reflect.ValueOf(dst).Type() != reflect.ValueOf(src).Type() {
		panic(fmt.Sprintf("src %s not target %s", reflect.ValueOf(src).Type().Name(),
			reflect.ValueOf(dst).Type().Name()))
	}
	jmpCode := buildJmpDirective(src)
	dstAddr := reflect.ValueOf(dst).Pointer()
	nextPage := GetPageStartAddress(dstAddr) + uintptr(syscall.Getpagesize())
	for len(jmpCode) > 0 {
		// 页不够大，则截取代码段
		codeLen := len(jmpCode)
		if nextPage-dstAddr < uintptr(len(jmpCode)) {
			codeLen = int(nextPage - dstAddr)
		}
		modifyBinary(dstAddr, jmpCode[:codeLen])
		// 更新地址和下一页地址
		jmpCode = jmpCode[codeLen:]
		dstAddr = nextPage
		nextPage = nextPage + uintptr(syscall.Getpagesize())
	}
}

// GetPageStartAddresss 获取地址的页面地址
func GetPageStartAddress(addr uintptr) uintptr {
	return addr & ^(uintptr(syscall.Getpagesize()) - 1)
}

// 获取页面
func GetPage(addr uintptr) []byte {
	pageAddr := GetPageStartAddress(addr)
	page := reflect.SliceHeader{
		Data: pageAddr,
		Len:  syscall.Getpagesize(),
		Cap:  syscall.Getpagesize(),
	}
	return *(*[]byte)(unsafe.Pointer(&page))
}

// GetEntryCodeSegment 获取代码实际地址
func GetEntryCodeSegment(addr uintptr) []byte {
	page := GetPage(addr)
	startAddr := GetPageStartAddress(addr)
	offset := addr - startAddr
	return page[offset:]
}

// reflect.Value的镜像
type valueMirror struct {
	_    unsafe.Pointer // rtype
	data unsafe.Pointer
	_    uintptr // flag
}

func GetFuncPtr(f interface{}) unsafe.Pointer {
	if reflect.TypeOf(f).Kind() != reflect.Func {
		panic("f not func")
	}
	value := reflect.ValueOf(f)
	funcPtr := (*(*valueMirror)(unsafe.Pointer(&value))).data
	return funcPtr
}

// https://github.com/agiledragon/gomonkey/blob/master/jmp_amd64.go
func buildJmpDirective(f interface{}) []byte {
	funcPtr := GetFuncPtr(f)
	d0 := byte(uintptr(funcPtr))
	d1 := byte(uintptr(funcPtr) >> 8)
	d2 := byte(uintptr(funcPtr) >> 16)
	d3 := byte(uintptr(funcPtr) >> 24)
	d4 := byte(uintptr(funcPtr) >> 32)
	d5 := byte(uintptr(funcPtr) >> 40)
	d6 := byte(uintptr(funcPtr) >> 48)
	d7 := byte(uintptr(funcPtr) >> 56)

	return []byte{
		0x48, 0xBA, d0, d1, d2, d3, d4, d5, d6, d7, // MOV rdx, double
		0xFF, 0x22, // JMP [rdx]
	}
}
