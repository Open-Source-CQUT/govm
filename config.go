package govm

import (
	"github.com/pelletier/go-toml/v2"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

const (
	configDir  = ".govm"
	configFile = "config.toml"
)

type Config struct {
	VersionURL  string `toml:"versionURL,commented" comment:"where to query Go versions, default https://go.dev/dl/?mode=json&include=all"`
	DownloadURL string `toml:"downloadURL,commented" comment:"download URL for Go release archive, default https://dl.google.com/go/"`
	Proxy       string `toml:"proxy,commented" comment:"proxy for HTTP requests, default use system proxy"`
	Cache       string `toml:"cache,commented" comment:"where to cache downloaded archives, default $HOME/.govm/cache/"`
	Install     string `toml:"install,commented" comment:"where to install Go, windows: $HOME/AppData/Local/govm-store/root/ "`

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
	_EnvVersionURL  = "GOVM_VERSION"
	_GoDLVersionURL = "https://go.dev/dl/"
)

func GetVersionURL() (string, error) {
	envVersionURL, found := os.LookupEnv(_EnvVersionURL)
	if found {
		return envVersionURL, nil
	}
	config, err := ReadConfig()
	if err != nil {
		return "", err
	} else if config.VersionURL != "" {
		return config.VersionURL, nil
	}
	return _GoDLVersionURL, nil
}

const (
	_EnvDownloadURL = "GOVM_DOWNLOAD"
	// eg. https://dl.google.com/go/go1.22.5.linux-amd64.tar.gz
	_GoogleDownloadURL = "https://dl.google.com/go/"
	// eg. https://mirrors.aliyun.com/golang/go1.10.1.linux-amd64.tar.gz
	_AliCloudDL = "https://mirrors.aliyun.com/golang/"
	// eg. https://mirrors.ustc.edu.cn/golang/go1.10.2.freebsd-386.tar.gz.asc
	_USTCDownloadURL = "https://mirrors.ustc.edu.cn/golang/"
	// eg. https://go.dev/dl/go1.22.5.windows-amd64.msi.sha256
	_NJUDownloadURL = "https://mirrors.nju.edu.cn/golang/"
)

func GetDownloadURL() (string, error) {
	envDownloadURL, found := os.LookupEnv(_EnvDownloadURL)
	if found {
		return envDownloadURL, nil
	}
	config, err := ReadConfig()
	if err != nil {
		return "", err
	} else if config.DownloadURL != "" {
		return config.DownloadURL, err
	}
	return _GoogleDownloadURL, nil
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
	} else if config.Cache != "" {
		return config.Cache, nil
	}
	return filepath.Join(config.dir, _DefaultCache), err
}

func GetHttpClient() (*http.Client, error) {
	var (
		err    error
		proxy  string
		config *Config
		client = &http.Client{Timeout: time.Second * 10}
	)
	config, err = ReadConfig()
	if err != nil {
		return nil, err
	} else if config.Proxy != "" {
		proxy = config.Proxy
	}

	// get proxy from env
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
	} else if config.Install != "" {
		return config.Install, nil
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
