package config

import (
	"fmt"
	"github.com/DeAccountSystems/das-lib/common"
	"github.com/fsnotify/fsnotify"
	"github.com/scorpiotzh/mylog"
	"github.com/scorpiotzh/toolib"
)

var (
	Cfg CfgServer
	log = mylog.NewLogger("main", mylog.LevelDebug)
)

func InitCfg(configFilePath string) error {
	if configFilePath == "" {
		configFilePath = "./conf/config.yaml"
	}
	log.Info("config file path：", configFilePath)
	if err := toolib.UnmarshalYamlFile(configFilePath, &Cfg); err != nil {
		return fmt.Errorf("UnmarshalYamlFile err:%s", err.Error())
	}
	log.Info("config file info：", toolib.JsonString(Cfg))
	return nil
}

func AddCfgFileWatcher(configFilePath string) (*fsnotify.Watcher, error) {
	if configFilePath == "" {
		configFilePath = "./conf/config.yaml"
	}
	return toolib.AddFileWatcher(configFilePath, func() {
		log.Info("config file path：", configFilePath)
		if err := toolib.UnmarshalYamlFile(configFilePath, &Cfg); err != nil {
			log.Error("UnmarshalYamlFile err:", err.Error())
		}
		log.Info("config file info：", toolib.JsonString(Cfg))
	})
}

type CfgServer struct {
	Server struct {
		Net            common.DasNetType `json:"net" yaml:"net"`
		HttpServerAddr string            `json:"http_server_addr" yaml:"http_server_addr"`
	} `json:"server" yaml:"server"`
	Chain struct {
		CkbUrl             string `json:"ckb_url" yaml:"ckb_url"`
		IndexUrl           string `json:"index_url" yaml:"index_url"`
		CurrentBlockNumber uint64 `json:"current_block_number" yaml:"current_block_number"`
		ConfirmNum         uint64 `json:"confirm_num" yaml:"confirm_num"`
		ConcurrencyNum     uint64 `json:"concurrency_num" yaml:"concurrency_num"`
	} `json:"chain" yaml:"chain"`
	DB struct {
		Mysql DbMysql `json:"mysql" yaml:"mysql"`
	} `json:"db" yaml:"db"`
	Cache struct {
		Redis struct {
			Addr     string `json:"addr" yaml:"addr"`
			Password string `json:"password" yaml:"password"`
			DbNum    int    `json:"db_num" yaml:"db_num"`
		} `json:"redis" yaml:"redis"`
	} `json:"cache" yaml:"cache"`
	DasLib struct {
		THQCodeHash         string                            `json:"thq_code_hash" yaml:"thq_code_hash"`
		DasContractArgs     string                            `json:"das_contract_args" yaml:"das_contract_args"`
		DasContractCodeHash string                            `json:"das_contract_code_hash" yaml:"das_contract_code_hash"`
		MapDasContract      map[common.DasContractName]string `json:"map_das_contract" yaml:"map_das_contract"`
	} `json:"das_lib" yaml:"das_lib"`
}

type DbMysql struct {
	Addr        string `json:"addr" yaml:"addr"`
	User        string `json:"user" yaml:"user"`
	Password    string `json:"password" yaml:"password"`
	DbName      string `json:"db_name" yaml:"db_name"`
	MaxOpenConn int    `json:"max_open_conn" yaml:"max_open_conn"`
	MaxIdleConn int    `json:"max_idle_conn" yaml:"max_idle_conn"`
}
