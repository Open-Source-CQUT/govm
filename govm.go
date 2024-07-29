package govm

import (
	"context"
	"fmt"
	"github.com/mholt/archiver/v4"
	"github.com/schollz/progressbar/v3"
	"io"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"slices"
	"strings"
	"time"
)

const (
	_GoGithubURL = "https://github.com/golang/go"
)

var Buffer = make([]byte, 4096)

// match tags from ls-remote output
var matchVersion = regexp.MustCompile(`refs/tags/go.+`)

// GetRemoteVersions return all available go versions from remote repository.
func GetRemoteVersions(ascend bool) ([]string, error) {
	// sort order by version
	sortByVersion := "--sort=-version:refname"
	if ascend {
		sortByVersion = "--sort=version:refname"
	}

	// git ls--remote --tags --sort=version:refname _GoGithubURL
	lsCmd := exec.Command("git", "ls-remote", "--tags", sortByVersion, _GoGithubURL)
	output, err := lsCmd.Output()
	if err != nil {
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

func ChooseDownloadURL(version string) (string, string, error) {
	source, err := GetSource()
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

func DownloadProcessBar(length int64, description string, finishedTip string) *progressbar.ProgressBar {
	return progressbar.NewOptions64(length,
		progressbar.OptionSetWriter(os.Stdout),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(40),
		progressbar.OptionThrottle(65*time.Millisecond),
		progressbar.OptionShowCount(),
		progressbar.OptionShowElapsedTimeOnFinish(),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionSetDescription(description),
		progressbar.OptionOnCompletion(func() {
			fmt.Fprint(os.Stdout, finishedTip)
		}),
	)
}

// ExtractTarGzip extract tar.gz from reader and save to target.
func ExtractTarGzip(ctx context.Context, reader io.Reader, target string) error {
	gzip := archiver.Gz{}
	gzipReader, err := gzip.OpenReader(reader)
	if err != nil {
		return err
	}
	defer gzipReader.Close()

	tar := archiver.Tar{}

	return tar.Extract(ctx, gzipReader, nil, extractHandler(target))
}

// ExtractZip extract zip from reader and save to target.
func ExtractZip(ctx context.Context, reader io.Reader, target string) error {
	zip := archiver.Zip{}
	return zip.Extract(ctx, reader, nil, extractHandler(target))
}

func extractHandler(target string) archiver.FileHandler {
	return func(ctx context.Context, f archiver.File) error {
		targetPath := filepath.Join(target, f.NameInArchive)
		// mkdir if it is dir
		if f.IsDir() {
			return os.MkdirAll(targetPath, f.Mode())
		}
		// copy to target if is a file
		targetFile, err := OpenFile(targetPath, os.O_CREATE|os.O_RDWR, f.Mode())
		if err != nil {
			return err
		}
		defer targetFile.Close()
		archvieReader, err := f.Open()
		if err != nil {
			return err
		}
		defer archvieReader.Close()
		// copy file to target
		_, err = io.CopyBuffer(targetFile, archvieReader, Buffer)
		if err != nil {
			return err
		}
		return nil
	}
}

func LocalList(ascend bool) ([]string, error) {
	storeData, err := ReadStore()
	if err != nil {
		return nil, err
	}
	var localList []string
	for version, _ := range storeData.Root {
		localList = append(localList, version)
	}
	slices.SortFunc(localList, func(v1, v2 string) int {
		return -CompareVersion(v1, v2)
	})
	if ascend {
		slices.Reverse(localList)
	}
	return localList, nil
}
