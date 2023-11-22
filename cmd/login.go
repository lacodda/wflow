package cmd

import (
	"wflow/api"
	"wflow/config"
	"wflow/core"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

// LoginCmd provides the command for logging into the finlab server.
var LoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate to the finlab server.",
	Run: func(cmd *cobra.Command, args []string) {
		credentials := core.Credentials{}

		err := survey.Ask(getLoginQuestions(core.Credentials{}), &credentials)
		if err != nil {
			core.Danger("Prompt failed: %v\n", err.Error())
			return
		}

		token, err := api.SignIn(credentials)
		if err != nil {
			core.Danger("Error: %v\n", err.Error())
			return
		}
		config.SaveToken(token)
		core.Success("You successfully authenticated!\n")
	},
}

func getLoginQuestions(credentials core.Credentials) []*survey.Question {
	return []*survey.Question{
		{
			Name:     "email",
			Prompt:   &survey.Input{Message: "Login", Default: credentials.Email},
			Validate: survey.Required,
		},
		{
			Name:     "password",
			Prompt:   &survey.Password{Message: "Password"},
			Validate: survey.Required,
		},
	}
}
