package check

import (
	"syscall"
	"unsafe"
)

func ChkDisk(val float32) bool {
	minDiskSizeGB := float32(50.0)

	if val != 0 {
		minDiskSizeGB = float32(val)
	}

	var kernel32 = syscall.NewLazyDLL("kernel32.dll")
	var getDiskFreeSpaceEx = kernel32.NewProc("GetDiskFreeSpaceExW")

	lpFreeBytesAvailable := int64(0)
	lpTotalNumberOfBytes := int64(0)
	lpTotalNumberOfFreeBytes := int64(0)

	getDiskFreeSpaceEx.Call(
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("C:"))),
		uintptr(unsafe.Pointer(&lpFreeBytesAvailable)),
		uintptr(unsafe.Pointer(&lpTotalNumberOfBytes)),
		uintptr(unsafe.Pointer(&lpTotalNumberOfFreeBytes)))

	diskSizeGB := float32(lpTotalNumberOfBytes) / 1073741824

	if diskSizeGB > minDiskSizeGB {
		return true
	} else {
		return false
	}
}
