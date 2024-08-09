package main

import (
	"github.com/Open-Source-CQUT/govm"
	"github.com/spf13/cobra"
	"runtime"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show govm version",
	RunE: func(cmd *cobra.Command, args []string) error {
		govm.Printf("govm versoin %s %s/%s\n", Version, runtime.GOOS, runtime.GOARCH)
		return nil
	},
}
