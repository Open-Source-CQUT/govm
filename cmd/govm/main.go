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
	Version:       Version,
	Short:         "Go Version Manager",
	Long:          "Go Version Manager",
}

func init() {
	rootCmd.SetVersionTemplate(fmt.Sprintf("govm version {{.Version}} %s", runtime.GOOS+"/"+runtime.GOARCH))
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(uninstallCmd)
	rootCmd.AddCommand(useCmd)
	rootCmd.AddCommand(currentCmd)
	rootCmd.AddCommand(cleanCmd)
}

func main() {
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
