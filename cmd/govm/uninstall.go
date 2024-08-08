package main

import (
	"github.com/Open-Source-CQUT/govm"
	"github.com/Open-Source-CQUT/govm/pkg/errorx"
	"github.com/spf13/cobra"
	"os"
)

var uninstallCmd = &cobra.Command{
	Use:     "uninstall",
	Short:   "Uninstall specified Go version",
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

	store, err := govm.ReadStore()
	if err != nil {
		return err
	}
	foundV, exist := store.Versions[version]
	if !exist {
		return errorx.Warnf("%s not installed", v)
	}

	if store.Use == foundV.Version {
		store.Use = ""
	}
	err = govm.WriteStore(store)
	if err != nil {
		return err
	}

	// remove from store
	if err := os.RemoveAll(foundV.Path); err != nil {
		return err
	}
	govm.Tipf("Version %s uninstalled", v)

	return nil
}
