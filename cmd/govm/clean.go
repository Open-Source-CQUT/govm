package main

import (
	"github.com/Open-Source-CQUT/govm"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "clean local cache",
	RunE: func(cmd *cobra.Command, args []string) error {
		return RunClean()
	},
}

func RunClean() error {
	cacheDir, err := govm.GetCacheDir()
	if err != nil {
		return err
	}
	entries, err := os.ReadDir(cacheDir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		name := filepath.Join(cacheDir, entry.Name())
		err := os.RemoveAll(name)
		if err != nil {
			return err
		}
	}
	return nil
}
