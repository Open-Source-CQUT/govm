package main

import (
	"fmt"
	"github.com/Open-Source-CQUT/govm"
	"github.com/Open-Source-CQUT/govm/pkg/errorx"
	"github.com/spf13/cobra"
	"path"
	"path/filepath"
	"runtime"
)

const (
	UnixProfileTemplate = `export GOROOT="%s"
export PATH=$PATH:"%s"`
	PowerShellProfileTemplate = `$env:GOROOT = "%s"
$env:PATH += ";%s"`

	bash       = "bash"
	gitbash    = "gitbash"
	powershell = "powershell"
)

var (
	shell string
)

var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Show profile env",
	RunE: func(cmd *cobra.Command, args []string) error {
		script, err := RunProfile(shell)
		if err != nil {
			return err
		}
		govm.Println(script)
		return nil
	},
}

func init() {
	// set default shell for different platforms
	defaultShell := bash
	if runtime.GOOS == "windows" {
		defaultShell = gitbash
	}
	profileCmd.Flags().StringVar(&shell, "shell", defaultShell, "shell to use for profile")
}

func RunProfile(shell string) (string, error) {
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
	switch shell {
	case gitbash:
		return GitBashProfile(version)
	case bash:
		return BashProfile(version)
	case powershell:
		return PowershellProfile(version)
	default:
		return BashProfile(version)
	}
}

// GitBashProfile return a profile script for git bash, it must convert windows path to unix path.
func GitBashProfile(version govm.Version) (string, error) {
	if runtime.GOOS != "windows" {
		govm.Warnf("git bash for windows only support windows platform")
	}
	GOROOT := path.Join(govm.ToUnixPath(version.Path), "go")
	GOROOTBIN := `$GOROOT/bin`
	return fmt.Sprintf(UnixProfileTemplate, GOROOT, GOROOTBIN), nil
}

// PowershellProfile return a profile script for powershell
func PowershellProfile(version govm.Version) (string, error) {
	return fmt.Sprintf(PowerShellProfileTemplate, filepath.Join(version.Path, "go"), filepath.Join(version.Path, "go", "bin")), nil
}

// BashProfile return a profile script for bash, it is compatible with most linux shells.
func BashProfile(version govm.Version) (string, error) {
	if runtime.GOOS == "windows" {
		govm.Warnf("bash style profile not support for windows, you should use --shell=gitbash")
	}
	return fmt.Sprintf(UnixProfileTemplate, filepath.Join(version.Path, "go"), "$GOROOT/bin"), nil
}
