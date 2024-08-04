package main

import (
	"github.com/Open-Source-CQUT/govm"
	"github.com/Open-Source-CQUT/govm/pkg/errorx"
	"github.com/spf13/cobra"
	"os"
	"slices"
)

var uninstallCmd = &cobra.Command{
	Use:     "uninstall",
	Short:   "uninstall specified Go version",
	Aliases: []string{"ui"},
	RunE: func(cmd *cobra.Command, args []string) error {
		var version string
		if len(args) > 0 {
			version = args[0]
		}
		return RunUninstall(version)
	},
}

func RunUninstall(v string) error {
	if v == "" {
		return errorx.Warn("no version specified")
	}

	// check versions
	version, ok := govm.CheckVersion(v)
	if !ok {
		return errorx.Warnf("invliad version %s", v)
	}

	// check if is already installed
	locals, err := govm.GetLocalVersions(false)
	installed := slices.ContainsFunc(locals, func(v govm.Version) bool {
		return v.Version == version
	})
	if !installed {
		return errorx.Warnf("%s not installed", v)
	}

	storeData, err := govm.ReadStore()
	if err != nil {
		return err
	}
	removedPath := storeData.Versions[version].Path
	delete(storeData.Versions, version)
	err = govm.WriteStore(storeData)
	if err != nil {
		return err
	}

	// remove from store
	if err := os.RemoveAll(removedPath); err != nil {
		return err
	}
	govm.Tipf("%s uninstalled", v)

	return nil
}
