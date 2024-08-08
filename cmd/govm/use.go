package main

import (
	"github.com/Open-Source-CQUT/govm"
	"github.com/Open-Source-CQUT/govm/pkg/errorx"
	"github.com/spf13/cobra"
	"os"
	"slices"
)

var useCmd = &cobra.Command{
	Use:   "use",
	Short: "use specified Go version",
	RunE: func(cmd *cobra.Command, args []string) error {
		var version string
		if len(args) > 0 {
			version = args[0]
		}
		return RunUse(version)
	},
}

func RunUse(v string) error {
	// check version
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
	i := slices.IndexFunc(localList, func(v govm.Version) bool {
		return v.Version == version
	})
	if i == -1 {
		return errorx.Warnf(`%s not found in local, use comamnd "govm install %s" to install`, version, version)
	}
	using := localList[i]

	// update store.toml
	storeData, err := govm.ReadStore()
	if err != nil {
		return err
	}
	storeData.Use = using.Version
	if err := govm.WriteStore(storeData); err != nil {
		return err
	}

	// update symlink
	currentLink, err := govm.GetRootSymLink()
	if err != nil {
		return err
	}
	os.Remove(currentLink)
	err = os.Symlink(using.Path, currentLink)
	if err != nil {
		return err
	}

	govm.Tipf("Use %s now", version)
	return nil
}
