package main

import (
	//"fmt"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

// Login is its on controller because it acts differently than others
type LoginController struct {
	Api        *Api
	OauthToken string
	Username   string
}

var configPath = os.Getenv("HOME") + "/.mongohq"
var credentialFile = configPath + "/credentials"

var Email, OauthToken string

func (c *LoginController) login() error {
	fmt.Println("Enter your MongoHQ credentials.")
	username := prompt("Email")
	password, err := safeGetPass("Password (typing will be hidden): ")

	if err != nil {
		return errors.New("Error returning password.  We may not be compliant with your system yet.  Please send us a message telling us about your system to support@mongohq.com.")
	}

	oauthToken, err := c.Api.Authenticate(username, password, "")

	return c.processAuthenticationResponse(username, password, oauthToken, err)
}

func (c *LoginController) processAuthenticationResponse(username, password, oauthToken string, err error) error {
	if err != nil {
		if err.Error() == "2fa token required" {
			twoFactorToken := prompt("2fa token")
			oauthToken, err := c.Api.Authenticate(username, password, twoFactorToken)
			return c.processAuthenticationResponse(username, password, oauthToken, err)
		} else {
			return err
		}
	} else {
		err = c.storeCredentials(username, oauthToken)

		if err != nil {
			return err
		} else {
			fmt.Println("\nAuthentication complete.\n\n")

			c.Api.OauthToken = oauthToken

			accounts, err := c.Api.GetAccounts()
			if err != nil {
				return errors.New("Error returning accounts after authentication.  Seems like something with authentication may have failed.  Please try again.")
			}

			config := getConfig()
			if len(accounts) == 1 {
				config.AccountSlug = accounts[0].Slug
			} else {
				fmt.Println("== Accounts")
				for _, account := range accounts {
					fmt.Println("  " + account.Slug)
				}
				accountSlug := prompt("Which account should be default? (Can be changed later with config:account)")
				config.AccountSlug = accountSlug
			}
			config.Save()

			return nil
		}
	}
}

func (c *LoginController) storeCredentials(username, oauth string) error {
	credentials := make(map[string]interface{})
	credentials["email"] = username
	credentials["oauth_token"] = oauth

	jsonText, _ := json.Marshal(credentials)

	err := os.MkdirAll(configPath, 0700)

	if err != nil {
		return errors.New("Error creating directory " + configPath)
	}

	err = ioutil.WriteFile(credentialFile, jsonText, 0400)

	if err != nil {
		err = errors.New("Error writing credentials to " + credentialFile)
	}

	return err
}

func (c *LoginController) readCredentialFile() (jsonResponse map[string]interface{}, err error) {
	if _, err := os.Stat(credentialFile); os.IsNotExist(err) { // check if file exists
		return nil, errors.New("Credential file does not exist.")
	} else {
		jsonText, err := ioutil.ReadFile(credentialFile)
		_ = json.Unmarshal(jsonText, &jsonResponse)

		c.Api.OauthToken = jsonResponse["oauth_token"].(string)

		return jsonResponse, err
	}
}

func (l *LoginController) RequireAuth() {
	l.verifyAuth()
}

func (c *LoginController) Logout() {
	_, err := c.Api.restDelete(c.Api.apiUrl("/authorization"))

	os.Remove(credentialFile)
	os.Remove(configFile)

	if err != nil {
		fmt.Println("Error deleting authorization token.  You will need to do that manually from the MongoHQ UI.")
	} else {
		fmt.Println("Logout successful.")
	}
}

func (c *LoginController) verifyAuth() {
	_, err := c.readCredentialFile()
	if err != nil {
		err := c.login()

		if err != nil {
			fmt.Println("\n" + err.Error() + "\n")
			os.Exit(1)
		}
	}

	c.Api.Config = getConfig()
}
