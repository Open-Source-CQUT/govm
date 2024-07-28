package govm

import (
	"github.com/pelletier/go-toml/v2"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
)

const (
	configDir  = ".govm"
	configFile = "config.toml"
)

type Config struct {
	Source       string `toml:"source"`
	Cache        string `toml:"cache"`
	Proxy        string `toml:"proxy"`
	Installation string `toml:"installation"`

	dir string
}

func GetConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	// config dir
	configDir := filepath.Join(homeDir, configDir)
	err = os.MkdirAll(configDir, 0644)
	if err != nil {
		return "", err
	}
	return configDir, nil
}

func OpenConfigFile() (*os.File, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return nil, err
	}
	// located at $HOME/.govm
	cfgFile, err := OpenFile(filepath.Join(configDir, configFile), os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	return cfgFile, nil
}

// ReadConfig read config from config file.
func ReadConfig() (*Config, error) {
	cfgFile, err := OpenConfigFile()
	if err != nil {
		return nil, err
	}
	defer cfgFile.Close()

	var config Config
	err = toml.NewDecoder(cfgFile).Decode(&config)
	if err != nil {
		return nil, err
	}
	config.dir = filepath.Dir(cfgFile.Name())
	return &config, nil
}

// WriteConfig write config into config file.
func WriteConfig(cfg *Config) error {
	cfgFile, err := OpenConfigFile()
	if err != nil {
		return err
	}
	defer cfgFile.Close()
	return toml.NewEncoder(cfgFile).Encode(cfg)
}

const (
	_EnvSourceKey = "GOVM_SOURCE"
	// eg. https://go.dev/dl/go1.22.5.linux-amd64.msi
	_GoSource = "https://go.dev/dl/"
	// eg. https://dl.google.com/go/go1.22.5.linux-amd64.tar.gz
	_GoogleSource = "https://dl.google.com/go/"
	// eg. https://mirrors.aliyun.com/golang/go1.10.1.linux-amd64.tar.gz
	_AliCloudSource = "https://mirrors.aliyun.com/golang/"
)

func GetSource() (string, error) {
	envProxy, found := os.LookupEnv(_EnvSourceKey)
	if found {
		return envProxy, nil
	}
	config, err := ReadConfig()
	if err != nil {
		return "", err
	} else if config.Source != "" {
		return config.Source, err
	}
	return _GoSource, nil
}

const (
	_EnvCacheKey  = "GOVM_CACHE"
	_DefaultCache = "cache"
)

func GetCacheDir() (string, error) {
	envCache, found := os.LookupEnv(_EnvCacheKey)
	if found {
		return envCache, nil
	}
	config, err := ReadConfig()
	if err != nil {
		return "", err
	}
	if config.Cache != "" {
		return "", err
	}
	return filepath.Join(config.dir, _DefaultCache), err
}

const (
	_EnvProxyKey = "GOVM_PROXY"
)

func GetHttpClient() (*http.Client, error) {
	var (
		err    error
		proxy  string
		config *Config
		client = &http.Client{}
	)
	envProxy, found := os.LookupEnv(_EnvProxyKey)
	if found {
		proxy = envProxy
		goto ret
	}
	config, err = ReadConfig()
	if err != nil {
		return nil, err
	} else if config.Proxy != "" {
		proxy = config.Proxy
	}
	goto ret

ret:
	if proxy == "" {
		client.Transport = &http.Transport{Proxy: http.ProxyFromEnvironment}
		return client, nil
	}
	proxyURL, err := url.Parse(proxy)
	if err != nil {
		return nil, err
	}
	client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyURL)}
	return client, nil
}

const (
	_EnvInstallKey = "GOVM_INSTALL"
	// linux, macos, bsd
	_DefaultLinuxInstallation = "/usr/local/govm-store"
)

func GetInstallation() (string, error) {
	envInstallation, found := os.LookupEnv(_EnvInstallKey)
	if found {
		return envInstallation, nil
	}

	config, err := ReadConfig()
	if err != nil {
		return "", err
	} else if config.Installation != "" {
		return config.Installation, err
	}

	// windows: $HOME/AppData/Local/govm-store
	// linux, macOS, bsd and others: /usr/local/govm-store
	if runtime.GOOS == "windows" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(homeDir, "AppData", "Local", "govm-store"), nil
	}
	return _DefaultLinuxInstallation, err
}

func GetRootStore() (string, error) {
	installation, err := GetInstallation()
	if err != nil {
		return "", err
	}
	return filepath.Join(installation, "root"), nil
}
