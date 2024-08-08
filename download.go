package govm

import (
	"context"
	"fmt"
	"github.com/Open-Source-CQUT/govm/pkg/errorx"
	"github.com/mholt/archiver/v4"
	"github.com/schollz/progressbar/v3"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

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

var _Buffer = make([]byte, 4096)

// Extract extracts archive file to specified target path, only support .zip and .tar.gz
func Extract(archive *os.File, target string) error {
	name := archive.Name()
	ctx := context.Background()
	if strings.HasSuffix(name, "tar.gz") {
		return ExtractTarGzip(ctx, archive, target)
	} else if strings.HasSuffix(name, "zip") {
		return ExtractZip(ctx, archive, target)
	}
	return errorx.Errorf("unsupported archive format: %s", filepath.Ext(name))
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
			return os.MkdirAll(targetPath, 0755)
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
		_, err = io.CopyBuffer(targetFile, archvieReader, _Buffer)
		if err != nil {
			return err
		}
		return nil
	}
}
