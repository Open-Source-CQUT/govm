package main

import (
	"errors"
	"fmt"
	"github.com/Open-Source-CQUT/govm"
	"github.com/spf13/cobra"
	"slices"
)

var useCmd = &cobra.Command{
	Use:   "use",
	Short: "use specified version of go",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("not version specified")
		} else {
			version := args[0]
			return RunUse(version)
		}
	},
}

func RunUse(version string) error {
	checkV, err := govm.CheckVersion(version)
	if err != nil {
		return err
	}
	version = checkV
	// try to find from local
	localList, err := govm.LocalList(false)
	if err != nil {
		return err
	}

	if !slices.Contains(localList, version) {
		return fmt.Errorf("%s not found in local", version)
	}

	storeData, err := govm.ReadStore()
	if err != nil {
		return err
	}
	storeData.Use = version
	if err := govm.WriteStore(storeData); err != nil {
		return err
	}

	//TODO update profile
	govm.Printf("using version %s now", version)
	return nil
}
