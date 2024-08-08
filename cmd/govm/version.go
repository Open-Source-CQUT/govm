package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"runtime"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show govm version",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("govm versoin %s %s/%s\n", Version, runtime.GOOS, runtime.GOARCH)
		return nil
	},
}
