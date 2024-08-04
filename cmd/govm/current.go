package main

import (
	"fmt"
	"github.com/Open-Source-CQUT/govm"
	"github.com/Open-Source-CQUT/govm/pkg/errorx"
	"github.com/spf13/cobra"
)

var currentCmd = &cobra.Command{
	Use:   "current",
	Short: "show current using version",
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
	store, err := govm.ReadStore()
	if err != nil {
		return "", err
	}
	if store.Use == "" {
		return "", errorx.Warn("no using version")
	}
	return store.Use, nil
}
