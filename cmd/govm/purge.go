package main

import (
	"github.com/Open-Source-CQUT/govm"
	"github.com/spf13/cobra"
	"os"
)

var purgeCmd = &cobra.Command{
	Use:   "purge",
	Short: "Purge all the cache and local installed versions",
	RunE: func(cmd *cobra.Command, args []string) error {
		return RunClean()
	},
}

func RunClean() error {
	store, err := govm.ReadStore()
	if err != nil {
		return err
	}
	storeDir, err := govm.GetInstallation()
	if err != nil {
		return err
	}
	err = os.RemoveAll(storeDir)
	if err != nil {
		return err
	}
	govm.Tipf("Purged versions %d", len(store.Versions))
	return nil
}
