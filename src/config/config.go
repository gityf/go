package config

import (
	"io/ioutil"
	"encoding/json"
)

type ServerConfigItem struct {
	// http listem port
	Port          int    `json:"port"`
	// http pref listen port
	PprofPort     int    `json:"profport"`
	// limit qps count for system self-protect
	Limitqps      int    `json:"limitqps"`
	// dispatch Concurrency count
	EventConcurrency       int `json:"EventConcurrency"`
	// dispatch message length
	EventMessageQueueLen   int `json:"EventMessageQueueLen"`
}

type ServerConfigAll struct {
	Profile        string `json:"profile"`
	Dev  ServerConfigItem `json:"dev"`
	Test ServerConfigItem `json:"test"`
	Sbox ServerConfigItem `json:"sbox"`
	Prod ServerConfigItem `json:"prod"`
}

var CfgAll ServerConfigAll
var Cfg ServerConfigItem

func ParseConf(file string) (err error) {
	cnt, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}

	err = json.Unmarshal(cnt, &CfgAll);
	switch CfgAll.Profile {
	case "dev":
		Cfg = CfgAll.Dev
	case "test":
		Cfg = CfgAll.Test
	case "sbox":
		Cfg = CfgAll.Sbox
	case "prod":
		Cfg = CfgAll.Prod
	default:
		Cfg = CfgAll.Dev
	}
	return
}