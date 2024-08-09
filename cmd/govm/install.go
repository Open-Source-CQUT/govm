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
)

var (
	use bool
)

var installCmd = &cobra.Command{
	Use:     "install",
	Short:   "Install specified Go version",
	Aliases: []string{"i"},
	RunE: func(cmd *cobra.Command, args []string) error {
		var version string
		if len(args) > 0 {
			version = args[0]
		}
		return RunInstall(version, use)
	},
}

func init() {
	installCmd.Flags().BoolVar(&use, "use", false, "install then use")
}

func RunInstall(v string, use bool) error {
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

	// if not specified, use the latest stable version
	if (version == "") && len(remoteVersions) > 0 {
		for _, remoteVersion := range remoteVersions {
			if remoteVersion.Stable {
				downloadVersion = remoteVersion
				break
			}
		}
	} else {
		// find the specified version
		for _, remoteVersion := range remoteVersions {
			if remoteVersion.Version == version {
				downloadVersion = remoteVersion
				break
			}
		}
		if downloadVersion.Filename == "" {
			return errorx.Warnf("no matching version found for %s", version)
		}
	}

	// check if is already installed
	storeData, err := govm.ReadStore()
	if err != nil {
		return err
	}
	foundV, exist := storeData.Versions[downloadVersion.Version]
	if exist {
		_, err := os.Stat(foundV.Path)
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return err
		} else if err == nil {
			return errorx.Warnf("%s already installed", downloadVersion.Version)
		}
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
	_ = os.MkdirAll(storePath, 0755)

	// extract
	govm.Tipf("Extract %s to local store", filepath.Base(archiveFile.Name()))
	if err = govm.Extract(archiveFile, storePath); err != nil {
		return err
	}

	// store meta info
	downloadVersion.Path = storePath
	err = govm.AppendVersion(downloadVersion)
	if err != nil {
		return err
	}
	_ = archiveFile.Close()
	govm.Tipf("Remove archive from cache")
	err = os.RemoveAll(archiveFile.Name())
	if err != nil {
		return err
	}
	govm.Tipf("Version %s installed", downloadVersion.Version)
	// whether to use
	if use {
		return RunUse(downloadVersion.Version)
	}
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
	// find in cache
	_, err = os.Stat(cacheFilename)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	} else if err == nil { // found in cache
		govm.Tipf("Found %s from cache", version.Version)
		return govm.OpenFile(cacheFilename, os.O_CREATE|os.O_RDWR, 0755)
	}

	// find version from remote
	govm.Tipf("Fetch %s from %s", version.Version, downloadURL)

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
	cacheFile, err := govm.OpenFile(cacheFilename, os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		return nil, err
	}

	processBar := io.Discard
	if !silence {
		// create progress bar
		processBar = (io.Writer)(govm.DownloadProcessBar(response.ContentLength,
			fmt.Sprintf("Downloading %s", filename), "\n"))
	}

	// sha256 hash writer
	hash := sha256.New()

	// copy to cache file
	_, err = io.CopyBuffer(io.MultiWriter(hash, processBar, cacheFile), response.Body, make([]byte, 4096))
	if err != nil {
		return cacheFile, err
	}

	if version.Sha256 != "" {
		h64 := make([]byte, 64)
		hex.Encode(h64, hash.Sum(nil))
		if !bytes.Equal(h64, govm.String2bytes(version.Sha256)) {
			return cacheFile, errorx.Error("sha256 hash check failed")
		}
	}

	// seek to start
	_, err = cacheFile.Seek(0, io.SeekStart)
	if err != nil {
		return cacheFile, err
	}
	return cacheFile, nil
}
