package go_reflect

import "syscall"

func modifyBinary(dstAddr uintptr, jmpCode []byte) {
	page := GetPage(dstAddr)

	if err := syscall.Mprotect(page, syscall.PROT_READ|syscall.PROT_WRITE|syscall.PROT_EXEC); err != nil {
		panic(err)
	}
	defer func() {
		if err := syscall.Mprotect(page, syscall.PROT_READ|syscall.PROT_EXEC); err != nil {
			panic(err)
		}
	}()
	codeSegment := GetEntryCodeSegment(dstAddr)
	copy(codeSegment, jmpCode)
}
