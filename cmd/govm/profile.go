package main

import (
	"fmt"
	"github.com/Open-Source-CQUT/govm"
	"github.com/spf13/cobra"
	"os"
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
	profilepath, err := govm.GetProfile()
	if err != nil {
		return "", err
	}
	profile, err := os.ReadFile(profilepath)
	if err != nil {
		return "", err
	}
	return govm.Bytes2string(profile), nil
}
