package govm

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os/exec"
	"regexp"
	"runtime"
	"slices"
	"strings"
	"time"
)

type GoVersion struct {
	Version string    `json:"version"`
	Stable  bool      `json:"stable"`
	Files   []Version `json:"files"`
}

type Version struct {
	Filename string `json:"filename"`
	Os       string `json:"os"`
	Arch     string `json:"arch"`
	Version  string `json:"version"`
	Sha256   string `json:"sha256"`
	Size     int64  `json:"size"`
	Kind     string `json:"kind"`
}

// GetRemoteVersion returns all available go versions from versionURL without git.
func GetRemoteVersion(ascend bool) ([]GoVersion, error) {
	versionURL, err := GetVersionURL()
	if err != nil {
		return nil, err
	}
	httpClient, err := GetHttpClient()
	httpClient.Timeout = time.Second * 10
	if err != nil {
		return nil, err
	}
	response, err := httpClient.Get(fmt.Sprintf("%s?mode=json&include=all", versionURL))
	if err != nil {
		return nil, err
	}
	var versions []GoVersion
	err = json.NewDecoder(response.Body).Decode(&versions)
	if err != nil {
		return nil, err
	}
	if ascend {
		slices.Reverse(versions)
	}
	return versions, nil
}

func ChooseDownloadURL(version string) (string, string, error) {
	source, err := GetDownloadURL()
	if err != nil {
		return "", "", err
	}
	// the os and arch is same as current tool
	os := runtime.GOOS
	arch := runtime.GOARCH
	var ext string
	if os == "windows" {
		ext = "zip"
	} else {
		ext = "tar.gz"
	}
	filename := fmt.Sprintf("%s.%s-%s.%s", version, os, arch, ext)
	dlurl, err := url.JoinPath(source, filename)
	if err != nil {
		return "", "", err
	}
	return dlurl, filename, err
}

var Buffer = make([]byte, 4096)

const (
	_GoGithubURL = "https://github.com/golang/go"
)

// match tags from ls-remote output
var matchVersion = regexp.MustCompile(`refs/tags/go.+`)

// GetGitRemoteVersions return all available go versions from remote git repository.
// This way need to install git in local.
func GetGitRemoteVersions(ascend bool) ([]string, error) {
	// sort order by version
	sortByVersion := "--sort=-version:refname"
	if ascend {
		sortByVersion = "--sort=version:refname"
	}

	// git ls--remote --tags --sort=version:refname _GoGithubURL
	lsCmd := exec.Command("git", "ls-remote", "--tags", sortByVersion, _GoGithubURL)
	output, err := lsCmd.Output()
	if err != nil {
		execErr := new(exec.ExitError)
		if errors.As(err, &execErr) {
			return nil, errors.New(bytes2string(execErr.Stderr))
		}
		return nil, err
	}

	// find all matched tags by regex
	outputStr := bytes2string(output)
	matches := matchVersion.FindAllString(outputStr, -1)

	// trim prefix
	var list []string
	for _, match := range matches {
		list = append(list, strings.TrimPrefix(match, "refs/tags/"))
	}
	return list, nil
}
