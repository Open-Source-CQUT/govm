package govm

import (
	"github.com/pelletier/go-toml/v2"
	"os"
	"path/filepath"
)

const (
	storeFile = "store.toml"
)

type Store struct {
	Use  string             `toml:"use"`
	Root map[string]Version `toml:"root"`
}

func ReadStore() (*Store, error) {
	store, err := GetInstallation()
	if err != nil {
		return nil, err
	}
	storeFile, err := OpenFile(filepath.Join(store, storeFile), os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	defer storeFile.Close()
	var storeData Store
	err = toml.NewDecoder(storeFile).Decode(&storeData)
	if err != nil {
		return nil, err
	}

	// initialize
	if storeData.Root == nil {
		storeData.Root = make(map[string]Version)
	}
	return &storeData, nil
}

func WriteStore(storeData *Store) error {
	store, err := GetInstallation()
	if err != nil {
		return err
	}
	storeFile, err := OpenFile(filepath.Join(store, storeFile), os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer storeFile.Close()
	return toml.NewEncoder(storeFile).Encode(storeData)
}
