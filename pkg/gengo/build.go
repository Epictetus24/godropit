package gengo

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/fatih/color"
)

var Env []string

var Arch string

var Garble bool

var compilerBin string

func init() {
	GoGetEnv()

	Env = []string{
		fmt.Sprintf("CC=%s", ReadEnv("CC")),
		fmt.Sprintf("CGO_ENABLED=%s", "0"),
		//fmt.Sprintf("GOCACHE=%s", ReadEnv("GOCACHE")),
		//fmt.Sprintf("GOMODCACHE=%s", ReadEnv("GOMODCACHE")),
		fmt.Sprintf("GOPRIVATE=%s", ReadEnv("GOPRIVATE")),
		fmt.Sprintf("PATH=%s:%s", path.Join(ReadEnv("GOVERSION"), "bin"), os.Getenv("PATH")),
		//fmt.Sprintf("GOPATH=%s", ReadEnv("GOPATH")),
		//fmt.Sprintf("GOROOT=%s", ReadEnv("GOROOT")),
		fmt.Sprintf("HOME=%s", os.Getenv("HOME")),
	}
}

func buildInstruct(outdir, fname string, dll bool, x86 bool) error {
	if Garble {
		compilerBin = "garble"
	} else {
		compilerBin = "go"
	}

	var err error

	if err != nil {
		log.Fatalf("Error changing to output dir %v.\n", err)
	}

	fmt.Printf("Shellcode files are now in %s.\n", outdir)

	if err != nil {
		log.Fatalf("Error getting working dir %v.\n", err)
	}

	fmt.Printf("Getting go env and saving to %s\n", outdir+"/goenv.txt")
	GoGetEnv()
	os.Remove("go.mod")
	os.Remove("go.sum")
	if err != nil {
		log.Printf("Tried to clean any previous go.mod files but failed:\n%v.\n", err)
	}

	//go mod init the directory

	Arch = "amd64"
	if x86 {
		Arch = "386"
	}

	Env = append(Env, fmt.Sprintf("GOARCH=%s", Arch))
	Env = append(Env, fmt.Sprintf("GOOS=%s", "windows"))
	Env = append(Env, fmt.Sprintf("GOTMPDIR=%s", outdir))

	var out bytes.Buffer
	var stderr bytes.Buffer

	//Warn about usage if garble is chosen over go compiler.
	if Garble {
		fmt.Println("Using garble to compile. Ensure garble is in your PATH and is the correct version for your go installation.")
	} else {
		fmt.Println("Using go for executable. Ensure go is in your PATH.")
	}
	initCmd := exec.Command("go", "mod", "init", strings.ReplaceAll(fname, ".go", ""))
	initCmd.Env = Env
	//initCmd.Dir = outdir

	initCmd.Stdout = &out
	initCmd.Stderr = &stderr
	err = initCmd.Run()
	if err != nil {
		color.Red("[gengo] Woops, something went wrong with go mod init, soz.\n")
		color.Red(fmt.Sprint(err) + ": " + stderr.String())
		return err
	}

	initCmd.Wait()
	//go mod tidy to get depencies
	tidyCmd := exec.Command("go", "mod", "tidy")
	tidyCmd.Env = Env
	tidyCmd.Stdout = &out
	tidyCmd.Stderr = &stderr
	//tidyCmd.Dir = outdir
	err = tidyCmd.Run()
	if err != nil {
		color.Red("[gengo] Woops, something went wrong with go mod tidy:\n")
		color.Red(fmt.Sprint(err) + ": " + stderr.String())
		return err
	}

	tidyCmd.Wait()

	getWinmod := exec.Command("go", "get", "golang.org/x/sys/windows")
	getWinmod.Env = Env
	getWinmod.Stdout = &out
	getWinmod.Stderr = &stderr
	err = getWinmod.Run()
	if err != nil {
		color.Red("[gengo] Woops, something went wrong with go mod tidy:\n")
		color.Red(fmt.Sprint(err) + ": " + stderr.String())
		return err
	}

	color.Green("Prep done! go.mod and go.sum may need updating before compilation.")

	return nil

}

// Compiles the resulting go file into the appropriate format.
func buildFileGo(outdir, fname string, dll bool, x86 bool) (bool, error) {

	var command []string
	var err error

	if Garble {
		compilerBin = "garble"
	} else {
		compilerBin = "go"
	}

	if err != nil {
		log.Fatalf("Error changing to output dir %v.\n", err)
	}

	if err != nil {
		log.Fatalf("Error getting working dir %v.\n", err)
	}

	var out bytes.Buffer
	var stderr bytes.Buffer
	outname := fmt.Sprintf("%s_%s", strings.ReplaceAll(fname, ".go", ""), Arch)

	//modify args based on compiler and output format:
	if Garble {
		if dll {

			outname = outname + ".dll"
			command = []string{"-tiny", "-literals", "-seed=random", "build", "-o", outname, "-trimpath", "-buildmode=c-shared", `-ldflags=-w -s -H=windowsgui`, fname}
		} else {
			outname = outname + ".exe"
			command = []string{"-tiny", "-literals", "-seed=random", "build", "-o", outname, "-trimpath", `-ldflags=-w -s -H=windowsgui`, fname}
		}

	} else {
		if dll {

			outname = outname + ".dll"
			command = []string{"build", "-o", outname, "-trimpath", "-buildmode=c-shared", `-ldflags=-w -s -H=windowsgui`, fname}
		} else {
			outname = outname + ".exe"
			command = []string{"build", "-o", outname, "-trimpath", `-ldflags=-w -s -H=windowsgui`, fname}
		}
	}

	cmd := exec.Command(compilerBin, command...)
	cmd.Env = Env
	cmd.Dir = outdir
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if dll {
		cmd.Env[0] = fmt.Sprintf("CC=%s", "x86_64-w64-mingw32-gcc")
		cmd.Env[1] = fmt.Sprintf("CGO_ENABLED=%s", "1")
	}
	if Debug {
		PrintBuild(cmd.Env[0], cmd.Env[1], compilerBin, command, outdir)
	}

	err = cmd.Run()
	if err != nil {
		color.Red("[gengo] Woops, something went wrong with compiling, soz.\n")
		color.Red(fmt.Sprint(err) + ": " + stderr.String())
		return false, err
	}

	cmd.Wait()
	if Garble {
		color.Green("Dropper compiled with Garble compiler, find it at %s/%s.\n", outdir, outname)
	} else {
		color.Green("Dropper compiled with regular go compiler, find it at %s/%s.\n", outdir, outname)
	}

	return true, nil

}

func PrintBuild(cc, cgo, compileBin string, args []string, outdir string) {

	buildcmd := fmt.Sprintf("%s %s %s %s\n", cc, cgo, compileBin, strings.Join(args[:], " "))

	color.Cyan("To compile yourself, navigate to %s and run:\n\n", outdir)
	fmt.Printf("%s\n\n", buildcmd)
	if !Garble {
		color.Cyan("Or if you're really cool use garble: https://github.com/burrowers/garble")
	}

}
func buildWasm(outdir, fname string) error {
	Arch = "wasm"

	Env = append(Env, fmt.Sprintf("GOARCH=%s", Arch))
	Env = append(Env, fmt.Sprintf("GOOS=%s", "js"))

	var command []string
	var err error

	if err != nil {
		log.Fatalf("Error changing to output dir %v.\n", err)
	}

	fmt.Printf("Shellcode files are now in %s.\n", outdir)
	fmt.Println("Using go for executable. Ensure go is in your PATH.")
	goBinPath := "go"

	if err != nil {
		log.Fatalf("Error getting working dir %v.\n", err)
	}

	fmt.Printf("Getting go env and saving to %s\n", outdir+"/goenv.txt")
	GoGetEnv()
	os.Remove("go.mod")
	os.Remove("go.sum")
	if err != nil {
		log.Printf("Tried to clean any previous go.mod files but failed:\n%v.\n", err)
	}

	//go mod init the directory

	var out bytes.Buffer
	var stderr bytes.Buffer

	outname := fmt.Sprintf("%s_%s", strings.ReplaceAll(fname, ".go", ""), Arch)

	initCmd := exec.Command(goBinPath, "mod", "init", strings.ReplaceAll(fname, ".go", ""))
	initCmd.Env = Env

	initCmd.Stdout = &out
	initCmd.Stderr = &stderr
	initCmd.Dir = outdir
	err = initCmd.Run()
	if err != nil {
		color.Red("[gengo] Woops, something went wrong with go mod init, soz.\n")
		color.Red(fmt.Sprint(err) + ": " + stderr.String())
		return err
	}

	initCmd.Wait()
	//go mod tidy to get depencies
	tidyCmd := exec.Command(goBinPath, "mod", "tidy")
	tidyCmd.Env = Env
	tidyCmd.Stdout = &out
	tidyCmd.Stderr = &stderr
	tidyCmd.Dir = outdir
	err = tidyCmd.Run()
	if err != nil {
		color.Red("[gengo] Woops, something went wrong with go mod tidy:\n")
		color.Red(fmt.Sprint(err) + ": " + stderr.String())
		return err
	}

	tidyCmd.Wait()

	if err != nil {
		log.Fatalf("Error changing to output dir %v.\n", err)
	}

	if err != nil {
		log.Fatalf("Error getting working dir %v.\n", err)
	}

	outname = fmt.Sprintf("%s_%s", strings.ReplaceAll(fname, ".go", ""), Arch)

	outname = outname + ".wasm"
	command = []string{"build", "-o", outname, "-trimpath", `-ldflags=-w -s -H=windowsgui`, fname}

	cmd := exec.Command(goBinPath, command...)
	cmd.Env = Env
	cmd.Dir = outdir
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		color.Red("[gengo] Woops, something went wrong with compiling, soz.\n")
		color.Red(fmt.Sprint(err) + ": " + stderr.String())
		return err
	}

	cmd.Wait()

	color.Green("Dropper compiled with regular normal go compiler, find it at %s/%s.\n", outdir, outname)

	return nil

}
