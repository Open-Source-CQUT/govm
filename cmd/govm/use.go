package main

import (
	"github.com/Open-Source-CQUT/govm"
	"github.com/Open-Source-CQUT/govm/pkg/errorx"
	"github.com/spf13/cobra"
	"slices"
)

var useCmd = &cobra.Command{
	Use:   "use",
	Short: "Use specified Go version",
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
	config, err := govm.ReadConfig()
	if err != nil {
		return err
	}
	config.Use = using.Version
	if err := govm.WriteConfig(config); err != nil {
		return err
	}

	govm.Tipf("Use %s now", version)
	return nil
}
