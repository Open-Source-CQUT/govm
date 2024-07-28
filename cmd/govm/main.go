package main

import "github.com/spf13/cobra"

var (
	Version string
)

var rootCmd = &cobra.Command{
	Use:     "govm",
	Version: Version,
	Short:   "Go Version Manager",
	Long:    "Go Version Manager",
}

func init() {
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(searchCmd)
}

func main() {
	rootCmd.Execute()
}
