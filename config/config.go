package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

type Config struct {
	BootsrapServer []string          `json:"bootsrapServer"`
	Topic          string            `json:"topic"`
	Producer       ProducerConfig    `json:"producer"`
	Consumer       ConsumerConfig    `json:"consumer"`
	Simple         SimpleModeConfig  `json:"simple"`
	AllTime        AllTimeModeConfig `json:"allTime"`
	LogFile        string            `json:"logFile"`
}

type ProducerConfig struct {
	ClientID string `json:"clientId"`
	Acks     string `json:"acks"`
}

type ConsumerConfig struct {
	GroupID         string `json:"groupId"`
	AutoOffsetReset string `json:"autoOffsetReset"`
}

type SimpleModeConfig struct {
	MsgCount       int `json:"msgCount"`
	MsgIntervalSec int `json:"msgIntervalSec"`
}

type AllTimeModeConfig struct {
	MsgIntervalSec int    `json:"msgIntervalSec"`
	RttFile        string `json:"rttFile"`
}

func (config *Config) InitConfig(configFilePath string) (*Config, error) {
	if configFilePath == "" {
		return nil, errors.New("[CONFIG] config file path is empty!")
	}
	configRawData, err := os.Open(configFilePath)
	if err != nil {
		return nil, errors.New("[CONFIG] config file path is wrong!")
	}
	configByteValue, err := ioutil.ReadAll(configRawData)
	if err != nil {
		return nil, errors.New("[CONFIG] config read io error!")
	}

	configData := &Config{}
	if err := json.Unmarshal(configByteValue, configData); err != nil {
		return nil, errors.New("[CONFIG] config json parse(unmarshal) error!")
	}
	return configData, nil
}
