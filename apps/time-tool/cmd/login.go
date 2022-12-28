package cmd

import (
	"errors"
	"finlab/apps/time-tool/api"
	"finlab/apps/time-tool/config"
	"finlab/apps/time-tool/core"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// LoginCmd provides the command for logging into the finlab server.
var LoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate to the finlab server.",
	Run: func(cmd *cobra.Command, args []string) {
		loginPrompt := promptui.Prompt{
			Label:    "Login",
			Validate: validateLogin,
		}

		passwordPrompt := promptui.Prompt{
			Label:    "Password",
			Validate: validatePassword,
			Mask:     '*',
		}

		login, errL := loginPrompt.Run()
		password, errP := passwordPrompt.Run()
		promptErr := core.NotNil(errL, errP)

		if promptErr != nil {
			core.Danger("Prompt failed: %v\n", promptErr)
			return
		}

		token, err := api.SignIn(login, password)
		if err != nil {
			core.Danger("Error: %v\n", err.Error())
			return
		}
		config.SaveToken(token)
		core.Success("You successfully authenticated!\n")
	},
}

func validateLogin(input string) error {
	if len(input) < 3 {
		return errors.New("LOGIN MUST HAVE MORE THAN 3 CHARACTERS")
	}
	return nil
}

func validatePassword(input string) error {
	if len(input) == 0 {
		return errors.New("PASSWORD MUST HAVE MORE THAN 1 CHARACTERS")
	}
	return nil
}
