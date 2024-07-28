package main

import (
	"fmt"
	"github.com/Open-Source-CQUT/govm"
	"github.com/spf13/cobra"
	"regexp"
	"strings"
)

var (
	// number of lines to display, list all if n=-1
	lines int
	// sort by version in ascending order
	ascend bool
	// only show beta versions
	showBeta bool
	// only show rc versions
	showRc bool
	// only show release versions
	showRelease bool
	// show count of matching versions
	showCount bool
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
		result, err := RunSearch(pattern, lines, ascend, showRelease, showBeta, showRc)
		if err != nil {
			return err
		}
		if showCount {
			fmt.Println(len(result))
		} else {
			fmt.Println(strings.Join(result, "\n"))
		}
		return nil
	},
}

func init() {
	searchCmd.Flags().IntVarP(&lines, "lines", "n", defaultLines, "number of lines to display, list all if n=-1")
	searchCmd.Flags().BoolVarP(&showCount, "count", "c", false, "show count of matching versions")
	searchCmd.Flags().BoolVar(&ascend, "ascend", false, "sort by version in ascending order")
	searchCmd.Flags().BoolVar(&showBeta, "beta", false, "only show beta versions")
	searchCmd.Flags().BoolVar(&showRc, "rc", false, "only show rc versions")
	searchCmd.Flags().BoolVar(&showRelease, "release", false, "only show release versions")
}

func RunSearch(pattern string, lines int, ascend, release, beta, rc bool) ([]string, error) {
	// get remote versions
	remoteVersions, err := govm.GetRemoteVersions(ascend)
	if err != nil {
		return nil, err
	}
	return Filter(remoteVersions, pattern, lines, release, beta, rc)
}

func Filter(versions []string, pattern string, lines int, release, beta, rc bool) ([]string, error) {
	// match versions with pattern
	patternMatch := regexp.MustCompile(fmt.Sprintf(`(%s)`, pattern))
	var matchedVersions []string
	if len(pattern) > 0 {
		for _, version := range versions {
			if len(matchedVersions) >= lines {
				break
			}
			if patternMatch.MatchString(version) {
				matchedVersions = append(matchedVersions, version)
			}
		}
	} else {
		matchedVersions = versions
	}

	// filter by beta or rc
	var result []string
	if release || beta || rc {
		for _, version := range matchedVersions {
			if (beta && strings.Contains(version, "beta")) ||
				(rc && strings.Contains(version, "rc")) ||
				(release && !strings.Contains(version, "beta") && !strings.Contains(version, "rc")) {
				result = append(result, version)
			}
		}
	} else {
		result = matchedVersions
	}

	// truncate to specified lines
	if lines == 0 {
		lines = defaultLines
	} else if lines < 0 {
		lines = len(result)
	} else if lines > len(result) {
		lines = len(result)
	}
	return result[:lines], nil
}
