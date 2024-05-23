package gengo

import (
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
)

var Debug bool

func GoGetEnv() {
	cmd := exec.Command("go", "env")

	// open the out file for writing
	outfile, err := os.Create("./goenv.txt")
	if err != nil {
		log.Fatalf("Error getting go env: %v\n", err)
	}
	defer outfile.Close()
	cmd.Stdout = outfile

	err = cmd.Start()
	if err != nil {
		log.Fatalf("Error getting go env: %v\n", err)
	}
	cmd.Wait()
}

func ReadEnv(envar string) string {

	data, err := os.ReadFile("./goenv.txt")
	if err != nil {
		color.Red("Could not read goenv.txt, compilation may fail.")
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.Contains(line, envar) {
			envarout := strings.Split(line, "=")
			envarout[1] = strings.Trim(envarout[1], "\"")

			return envarout[1]

		}

	}

	color.Red("Couldn't Find go env %s, compilation may fail.", envar)

	return "failed"

}
