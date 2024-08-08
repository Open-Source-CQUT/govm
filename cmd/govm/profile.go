package main

import (
	"fmt"
	"github.com/Open-Source-CQUT/govm"
	"github.com/spf13/cobra"
)

var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "show profile env",
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
	return govm.GetProfileContent()
}
