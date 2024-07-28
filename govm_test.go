package govm

import (
	"fmt"
	"testing"
)

func TestRemoteVersionAscend(t *testing.T) {
	versions, err := GetRemoteVersions(true)
	if err != nil {
		t.Error(err)
		return
	}
	for i, version := range versions {
		fmt.Println(i, version)
	}
}

func TestRemoteVersionDescend(t *testing.T) {
	versions, err := GetRemoteVersions(false)
	if err != nil {
		t.Error(err)
		return
	}
	for i, version := range versions {
		fmt.Println(i, version)
	}
}
