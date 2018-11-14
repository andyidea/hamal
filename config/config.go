package config

import (
	"encoding/json"
	"fmt"
	"github.com/koding/multiconfig"
)

type Config struct {
	ApiKey    string  `json:"api_key"`
	ApiSecret string  `json:"api_secret"`
	Interval  float64 `json:"interval"`
	Amount    float64 `json:"amount"`
}

var Instance Config

func readConfigFile(obj interface{}, filePath string) error {
	m := multiconfig.New()
	m.Loader = &multiconfig.JSONLoader{Path: filePath}
	err := m.Load(obj)
	if err != nil {
		return err
	}

	return m.Validate(obj)
}

func init() {
	var configFilePath = "config.json"
	if err := readConfigFile(&Instance, configFilePath); err != nil {
		panic(err)
	}

	ind, err := json.MarshalIndent(Instance, "", "	")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Load Config from %s :\n", configFilePath)
	fmt.Println(string(ind))
	fmt.Println("-----------------------------")
}
