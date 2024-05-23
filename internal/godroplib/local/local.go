package local

import (
	"log"
	"strings"

	"github.com/kopp0ut/godropit/pkg/dropfmt"
)

var Droppers = []string{"CreateFiber", "CreateThread", "CreateThreadNative", "EtwpCreateETWThread", "NtQueueAPCThreadExLocal", "goSyscall", "UUIDFromStringA", "BananaPhone", "EnumerateChildWindows", "EnumerateLoadedModules", "CreateThreadPoolWait"}

var Hold = `
	for {

	}
`

func SelectLocal(selected string) (Dlls, Inject, Import, Extra string) {

	_, selected, _ = dropfmt.PromptList(Droppers, "Select the local dropper you would like to use:")

	switch strings.ToLower(selected) {
	case "createfiber":
		Dlls = CreateFiberDlls
		Inject = CreateFiber
		Import = CreateFiberImports
		Extra = ""
	case "createthread":
		Dlls = CreateThreadDlls
		Inject = CreateThread
		Import = CreateThreadImports
		Extra = ""
	case "createthreadnative":
		Dlls = CreateThreadNativeDlls
		Inject = CreateThreadNative
		Import = CreateThreadNativeImports
		Extra = ""
	case "etwpcreateetwthread":
		Dlls = EtwpCreateETWThreadDlls
		Inject = EtwpCreateETWThread
		Import = EtwpCreateETWThreadImports
		Extra = ""
	case "ntqueueapcthreadexlocal":
		Dlls = NtQueueAPCThreadExLocalDlls
		Inject = NtQueueAPCThreadExLocal
		Import = NtQueueAPCThreadExLocalImports
		Extra = NtQueueAPCThreadExLocalExtra
	case "gosyscall":
		Dlls = goSyscallDlls
		Inject = goSyscall
		Import = goSyscallImports
	case "uuidfromstringa":
		Dlls = UUIDFromStringADlls
		Inject = UUIDFromStringA
		Import = UUIDFromStringAImports
		Extra = UUIDFromStringAExtra
	case "bananaphone":
		Dlls = BananaPhoneDlls
		Inject = BananaPhone
		Import = BananaPhoneImports
		Extra = BananaPhoneExtra
	case "enumeratechildwindows":
		Dlls = EnumChildWindowsDlls
		Inject = EnumChildWindows
		Import = EnumChildWindowsImports
		Extra = ""
	case "enumerateloadedmodules":
		Dlls = EnumerateLoadedModulesDlls
		Inject = EnumerateLoadedModules
		Import = EnumerateLoadedModulesImports
		Extra = ""
	case "createthreadpoolwait":
		Dlls = CreateThreadPoolWaitDlls
		Inject = CreateThreadPoolWait
		Import = CreateThreadPoolWaitImports
		Extra = ""
	default:
		log.Fatalf("Error: Method '%s' not found.\nPlease use one of the following local methods: %v\n", selected, Droppers)

	}

	return Dlls, Inject, Import, Extra
}
