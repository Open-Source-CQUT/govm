package govm

import (
	"encoding/json"
	"fmt"
	"net/url"
	"runtime"
	"slices"
)

type GoVersion struct {
	Version string    `json:"version"`
	Stable  bool      `json:"stable"`
	Files   []Version `json:"files"`
}

type Version struct {
	Filename string `toml:"filename" json:"filename"`
	Os       string `toml:"os" json:"os"`
	Arch     string `toml:"arch" json:"arch"`
	Version  string `toml:"version" json:"version"`
	Sha256   string `toml:"sha256" json:"sha256"`
	Size     uint64 `toml:"size" json:"size"`
	Kind     string `toml:"kind" json:"kind"`
	Path     string `toml:"path"`
	Using    bool
}

// GetRemoteVersion returns all available go versions from versionURL without git.
// If unstable is false, it only returns stable versions, otherwise it returns all versions that includes unstable versions.
func GetRemoteVersion(ascend, unstable bool) ([]Version, error) {
	versionURL, err := GetVersionListAPI()
	if err != nil {
		return nil, err
	}
	httpClient, err := GetHttpClient()
	if err != nil {
		return nil, err
	}
	// get all versions from versionURL
	response, err := httpClient.Get(fmt.Sprintf("%s?mode=json&include=all", versionURL))
	if err != nil {
		return nil, err
	}
	var versions []GoVersion
	err = json.NewDecoder(response.Body).Decode(&versions)
	if err != nil {
		return nil, err
	}

	// filter by current os and arch
	var filterVersions []Version
	for _, goversion := range versions {
		if !unstable && !goversion.Stable {
			continue
		}
		for _, version := range goversion.Files {
			if version.Kind == "archive" &&
				version.Os == runtime.GOOS &&
				version.Arch == runtime.GOARCH {
				filterVersions = append(filterVersions, version)
			}
		}
	}

	// sort by version
	slices.SortFunc(filterVersions, func(v1, v2 Version) int {
		if ascend {
			return CompareVersion(v1.Version, v2.Version)
		}
		return -CompareVersion(v1.Version, v2.Version)
	})

	return filterVersions, nil
}

func ChooseDownloadURL(version string) (string, string, error) {
	source, err := GetMirror()
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

// GetLocalVersions returns the local versions from store.
func GetLocalVersions(ascend bool) ([]Version, error) {
	storeData, err := ReadStore()
	if err != nil {
		return nil, err
	}
	var localList []Version
	for _, v := range storeData.Root {
		if storeData.Use == v.Version {
			v.Using = true
		}
		localList = append(localList, v)
	}
	slices.SortFunc(localList, func(v1, v2 Version) int {
		if ascend {
			return CompareVersion(v1.Version, v2.Version)
		}
		return -CompareVersion(v1.Version, v2.Version)
	})
	return localList, nil
}
