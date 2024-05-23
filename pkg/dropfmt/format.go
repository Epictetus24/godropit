package dropfmt

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/fatih/color"
)

//Format

func FormatSC(shellcodepath, outpath, name string, aes, c, cs, golang, b64 bool) {

	scFile, err := os.Open(shellcodepath)
	if err != nil {
		color.Red("Error getting shellcode: %v\n", err)
		log.Fatalln("Exiting...")
	}

	_, err = ioutil.ReadAll(scFile)
	if err != nil {
		color.Red("Error reading shellcode: %v\n", err)
		log.Fatalln("Exiting...")
	}

}

func WriteOutfile(buffer []byte, outpath, name, ext string) error {

	return nil
}

func MakeCByteArray([]byte) string {
	var array string

	return array
}

func MakeCSByteArray([]byte) string {
	var array string

	return array
}
