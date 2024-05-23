package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "godropit",
	Short: "Quick and low equity dropper generator.",
	Long: `Use to generate golang droppers with encrypted shellcode.
	For example a new child process dropper could be created with:

	./godropit new child -i <shellcode.bin|evil.exe>
	
	`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}
var debug bool

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	newCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Show debug information and keep source files after compilation.")
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.godropit.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
}
