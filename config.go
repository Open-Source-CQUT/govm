package govm

import (
	"fmt"
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
	// where to query Go versions, default https://go.dev/dl/?mode=json&include=all
	ListAPI string `toml:"listapi,commented"`
	// download URL for Go release archive, default https://dl.google.com/go/
	Mirror string `toml:"mirror,commented"`
	// proxy for HTTP requests, default use system proxy
	Proxy string `toml:"proxy,commented"`
	// where to cache downloaded archives, default $HOME/.govm/cache/
	Cache string `toml:"cache,commented"`
	// where to install Go, windows: $HOME/AppData/Local/govm-store/root/
	Install string `toml:"install,commented"`

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
	_GoDLVersionURL = "https://go.dev/dl/?mode=json&include=all"
)

func GetVersionListAPI() (string, error) {
	envVersionURL, found := os.LookupEnv(_EnvVersionURL)
	if found {
		return envVersionURL, nil
	}
	config, err := ReadConfig()
	if err != nil {
		return "", err
	} else if config.ListAPI != "" {
		return config.ListAPI, nil
	}
	return _GoDLVersionURL, nil
}

const (
	_EnvMirror = "GOVM_MIRROR"
	// eg. https://dl.google.com/go/go1.22.5.linux-amd64.tar.gz
	_GoogleMirror = "https://dl.google.com/go/"
	// eg. https://mirrors.aliyun.com/golang/go1.10.1.linux-amd64.tar.gz
	_AliCloudMirror = "https://mirrors.aliyun.com/golang/"
	// eg. https://mirrors.nju.edu.cn/golang/go1.22.5.windows-amd64.msi.sha256
	_NJUDownloadURL = "https://mirrors.nju.edu.cn/golang/"
)

func GetMirror() (string, error) {
	envDownloadURL, found := os.LookupEnv(_EnvMirror)
	if found {
		return envDownloadURL, nil
	}
	config, err := ReadConfig()
	if err != nil {
		return "", err
	} else if config.Mirror != "" {
		return config.Mirror, err
	}
	return _GoogleMirror, nil
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
		return filepath.Join(homeDir, "/AppData/Local/govm"), nil
	}
	return _DefaultLinuxInstallation, err
}

const _Profile = "profile"

func GetProfile() (string, error) {
	dir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, _Profile), nil
}

func GetProfileContent() (string, error) {
	rootSymLink, err := GetRootSymLink()
	if err != nil {
		return "", err
	}

	tmpl := `export GOROOT="%s"
export PATH=$PATH:$GOROOT/bin`
	return fmt.Sprintf(tmpl, filepath.Join(rootSymLink, "go")), err
}
