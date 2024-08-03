package main

import (
	"fmt"
	"github.com/Open-Source-CQUT/govm"
	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"
	"regexp"
)

var (
	// number of lines to display, list all if n=-1
	lines int
	// sort by version in ascending order
	ascend bool
	// show count of matching versions
	showCount bool
	// include unstable versions
	unstable bool
)

const defaultLines = 20

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "search available go versions from remote repository",
	RunE: func(cmd *cobra.Command, args []string) error {
		var pattern string
		if len(args) > 0 {
			pattern = args[0]
		}
		result, err := RunSearch(pattern, lines, ascend, unstable)
		if err != nil {
			return err
		}
		if showCount {
			fmt.Println(len(result))
		} else {
			for _, version := range result {
				fmt.Printf("%6s\t%-10s\t%8s\n", version.Sha256[:6], version.Version, humanize.Bytes(version.Size))
			}
		}
		return nil
	},
}

func init() {
	searchCmd.Flags().IntVarP(&lines, "lines", "n", defaultLines, "number of lines to display, list all if n=-1")
	searchCmd.Flags().BoolVarP(&showCount, "count", "c", false, "show count of matching versions")
	searchCmd.Flags().BoolVar(&ascend, "ascend", false, "sort by version in ascending order")
	searchCmd.Flags().BoolVar(&unstable, "unstable", false, "show unstable versions")
}

func RunSearch(pattern string, lines int, ascend, unstable bool) ([]govm.Version, error) {
	// get remote versions
	remoteVersions, err := govm.GetRemoteVersion(ascend, unstable)
	if err != nil {
		return nil, err
	}
	return Filter(remoteVersions, pattern, lines)
}

func Filter(versions []govm.Version, pattern string, lines int) ([]govm.Version, error) {
	// match versions with pattern
	patternMatch := regexp.MustCompile(fmt.Sprintf(`(%s)`, pattern))
	var matchedVersions []govm.Version
	if len(pattern) > 0 {
		for _, version := range versions {
			if patternMatch.MatchString(version.Version) {
				matchedVersions = append(matchedVersions, version)
			}
		}
	} else {
		matchedVersions = versions
	}

	// truncate to specified lines
	if lines == 0 {
		lines = defaultLines
	} else if lines < 0 || lines > len(matchedVersions) {
		lines = len(matchedVersions)
	}
	return matchedVersions[:lines], nil
}
