package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

type AccessToken struct {
	AccessToken string `json:"access_token"`
}

const (
	host string = "http://localhost:3333/api/auth/login"
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

		body := signIn(login, password)
		saveAccessToken(body)
	},
}

func signIn(login string, password string) []byte {
	data := url.Values{
		"email":    {login},
		"password": {password},
	}

	resp, err := http.PostForm(host, data)

	if err != nil {
		log.Fatal(err)
	}

	body, _ := ioutil.ReadAll(resp.Body)

	return body
}

func saveAccessToken(body []byte) {
	Save("access-token.json", getAccessToken(body))

	cyan := color.New(color.FgCyan)
	cyan.Printf("You successfully authenticated!")
}

func getAccessToken(body []byte) AccessToken {
	result := AccessToken{}
	json.Unmarshal([]byte(body), &result)

	return result
}

func PrettyPrint(b []byte) []byte {
	var out bytes.Buffer
	err := json.Indent(&out, b, "", "  ")

	if err != nil {
		log.Fatalln(err)
	}

	return out.Bytes()
}
