package govm

import (
	"github.com/stretchr/testify/assert"
	"slices"
	"testing"
)

func TestGetRemoteVersionAscend(t *testing.T) {
	versions, err := GetRemoteVersion(true, true)
	assert.NoError(t, err)
	sorted := slices.IsSortedFunc(versions, func(v1, v2 Version) int {
		return CompareVersion(v1.Version, v2.Version)
	})
	assert.Truef(t, sorted, "Version list is not sorted in ascending order")
}

func TestGetRemoteVersionDescend(t *testing.T) {
	versions, err := GetRemoteVersion(false, true)
	assert.NoError(t, err)
	sorted := slices.IsSortedFunc(versions, func(v1, v2 Version) int {
		return -CompareVersion(v1.Version, v2.Version)
	})
	assert.Truef(t, sorted, "Version list is not sorted in descend order")
}
