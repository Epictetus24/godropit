package dropfmt

import (
	"github.com/manifoldco/promptui"
)

func PromptList(list []string, prompt string) (int, string, error) {
	var err error

	payloadchoice := promptui.Select{
		Label: prompt,
		Items: list,
	}

	choiceInt, choiceStr, err := payloadchoice.Run()
	if err != nil {
		return 0, "", err
	}
	return choiceInt, choiceStr, nil
}
