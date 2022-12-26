package cmd

import (
	"errors"
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// LoginCmd provides the command for logging into the FA server using the Device Flow.
var LoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate to the FA server using the OAuth Device Flow.",
	Run: func(cmd *cobra.Command, args []string) {
		templates := &promptui.PromptTemplates{
			Prompt:  "{{ . }} ",
			Valid:   "{{ . | green }} ",
			Invalid: "{{ . | red }} ",
			Success: "{{ . | bold }} ",
		}

		validateLogin := func(input string) error {
			if len(input) < 3 {
				return errors.New("Login must have more than 3 characters")
			}
			return nil
		}

		validatePassword := func(input string) error {
			if len(input) < 0 {
				return errors.New("Password must have more than 3 characters")
			}
			return nil
		}

		loginPrompt := promptui.Prompt{
			Label:     "Login",
			Validate:  validateLogin,
			Templates: templates,
		}

		passwordPrompt := promptui.Prompt{
			Label:     "Password",
			Validate:  validatePassword,
			Mask:      '*',
			Templates: templates,
		}

		login, err := loginPrompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}
		password, err := passwordPrompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		token, err := SignIn(login, password)
		if err != nil {
			Danger.Printf(err.Error())
			return
		}
		SaveToken(token)
		Success.Printf("You successfully authenticated!")
	},
}
