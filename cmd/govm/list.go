package main

import (
	"github.com/Open-Source-CQUT/govm"
	"github.com/spf13/cobra"
	"strings"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "list local Go versions",
	Aliases: []string{"ls"},
	RunE: func(cmd *cobra.Command, args []string) error {
		var pattern string
		if len(args) > 0 {
			pattern = args[0]
		}
		localList, err := RunList(pattern, lines, ascend, showRelease, showBeta, showRc)
		if err != nil {
			return err
		}
		if showCount {
			govm.Println(len(localList))
		} else {
			govm.Println(strings.Join(localList, "\n"))
		}
		return nil
	},
}

func init() {
	listCmd.Flags().IntVarP(&lines, "lines", "n", defaultLines, "number of lines to display, list all if n=-1")
	listCmd.Flags().BoolVarP(&showCount, "count", "c", false, "show count of matching versions")
	listCmd.Flags().BoolVar(&ascend, "ascend", false, "sort by version in ascending order")
	listCmd.Flags().BoolVar(&showBeta, "beta", false, "only show beta versions")
	listCmd.Flags().BoolVar(&showRc, "rc", false, "only show rc versions")
	listCmd.Flags().BoolVar(&showRelease, "release", false, "only show release versions")
}

func RunList(pattern string, lines int, ascend, release, beta, rc bool) ([]string, error) {
	localList, err := govm.LocalList(ascend)
	if err != nil {
		return nil, err
	}
	return Filter(localList, pattern, lines, release, beta, rc)
}
