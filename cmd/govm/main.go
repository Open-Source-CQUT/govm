package main

import "github.com/spf13/cobra"

var (
	Version string
)

var rootCmd = &cobra.Command{
	Use:          "govm",
	SilenceUsage: true,
	Version:      Version,
	Short:        "Go Version Manager",
	Long:         "Go Version Manager",
}

func init() {
	rootCmd.SetVersionTemplate("{{.Version}}")
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(listCmd)
}

func main() {
	rootCmd.Execute()
}
