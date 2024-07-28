package govm

import (
	"os/exec"
	"regexp"
	"strings"
)

const (
	// eg. https://go.dev/dl/go1.22.5.linux-amd64.msi
	_GoPrefix = "https://go.dev/dl/"
	// eg. https://dl.google.com/go/go1.22.5.linux-amd64.tar.gz
	_GoogleDlPrefix = "https://dl.google.com/go/"
	// eg. https://mirrors.aliyun.com/golang/go1.10.1.linux-amd64.tar.gz
	_AliCloudDlPrefix = "https://mirrors.aliyun.com/golang/"

	_GoGithubURL = "https://github.com/golang/go"
)

// match tags from ls-remote output
var matchVersion = regexp.MustCompile(`refs/tags/go.+`)

// GetRemoteVersions return all available go versions from remote repository.
func GetRemoteVersions(ascend bool) ([]string, error) {
	// sort order by version
	sortByVersion := "--sort=-version:refname"
	if ascend {
		sortByVersion = "--sort=version:refname"
	}

	lsCmd := exec.Command("git", "ls-remote", "--tags", sortByVersion, _GoGithubURL)
	output, err := lsCmd.Output()
	if err != nil {
		return nil, err
	}
	outputStr := bytes2string(output)
	matches := matchVersion.FindAllString(outputStr, -1)

	// trim prefix
	var list []string
	for _, match := range matches {
		list = append(list, strings.TrimPrefix(match, "refs/tags/"))
	}
	return list, nil
}
