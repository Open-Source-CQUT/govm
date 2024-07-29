package main

import (
	"github.com/Open-Source-CQUT/govm"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var uninstallCmd = &cobra.Command{
	Use:     "uninstall",
	Short:   "uninstall specified version of go",
	Aliases: []string{"rm"},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			govm.Println("no version specified")
			return nil
		}
		for _, version := range args {
			err := RunUninstall(version)
			if err != nil {
				return err
			}
		}
		return nil
	},
}

func RunUninstall(version string) error {
	if !strings.HasPrefix(version, "go") {
		version = "go" + version
	}
	storeData, err := govm.ReadStore()
	if err != nil {
		return err
	}
	loc, found := storeData.Root[version]
	if !found {
		govm.Printf("%s not found in store", version)
		return nil
	}
	if err := os.RemoveAll(loc); err != nil {
		return err
	}
	delete(storeData.Root, version)
	return govm.WriteStore(storeData)
}
