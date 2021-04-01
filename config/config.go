package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"sync"
)

var once sync.Once
var config = &Config{}

type Config struct {
	IsRandomUserAgent bool   `yaml:"isRandomUseragent"`
	EID               string `yaml:"eid"`
	FP                string `yaml:"fp"`
	TrackID           string `yaml:"track_id"`
	RiskControl       string `yaml:"risk_control"`
	Timeout           int64  `yaml:"timeout"`
	EnableSendMessage bool   `yaml:"enableSendMessage"`
	Messenger         struct {
		Sckey string `yaml:"sckey"`
	} `yaml:"messenger"`
}

func Get() *Config {
	once.Do(func() {
		conf, err := ioutil.ReadFile("./conf/conf.yaml")
		if err != nil {
			log.Fatalln("read conf file: ", err)
		}
		err = yaml.Unmarshal(conf, config)
		if err != nil {
			log.Fatalln("unmarshal conf file: ", err)
		}
	})
	return config
}
