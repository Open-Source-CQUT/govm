package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/Open-Source-CQUT/govm"
	"github.com/Open-Source-CQUT/govm/pkg/errorx"
	"github.com/spf13/cobra"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"slices"
)

var installCmd = &cobra.Command{
	Use:     "install",
	Short:   "install specified version of go",
	Aliases: []string{"i"},
	RunE: func(cmd *cobra.Command, args []string) error {
		var version string
		if len(args) > 0 {
			version = args[0]
		}
		return RunInstall(version)
	},
}

func RunInstall(v string) error {
	var version string
	if v != "" {
		cv, ok := govm.CheckVersion(v)
		if !ok {
			return errorx.Warnf("invalid version: %s", v)
		}
		version = cv
	}

	// get available versions
	remoteVersions, err := govm.GetRemoteVersion(false, true)
	if err != nil {
		return err
	}

	var downloadVersion govm.Version

	// if not specified, use the latest available version
	if (version == "") && len(remoteVersions) > 0 {
		downloadVersion = remoteVersions[0]
	} else {
		filterVersions, err := Filter(remoteVersions, version, -1)
		if err != nil {
			return err
		}
		if len(filterVersions) == 0 {
			return errorx.Warnf("no matching version found for %s", version)
		}
		// find the latest from the matched versions
		// such as go1.22 -> go1.22.latest
		downloadVersion = slices.MaxFunc(filterVersions, func(v1, v2 govm.Version) int {
			return -govm.CompareVersion(v1.Version, v2.Version)
		})
	}

	// check if is already installed
	locals, err := govm.GetLocalVersions(false)
	if err != nil {
		return err
	}
	if slices.ContainsFunc(locals, func(v govm.Version) bool {
		return v.Version == downloadVersion.Version
	}) {
		return errorx.Warnf("%s already installed", version)
	}

	// download the specified version
	archiveFile, err := DownloadVersion(downloadVersion)
	if err != nil {
		if archiveFile != nil {
			archiveFile.Close()
			os.Remove(archiveFile.Name())
		}
		return err
	}
	defer archiveFile.Close()

	store, err := govm.GetStoreDir()
	if err != nil {
		return err
	}
	storePath := filepath.Join(store, downloadVersion.Version)

	// extract
	govm.Tipf("Extract %s to local store", filepath.Base(archiveFile.Name()))
	if err = govm.Extract(archiveFile, storePath); err != nil {
		return err
	}

	// store meta info
	storeData, err := govm.ReadStore()
	if err != nil {
		return err
	}
	downloadVersion.Path = storePath
	storeData.Root[downloadVersion.Version] = downloadVersion
	err = govm.WriteStore(storeData)
	if err != nil {
		return err
	}
	govm.Tipf("%s installed", downloadVersion.Version)
	return nil
}

// DownloadVersion downloads the specified version of go, and returns the downloaded archive file.
// The archive file be closed by the caller.
func DownloadVersion(version govm.Version) (*os.File, error) {

	// check if version already in the cache
	downloadURL, filename, err := govm.ChooseDownloadURL(version.Version)
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
		govm.Tipf("Found %s from cache", version.Version)
		return govm.OpenFile(cacheFilename, os.O_CREATE|os.O_RDWR, 0644)
	}

	// find version from remote
	govm.Tipf("Found %s from %s", version.Version, downloadURL)

	// not in cache, download from remote
	client, err := govm.GetHttpClient()
	if err != nil {
		return nil, err
	}
	// unset default timeout
	client.Timeout = 0
	request, _ := http.NewRequest("GET", downloadURL, nil)
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, errorx.Errorf("download failed: %s", response.Status)
	}
	cacheFile, err := govm.OpenFile(cacheFilename, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}

	// create progress bar
	processBar := govm.DownloadProcessBar(response.ContentLength,
		fmt.Sprintf("Downloading %s", filename), "\n")

	// sha256 hash writer
	hash := sha256.New()

	// copy to cache file
	_, err = io.CopyBuffer(io.MultiWriter(hash, processBar, cacheFile), response.Body, make([]byte, 4096))
	if err != nil {
		return cacheFile, err
	}

	h64 := make([]byte, 64)
	hex.Encode(h64, hash.Sum(nil))
	if !bytes.Equal(h64, govm.String2bytes(version.Sha256)) {
		return cacheFile, errorx.Error("sha256 hash check failed")
	}

	// seek to start
	_, err = cacheFile.Seek(0, io.SeekStart)
	if err != nil {
		return cacheFile, err
	}
	return cacheFile, nil
}
