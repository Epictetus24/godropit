package child

import (
	"log"
	"strings"

	"github.com/kopp0ut/godropit/pkg/dropfmt"
)

var Droppers = []string{"CreateProcess", "CreateProcessWithPipe", "EarlyBird"}

func SelectChild(selected string) (Dlls, Inject, Import string) {

	if selected == "" {
		_, selected, _ = dropfmt.PromptList(Droppers, "Select the child dropper you would like to use:")
	}

	switch strings.ToLower(selected) {
	case "createprocess":
		Dlls = CreateProcessDlls
		Inject = CreateProcess
		Import = CreateProcessImports
	case "createprocesswithpipe":
		Dlls = CreateProcWithPipeDlls
		Inject = CreateProcWithPipe
		Import = CreateProcWithPipeImports
	case "earlybird":
		Dlls = EarlyBirdDlls
		Inject = EarlyBird
		Import = EarlyBirdImports
	default:
		log.Fatalf("Error: Method '%s' not found.\nPlease use one of the following child methods: %v\n", selected, Droppers)

	}

	return Dlls, Inject, Import
}
