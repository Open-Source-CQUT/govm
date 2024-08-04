package govm

import (
	"fmt"
	"github.com/Open-Source-CQUT/govm/pkg/errorx"
	"github.com/pelletier/go-toml/v2"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	configDir  = ".govm"
	configFile = "config.toml"
)

type Config struct {
	// where to query Go versions, default https://go.dev/dl/?mode=json&include=all
	ListAPI string `toml:"listapi" mapstructure:"listapi" comment:"where to query Go versions, default https://go.dev/dl/?mode=json&include=all"`
	// download URL for Go release archive, default https://dl.google.com/go/
	Mirror string `toml:"mirror" mapstructure:"mirror" comment:"where to download Go release archive, default https://dl.google.com/go/"`
	// proxy for HTTP requests, default use system proxy
	Proxy string `toml:"proxy" mapstructure:"proxy" comment:"http proxy, default use system proxy"`
	// where to install Go, windows: $HOME/AppData/Local/govm/root/ other: /user/local/govm
	Install string `toml:"install" mapstructure:"install" comment:"where to install Go, for windows: $HOME/AppData/Local/govm/, for other: /usr/local/govm"`

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

// ReadConfig read config from config file.
func ReadConfig() (*Config, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return nil, err
	}
	// located at $HOME/.govm
	cfgFile, err := OpenFile(filepath.Join(configDir, configFile), os.O_CREATE|os.O_RDWR, 0644)
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
	configDir, err := GetConfigDir()
	if err != nil {
		return err
	}
	// located at $HOME/.govm
	cfgFile, err := OpenFile(filepath.Join(configDir, configFile), os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer cfgFile.Close()
	return toml.NewEncoder(cfgFile).Encode(cfg)
}

// DefaultConfig returns default config, only for config showing.
func DefaultConfig() (*Config, error) {
	config, err := ReadConfig()
	if err != nil {
		return nil, err
	}
	if config.ListAPI == "" {
		config.ListAPI = _GoDLVersionURL
	}
	if config.Mirror == "" {
		config.Mirror = _GoogleMirror
	}
	if config.Proxy == "" {
		config.Proxy = "(system proxy)"
	}
	if config.Install == "" {
		config.Install, _ = DefaultInstallation()
	}
	return config, nil
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
	_DefaultLinuxInstallation = "/usr/local/govm"
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
	return DefaultInstallation()
}

func DefaultInstallation() (string, error) {
	if runtime.GOOS == "windows" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(homeDir, "/AppData/Local/govm"), nil
	}
	return _DefaultLinuxInstallation, nil
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

const _RootDir = "root"

func GetRootSymLink() (string, error) {
	cfg, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(cfg, _RootDir), nil
}

const _StoreDir = "store"

func GetStoreDir() (string, error) {
	installation, err := GetInstallation()
	if err != nil {
		return "", err
	}
	return filepath.Join(installation, _StoreDir), nil
}

const (
	_DefaultCache = "cache"
)

func GetCacheDir() (string, error) {
	installation, err := GetInstallation()
	if err != nil {
		return "", err
	}
	return filepath.Join(installation, _DefaultCache), err
}

// Pair key=value pair format.
type Pair struct {
	Key, Value string
}

func ParsePair(a string) (Pair, error) {
	parts := strings.SplitN(a, "=", 2)
	if len(parts) != 2 {
		return Pair{}, errorx.Errorf("invalid key=value pair: %s", a)
	}
	return Pair{Key: parts[0], Value: parts[1]}, nil
}

func ParsePairList(args []string) ([]Pair, error) {
	pairs := make([]Pair, 0, len(args))
	for _, arg := range args {
		pair, err := ParsePair(arg)
		if err != nil {
			return nil, err
		}
		pairs = append(pairs, pair)
	}
	return pairs, nil
}
