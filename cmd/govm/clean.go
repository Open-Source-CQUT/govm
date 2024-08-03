package main

import (
	"github.com/Open-Source-CQUT/govm"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var cleanCmd = &cobra.Command{
	Use:     "clean",
	Short:   "clean local cache and redundant versions",
	Aliases: []string{"cl"},
	RunE: func(cmd *cobra.Command, args []string) error {
		return RunClean()
	},
}

func RunClean() error {
	// clean cache
	cacheDir, err := govm.GetCacheDir()
	if err != nil {
		return err
	}
	entries, err := os.ReadDir(cacheDir)
	if err != nil {
		return err
	}

	var cacheCount int
	for _, entry := range entries {
		name := filepath.Join(cacheDir, entry.Name())
		err := os.RemoveAll(name)
		if err != nil {
			return err
		}
	}
	govm.Tipf("Clean cache files %d", cacheCount)

	// remove redundant versions from store
	storeData, err := govm.ReadStore()
	if err != nil {
		return err
	}
	rootStore, err := govm.GetRootStore()
	if err != nil {
		return err
	}
	dirEntries, err := os.ReadDir(rootStore)
	if err != nil {
		return err
	}
	var redundantCount int
	for _, entry := range dirEntries {
		// it is a redundant version if not in store
		if _, ok := storeData.Root[entry.Name()]; !ok {
			err := os.RemoveAll(filepath.Join(rootStore, entry.Name()))
			if err != nil {
				return err

			}
			redundantCount++
		}

	}
	govm.Tipf("Clean redundant files %d", redundantCount)

	return nil
}
