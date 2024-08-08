package main

import (
	"errors"
	"fmt"
	"github.com/Open-Source-CQUT/govm/pkg/errorx"
	"github.com/spf13/cobra"
	"os"
	"runtime"
)

var (
	Version string
)

var rootCmd = &cobra.Command{
	Use:           "govm",
	SilenceUsage:  true,
	SilenceErrors: true,
	Short:         "govm is a tool to manage local Go versions",
}

func init() {
	rootCmd.SetVersionTemplate(fmt.Sprintf("govm version {{.Version}} %s", runtime.GOOS+"/"+runtime.GOARCH))
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(uninstallCmd)
	rootCmd.AddCommand(useCmd)
	rootCmd.AddCommand(currentCmd)
	rootCmd.AddCommand(cleanCmd)
	rootCmd.AddCommand(profileCmd)
	rootCmd.AddCommand(configCmd)
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
		}
	}()

	// execute command
	err := rootCmd.Execute()
	if err != nil {
		var kindError errorx.KindError
		if errors.As(err, &kindError) {
			// redirect to stdout
			if kindError.Kind != errorx.ErrorKind {
				_, _ = fmt.Fprintln(os.Stdout, kindError.Err.Error())
				return
			}
		}
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
	}
}
