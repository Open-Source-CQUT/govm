package main

import (
	"fmt"
	"github.com/Open-Source-CQUT/govm"
	"github.com/Open-Source-CQUT/govm/pkg/errorx"
	"github.com/spf13/cobra"
	"path/filepath"
	"runtime"
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
	using, err := govm.GetUsingVersion()
	if err != nil {
		return "", err
	}
	version, e := store.Versions[using]
	if !e {
		return "", errorx.Errorf("using version %s not exist", using)
	}
	tmpl := `export GOROOT="%s"
export PATH=$PATH:"$GOROOT/bin"`

	if runtime.GOOS == "windows" {
		tmpl = `export GOROOT="%s"
export PATH=$PATH:"$GOROOT\bin"`
	}
	return fmt.Sprintf(tmpl, filepath.Join(version.Path, "go")), nil
}
