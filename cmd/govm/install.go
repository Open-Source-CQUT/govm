package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/Open-Source-CQUT/govm"
	"github.com/spf13/cobra"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

var installCmd = &cobra.Command{
	Use:     "install",
	Short:   "install specified version of go",
	Aliases: []string{"i"},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 { // install latest version
			var latestVersion string
			versions, err := govm.GetRemoteVersions(false)
			if err != nil {
				return err
			}
			if len(versions) > 0 {
				latestVersion = versions[0]
			}
			return RunInstall(latestVersion)
		} else if len(args) == 1 { // specified one version
			return RunInstall(args[0])
		} else { // specified multiple versions
			for _, arg := range args {
				err := RunInstall(arg)
				if err != nil {
					return err
				}
			}
		}
		return nil
	},
}

func RunInstall(version string) error {
	if !strings.HasPrefix(version, "go") {
		version = "go" + version
	}
	if valid := govm.IsValidVersion(version); !valid {
		return fmt.Errorf("invalid version: %s", version)
	}
	// check if is already installed
	localList, err := govm.LocalList(false)
	if err != nil {
		return err
	}
	if slices.Contains(localList, version) {
		govm.Printf("%s already installed", version)
		return nil
	}
	// download the specified version
	archiveFile, err := DownloadVersion(version)
	if err != nil {
		return err
	}
	defer archiveFile.Close()
	store, err := govm.GetRootStore()
	if err != nil {
		return err
	}
	storePath := filepath.Join(store, version)
	if err = Extract(archiveFile, storePath); err != nil {
		return err
	}
	storeData, err := govm.ReadStore()
	if err != nil {
		return err
	}
	storeData.Root[version] = storePath
	err = govm.WriteStore(storeData)
	if err != nil {
		return err
	}
	govm.Printf("%s installed", version)
	return nil
}

// DownloadVersion downloads the specified version of go, and returns the downloaded archive file.
// The archive file be closed by the caller.
func DownloadVersion(version string) (*os.File, error) {

	// check if version already in the cache
	downloadURL, filename, err := govm.ChooseDownloadURL(version)
	if err != nil {
		return nil, err
	}
	cacheDir, err := govm.GetCacheDir()
	if err != nil {
		return nil, err
	}
	cacheFilename := filepath.Join(cacheDir, filename)
	_, err = os.Stat(cacheFilename)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	} else if err == nil { // found in cache
		govm.Printf("Found %s from cache", version)
		return govm.OpenFile(cacheFilename, os.O_CREATE|os.O_RDWR, 0644)
	}

	// find version from remote
	remoteVersions, err := govm.GetRemoteVersions(false)
	if err != nil {
		return nil, err
	}
	if !slices.Contains(remoteVersions, version) {
		return nil, fmt.Errorf("version %s not found", version)
	}
	govm.Printf("Found %s from remote", version)

	// not in cache, download from remote
	client, err := govm.GetHttpClient()
	if err != nil {
		return nil, err
	}
	request, _ := http.NewRequest("GET", downloadURL, nil)
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	cacheFile, err := govm.OpenFile(cacheFilename, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}

	// create progress bar
	processBar := govm.DownloadProcessBar(response.ContentLength,
		fmt.Sprintf("Downloading %s", filename), "\n")

	// copy to cache file
	_, err = io.CopyBuffer(io.MultiWriter(processBar, cacheFile), response.Body, make([]byte, 4096))
	if err != nil {
		return nil, err
	}
	// seek to start
	_, err = cacheFile.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}
	return cacheFile, nil
}

func Extract(archive *os.File, target string) error {
	name := archive.Name()
	ctx := context.Background()
	if strings.HasSuffix(name, "tar.gz") {
		return govm.ExtractTarGzip(ctx, archive, target)
	} else if strings.HasSuffix(name, "zip") {
		return govm.ExtractZip(ctx, archive, target)
	}
	return fmt.Errorf("unsupported archive format: %s", filepath.Ext(name))
}
