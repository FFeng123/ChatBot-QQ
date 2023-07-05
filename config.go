package main

import (
	"encoding/json"
	"os"

	"github.com/apex/log"
)

type ChatAPIArgsType struct {
	API_max_length  int     `json:"max_length "`
	API_top_p       float32 `json:"top_p"`
	API_temperature float32 `json:"temperature"`
}

type ConfigType struct {
	CQHttpAddr string `json:"api_cqhttp"`
	ChatAPI    string `json:"api_chat"`

	StaticMem [][2]string `json:"mem_static"`
	MemLength int         `json:"mem_length"`
	ThisQQ    int64       `json:"this_qq"`

	ChatAPIArgs ChatAPIArgsType `json:"chat_args"`
}

const ConfigFile = "config.json"

var Config ConfigType

func ReadConfig() {
	Config = ConfigType{
		CQHttpAddr: "ws://127.0.0.1:80",
		ChatAPI:    "http://127.0.0.1:80",
		StaticMem:  [][2]string{{"用户喵喵喵", "Bot喵喵喵"}},
		ThisQQ:     123123,
		MemLength:  50,
		ChatAPIArgs: ChatAPIArgsType{
			API_max_length:  2048,
			API_top_p:       0.95,
			API_temperature: 0.8,
		},
	}

	_, err := os.Stat(ConfigFile)
	if os.IsNotExist(err) {
		data, _ := json.Marshal(Config)
		log.Info("已创建默认配置")
		os.WriteFile(ConfigFile, data, 0660)
		os.Exit(0)
		return
	}

	d, err := os.ReadFile(ConfigFile)
	if err == nil {
		err = json.Unmarshal(d, &Config)
		if err == nil {
			return
		}
	}

	log.Error(err.Error())
	os.Exit(0)
}
