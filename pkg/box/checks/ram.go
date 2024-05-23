package check

import (
	"syscall"
	"unsafe"
)

type MEMORYSTATUSEX struct {
	dwLength                uint32
	dwMemoryLoad            uint32
	ullTotalPhys            uint64
	ullAvailPhys            uint64
	ullTotalPageFile        uint64
	ullAvailPageFile        uint64
	ullTotalVirtual         uint64
	ullAvailVirtual         uint64
	ullAvailExtendedVirtual uint64
}

func ChkRam() bool {
	var kernel32 = syscall.NewLazyDLL("kernel32.dll")
	var globalMemoryStatusEx = kernel32.NewProc("GlobalMemoryStatusEx")

	var memInfo MEMORYSTATUSEX
	memInfo.dwLength = uint32(unsafe.Sizeof(memInfo))
	globalMemoryStatusEx.Call(uintptr(unsafe.Pointer(&memInfo)))

	if memInfo.ullTotalPhys/1073741824 > 1 {
		return true
	} else {
		return false
	}
}
