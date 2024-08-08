package main

import (
	"fmt"
	"github.com/Open-Source-CQUT/govm"
	"github.com/Open-Source-CQUT/govm/pkg/errorx"
	"github.com/spf13/cobra"
	"path/filepath"
)

var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Show profile env",
	RunE: func(cmd *cobra.Command, args []string) error {
		script, err := RunProfile()
		if err != nil {
			return err
		}
		fmt.Println(script)
		return nil
	},
}

func RunProfile() (string, error) {
	store, err := govm.ReadStore()
	if err != nil {
		return "", err
	}
	use := store.Use
	if store.Use == "" {
		return "", errorx.Warn("no using version")
	}
	version, e := store.Versions[use]
	if !e {
		return "", errorx.Errorf("using version %s not exist", use)
	}
	tmpl := `export GOROOT="%s"
export PATH=$PATH:“$GOROOT/bin”`
	return fmt.Sprintf(tmpl, filepath.Join(version.Path, "go")), nil
}
