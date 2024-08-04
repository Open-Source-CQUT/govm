package main

import (
	"fmt"
	"github.com/Open-Source-CQUT/govm"
	"github.com/spf13/cobra"
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
		localList, err := RunList(pattern, lines, ascend)
		if err != nil {
			return err
		}
		if showCount {
			fmt.Println(len(localList))
		} else {
			for _, version := range localList {
				if version.Using {
					fmt.Printf("%s (*)\n", version.Version)
				} else {
					fmt.Printf("%s\n", version.Version)
				}
			}
		}
		return nil
	},
}

func init() {
	listCmd.Flags().IntVarP(&lines, "lines", "n", defaultLines, "number of lines to display, list all if n=-1")
	listCmd.Flags().BoolVarP(&showCount, "count", "c", false, "show count of matching versions")
	listCmd.Flags().BoolVar(&ascend, "ascend", false, "sort by version in ascending order")
}

func RunList(pattern string, lines int, ascend bool) ([]govm.Version, error) {
	localList, err := govm.GetLocalVersions(ascend)
	if err != nil {
		return nil, err
	}
	filterList, err := Filter(localList, pattern, lines)
	if err != nil {
		return nil, err
	}
	return filterList, nil
}
