package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "show profile scripts",
	RunE: func(cmd *cobra.Command, args []string) error {
		script, err := RunProfile()
		if err != nil {
			return err
		}
		fmt.Println(script)
		return nil
	},
}

func RunProfile() (string, error) {
	profile, err := useProfile()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("source %s", profile), nil
}
