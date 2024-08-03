package main

import (
	"github.com/Open-Source-CQUT/govm"
	"github.com/Open-Source-CQUT/govm/pkg/errorx"
	"github.com/spf13/cobra"
	"slices"
)

var useCmd = &cobra.Command{
	Use:   "use",
	Short: "use specified version of go",
	RunE: func(cmd *cobra.Command, args []string) error {
		var version string
		if len(args) > 0 {
			version = args[0]
		}
		return RunUse(version)
	},
}

func RunUse(v string) error {
	if v == "" {
		return errorx.Warn("no version specified")
	}
	version, ok := govm.CheckVersion(v)
	if !ok {
		return errorx.Warnf("invalid version: %s", v)
	}
	// try to find from local
	localList, err := govm.GetLocalVersions(false)
	if err != nil {
		return err
	}

	if !slices.ContainsFunc(localList, func(v govm.Version) bool {
		return v.Version == version
	}) {
		return errorx.Warnf("%s not found in local", version)
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
	govm.Tipf("Use %s now", version)
	return nil
}
