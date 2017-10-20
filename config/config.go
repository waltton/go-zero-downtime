package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Config struct {
	Server *Server `json:"server"`
}

type Server struct {
	Host               string `json:"host"`
	Port               string `json:"port"`
	IntervalToShutdown int    `json:"intervalToShutdown"` // in seconds
}

func LoadFromJSONFile(filename string) (cfg *Config, err error) {
	var data []byte

	if data, err = ioutil.ReadFile(filename); err != nil {
		return nil, fmt.Errorf("fail to read the file config file, error: %v", err)
	}

	cfg = new(Config)

	if err = json.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("fail to parse unmarshal data as json, error: %v", err)
	}

	return
}
