package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

type Config struct {
	BootstrapServer []string          `json:"bootstrapServer"`
	Topic           string            `json:"topic"`
	Producer        ProducerConfig    `json:"producer"`
	Consumer        ConsumerConfig    `json:"consumer"`
	Simple          SimpleModeConfig  `json:"simple"`
	AllTime         AllTimeModeConfig `json:"allTime"`
}

type ProducerConfig struct {
	ClientID      string `json:"clientId"`
	Acks          string `json:"acks"`
	ValueByteSize int    `json:"valueByteSize"`
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
	LogFile        string `json:"logFile"`
}

func InitConfig(configFilePath string) (*Config, error) {
	if configFilePath == "" {
		return nil, errors.New("config file path is empty!")
	}
	configRawData, err := os.Open(configFilePath)
	if err != nil {
		return nil, errors.New("config file path is wrong!")
	}
	configByteValue, err := ioutil.ReadAll(configRawData)
	if err != nil {
		return nil, errors.New("config read io error!")
	}

	configData := &Config{}
	if err := json.Unmarshal(configByteValue, configData); err != nil {
		return nil, errors.New("config json parse(unmarshal) error!")
	}

	if err := configData.isValidBootstrapServer(); err != nil {
		return nil, errors.New("bootstrap server is not valid")
	}
	if err := configData.isValidProducerAcks(); err != nil {
		return nil, err
	}

	if err := isValidString(configData.Topic); err != nil {
		return nil, err
	}
	if err := isValidString(configData.Producer.ClientID); err != nil {
		return nil, err
	}
	if err := isValidString(configData.Consumer.GroupID); err != nil {
		return nil, err
	}

	if err := configData.isValidProducerAcks(); err != nil {
		return nil, err
	}

	if configData.Simple.MsgCount <= 0 {
		return nil, errors.New("Error simple.msgCount is not valid")
	}

	if configData.Producer.ValueByteSize < 0 {
		return nil, errors.New("Error value byte size negative")
	}

	return configData, nil
}

func (config *Config) isValidBootstrapServer() error {
	if len(config.BootstrapServer) == 0 {
		return errors.New("bootstrap server list is empty")
	}
	return nil
}

func (config *Config) isValidProducerAcks() error {
	if config.Producer.Acks != "all" && config.Producer.Acks != "0" && config.Producer.Acks != "1" && config.Producer.Acks != "-1" {
		return errors.New("please check Producer.acks! acks value is [all , -1, 0 ,1] only")
	}
	return nil

}

func isValidString(value string) error {
	if value == "" {
		return errors.New("please check config some string is blank")
	}
	return nil
}

func isValidFile(filePath string) error {
	if _, err := os.Stat(filePath); err == nil {
		return nil
	} else if errors.Is(err, os.ErrNotExist) {
		return errors.New(filePath + "file not exists path check please!")
	} else {
		return err
	}
}

func (config *Config) GetBootStrapServer() string {
	var bootstrapServerStr = ""
	for i := 0; i < len(config.BootstrapServer); i++ {
		bootstrapServerStr += config.BootstrapServer[i]
		if i != len(config.BootstrapServer)-1 {
			bootstrapServerStr += ","
		}
	}
	return bootstrapServerStr
}
