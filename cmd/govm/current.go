package main

import (
	"fmt"
	"github.com/Open-Source-CQUT/govm"
	"github.com/Open-Source-CQUT/govm/pkg/errorx"
	"github.com/spf13/cobra"
)

var currentCmd = &cobra.Command{
	Use:   "current",
	Short: "Show current using Go version",
	RunE: func(cmd *cobra.Command, args []string) error {
		current, err := RunCurrent()
		if err != nil {
			return err
		}
		fmt.Println(current)
		return nil
	},
}

func RunCurrent() (string, error) {
	usingVersion, err := govm.GetUsingVersion()
	if err != nil {
		return "", err
	} else if usingVersion == "" {
		return "", errorx.Warn("no using version")
	}
	return usingVersion, nil
}
