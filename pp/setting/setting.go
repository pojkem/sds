package setting

import (
	ed25519crypto "crypto/ed25519"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/stratosnet/sds/framework/client/cf"
	"github.com/stratosnet/sds/relay/stratoschain"
	"github.com/stratosnet/sds/utils"
)

// REPROTDHTIME 1 hour
const REPROTDHTIME = 60 * 60

// MAXDATA max slice size
const MAXDATA = 1024 * 1024 * 3

// HTTPTIMEOUT  HTTPTIMEOUT second
const HTTPTIMEOUT = 20

// FILEHASHLEN
const FILEHASHLEN = 64

// IMAGEPATH
var IMAGEPATH = "./images/"

// ImageMap download image hash map
var ImageMap = &sync.Map{}

var VIDEOPATH = "./videos"

// DownProssMap download progress map
var DownProssMap = &sync.Map{}

// Config
var Config *config

// ConfigPath
var ConfigPath string

// IsLoad
var IsLoad bool

// UpLoadTaskIDMap
var UpLoadTaskIDMap = &sync.Map{}

// DownloadTaskIDMap
var DownloadTaskIDMap = &sync.Map{}

// socket map
var (
	UpMap     = make(map[string]interface{}, 0)
	DownMap   = make(map[string]interface{}, 0)
	ResultMap = make(map[string]interface{}, 0)
)

//  http code
var (
	FAILCode       = 500
	SUCCESSCode    = 0
	ShareErrorCode = 1002
	IsWindows      bool
)

const (
	Version = "v0.3.0"
	HD_PATH = "m/44'/606'/0'/0/0"
)

type SPBaseInfo struct {
	P2PAddress     string `yaml:"P2PAddress"`
	P2PPublicKey   string `yaml:"P2PPublicKey"`
	NetworkAddress string `yaml:"NetworkAddress"`
}

type config struct {
	Version                     uint32
	VersionShow                 string
	DownloadPathMinLen          int
	Port                        string       `yaml:"Port"`
	NetworkAddress              string       `yaml:"NetworkAddress"`
	Debug                       bool         `yaml:"Debug"`
	PPListDir                   string       `yaml:"PPListDir"`
	AccountDir                  string       `yaml:"AccountDir"`
	ScryptN                     int          `yaml:"scryptN"`
	ScryptP                     int          `yaml:"scryptP"`
	DefPassword                 string       `yaml:"DefPassword"`
	DefSavePath                 string       `yaml:"DefSavePath"`
	StorehousePath              string       `yaml:"StorehousePath"`
	DownloadPath                string       `yaml:"DownloadPath"`
	P2PAddress                  string       `yaml:"P2PAddress"`
	P2PPassword                 string       `yaml:"P2PPassword"`
	WalletAddress               string       `yaml:"WalletAddress"`
	WalletPassword              string       `yaml:"WalletPassword"`
	AutoRun                     bool         `yaml:"AutoRun"`  // is auto login
	Internal                    bool         `yaml:"Internal"` // is internal net
	IsWallet                    bool         `yaml:"IsWallet"` // is wallet
	BPURL                       string       `yaml:"BPURL"`    // bphttp
	IsCheckDefaultPath          bool         `yaml:"IsCheckDefaultPath"`
	IsLimitDownloadSpeed        bool         `yaml:"IsLimitDownloadSpeed"`
	LimitDownloadSpeed          uint64       `yaml:"LimitDownloadSpeed"`
	IsLimitUploadSpeed          bool         `yaml:"IsLimitUploadSpeed"`
	LimitUploadSpeed            uint64       `yaml:"LimitUploadSpeed"`
	IsCheckFileOperation        bool         `yaml:"IsCheckFileOperation"`
	IsCheckFileTransferFinished bool         `yaml:"IsCheckFileTransferFinished"`
	AddressPrefix               string       `yaml:"AddressPrefix"`
	P2PKeyPrefix                string       `yaml:"P2PKeyPrefix"`
	ChainId                     string       `yaml:"ChainId"`
	Token                       string       `yaml:"Token"`
	StratosChainUrl             string       `yaml:"StratosChainUrl"`
	StreamingCache              bool         `yaml:"StreamingCache"`
	RestPort                    string       `yaml:"RestPort"`
	InternalPort                string       `yaml:"InternalPort"`
	SPList                      []SPBaseInfo `yaml:"SPList"`
}

var ostype = runtime.GOOS

// LoadConfig
func LoadConfig(configPath string) error {
	ConfigPath = configPath
	Config = &config{}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		utils.Log("The config at location", configPath, "does not exist")
		return err
	}
	err := utils.LoadYamlConfig(Config, configPath)
	if err != nil {
		return err
	}

	Config.DownloadPathMinLen = 112

	if ostype == "windows" {
		IsWindows = true
		// imagePath = filepath.FromSlash(imagePath)
	} else {
		IsWindows = false
	}
	cf.SetLimitDownloadSpeed(Config.LimitDownloadSpeed, Config.IsLimitDownloadSpeed)
	cf.SetLimitUploadSpeed(Config.LimitUploadSpeed, Config.IsLimitUploadSpeed)

	// todo: we shouldn't call stratoschain package to setup a global variable
	stratoschain.Url = Config.StratosChainUrl
	// Initialize SPMap
	for _, sp := range Config.SPList {
		key := sp.P2PAddress
		if key == "" {
			key = "unknown"
		}
		SPMap.Store(key, sp)
	}

	return nil
}

// CheckLogin
func CheckLogin() bool {
	if WalletAddress == "" {
		utils.ErrorLog("please login")
		return false
	}
	return true
}

// GetSign
func GetSign(str string) []byte {
	sign := ed25519crypto.Sign(P2PPrivateKey, []byte(str))
	utils.DebugLog("GetSign == ", sign)
	return sign
}

// UpChan
var UpChan = make(chan string, 100)

// SetConfig SetConfig
func SetConfig(key, value string) bool {

	if !utils.CheckStructField(key, Config) {
		utils.Log("configuration not found")
		return false
	}

	f, err := os.Open(ConfigPath)
	defer f.Close()
	if err != nil {
		utils.ErrorLog("failed to change configuration file", err)
		return false
	}

	var contents []byte
	contents, err = ioutil.ReadAll(f)
	if err != nil {
		utils.ErrorLog("failed to read configuration file", err)
		return false
	}

	contentStrs := strings.Split(string(contents), "\n")
	newString := ""
	change := false
	keyStr := key + ":"
	for _, str := range contentStrs {
		ss := strings.Split(str, " ")
		if len(ss) > 0 && ss[0] == keyStr {
			if keyStr == "DownloadPath:" {
				if ostype == "windows" {
					value = value + `\`
				} else {
					value = value + `/`
				}
			}
			ns := key + ": " + value
			newString += ns
			newString += "\n"
			change = true
			continue
		}
		newString += str
		newString += "\n"
	}
	if !change {
		return false
	}

	if err = os.Truncate(ConfigPath, 0); err != nil {
		utils.ErrorLog("failed to change configuration file", err)
		return false
	}

	var configOS *os.File
	configOS, err = os.OpenFile(ConfigPath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0777)
	defer configOS.Close()

	if err != nil {
		utils.ErrorLog("failed to change configuration file", err)
		return false
	}

	_, err = configOS.WriteString(newString)
	if err != nil {
		utils.ErrorLog("failed to change configuration file", err)
		return false
	}

	LoadConfig(ConfigPath)
	if !(strings.Contains(strings.ToLower(key), "password") || strings.Contains(strings.ToLower(key), "pw")) {
		utils.Log("finished changing configuration file ", key+": ", value)
	}
	//prefix.SetConfig(Config.AddressPrefix)

	return true
}

func defaultConfig() *config {
	return &config{
		Version:                     3,
		VersionShow:                 Version,
		DownloadPathMinLen:          0,
		Port:                        ":18081",
		NetworkAddress:              "127.0.0.1",
		Debug:                       false,
		PPListDir:                   "./peers",
		AccountDir:                  "./accounts",
		ScryptN:                     4096,
		ScryptP:                     6,
		DefPassword:                 "",
		DefSavePath:                 "",
		StorehousePath:              "",
		DownloadPath:                "",
		P2PAddress:                  "",
		P2PPassword:                 "",
		WalletAddress:               "",
		WalletPassword:              "",
		AutoRun:                     true,
		Internal:                    false,
		IsWallet:                    true,
		BPURL:                       "",
		IsCheckDefaultPath:          false,
		IsLimitDownloadSpeed:        false,
		LimitDownloadSpeed:          0,
		IsLimitUploadSpeed:          false,
		LimitUploadSpeed:            0,
		IsCheckFileOperation:        false,
		IsCheckFileTransferFinished: false,
		AddressPrefix:               "st",
		P2PKeyPrefix:                "stsdsp2p",
		ChainId:                     "stratos-testnet-3",
		Token:                       "ustos",
		StratosChainUrl:             "http://127.0.0.1:1317",
		StreamingCache:              false,
		RestPort:                    "",
		InternalPort:                "",
		SPList:                      []SPBaseInfo{{NetworkAddress: "127.0.0.1:8888"}},
	}
}

func GenDefaultConfig(filePath string) error {
	cfg := defaultConfig()

	return utils.WriteConfig(cfg, filePath)
}
