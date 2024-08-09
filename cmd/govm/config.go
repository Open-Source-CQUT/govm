package main

import (
	"fmt"
	"github.com/Open-Source-CQUT/govm"
	"github.com/Open-Source-CQUT/govm/pkg/errorx"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"strings"
)

var (
	// write config
	write bool
	// unset config
	unset bool
)

var configCmd = &cobra.Command{
	Use:     "config",
	Aliases: []string{"cfg"},
	Short:   "Manage govm configs",
	RunE: func(cmd *cobra.Command, args []string) error {
		if write {
			return WriteConfig(args)
		} else if unset {
			return UnsetConfig(args)
		}
		config, err := ShowConfig()
		if err != nil {
			return err
		}
		govm.Println(config)
		return nil
	},
}

func init() {
	configCmd.Flags().BoolVarP(&write, "write", "w", false, "write configs with key=value pairs")
	configCmd.Flags().BoolVarP(&unset, "unset", "u", false, "unset configs with keys")
}

func WriteConfig(candidates []string) error {
	pairs, err := govm.ParsePairList(candidates)
	if err != nil {
		return err
	}

	// load config
	config, err := govm.ReadConfig()
	if err != nil {
		return err
	}

	// convert struct to map
	cfgMap := make(map[string]any)
	err = mapstructure.Decode(config, &cfgMap)
	if err != nil {
		return err
	}

	for _, pair := range pairs {
		if _, e := cfgMap[pair.Key]; !e {
			return errorx.Warnf("key %s not found in config", pair.Key)
		}
		if pair.Value == "" {
			return errorx.Warnf("invalid value for key: %s=%s", pair.Key, pair.Value)
		}
		cfgMap[pair.Key] = pair.Value
	}

	// convert map to struct
	err = mapstructure.Decode(&cfgMap, config)
	if err != nil {
		return err
	}

	// save config
	err = govm.WriteConfig(config)
	if err != nil {
		return err
	}

	return nil
}

func UnsetConfig(keys []string) error {
	// load config
	config, err := govm.ReadConfig()
	if err != nil {
		return err
	}

	// convert struct to map
	cfgMap := make(map[string]any)
	err = mapstructure.Decode(config, &cfgMap)
	if err != nil {
		return err
	}

	for _, key := range keys {
		if _, e := cfgMap[key]; !e {
			return errorx.Warnf("key %s not found in config", key)
		}
		delete(cfgMap, key)
	}

	// convert map to struct
	err = mapstructure.Decode(&cfgMap, config)
	if err != nil {
		return err
	}

	// save config
	err = govm.WriteConfig(config)
	if err != nil {
		return err
	}
	return nil
}

func ShowConfig() (string, error) {
	// load config
	config, err := govm.DefaultConfig()
	if err != nil {
		return "", err
	}
	// convert struct to map
	cfgMap := make(map[string]any)
	err = mapstructure.Decode(config, &cfgMap)
	if err != nil {
		return "", err
	}
	var result strings.Builder
	for k, v := range cfgMap {
		result.WriteString(fmt.Sprintf("%s=%v\n", k, v))
	}
	return result.String(), nil
}
