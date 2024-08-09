package main

import (
	"errors"
	"fmt"
	"github.com/Open-Source-CQUT/govm"
	"github.com/Open-Source-CQUT/govm/pkg/errorx"
	"github.com/spf13/cobra"
	"runtime"
)

var (
	Version = "untag"

	// do not show any tip, warn, error
	silence bool
)

var rootCmd = &cobra.Command{
	Use:           "govm",
	SilenceUsage:  true,
	SilenceErrors: true,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		govm.Silence = silence
	},
	Short: "govm is a tool to manage local Go versions",
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
	rootCmd.AddCommand(purgeCmd)
	rootCmd.AddCommand(profileCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.PersistentFlags().BoolVarP(&silence, "silence", "s", false, "Do not show any tip, warn, error")
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			govm.ErrPrintln(err)
		}
	}()

	// execute command
	err := rootCmd.Execute()
	if !silence && err != nil {
		var kindError errorx.KindError
		if errors.As(err, &kindError) {
			// redirect to stdout
			if kindError.Kind != errorx.ErrorKind {
				govm.Println(kindError.Err.Error())
				return
			}
		}
		govm.ErrPrintln(err.Error())
	}
}
