package main

import (
	"fmt"
	"os"
)

func (c *Controller) SetConfigAccount(slug string) {
	account, err := c.Api.GetAccount(slug)

	if err != nil {
		fmt.Println("Error retreiving account: " + err.Error())
		os.Exit(1)
	}

	config := getConfig()
	config.AccountSlug = account.Slug

	if err := config.Save(); err != nil {
		fmt.Println("Error setting default account: " + err.Error())
		os.Exit(1)
	}

	fmt.Println("Set default account to " + account.Slug)
}