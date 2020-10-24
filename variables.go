package main

import (
	"errors"
	"github.com/go-ini/ini"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

var (
	UserAgent = `Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36`
)

type userPwd struct {
	username string
	password string
}

type Config struct {
	Width   int
	Height  int
	TimeOut int
	Output  string
	Default map[string]string
	Another []map[string]string
}

func (c *Config) Init() {
	s := strings.Split(c.Default["WindowSize"], ",")
	c.TimeOut = strConvInt(c.Default["TimeOut"])
	c.Output = c.Default["Output"]
	if len(s) >= 2 {
		c.Width = strConvInt(s[0])
		c.Height = strConvInt(s[1])
	} else {
		logrus.Errorln(errors.New("width and height cannot be empty"))
	}
}

func ReadConfigFile(path string) *Config {
	cfg, err := ini.Load(path)
	if err != nil {
		logrus.Fatal(err)
	}
	config := &Config{}
	for _, section := range cfg.Sections() {
		if section.Name() == "DEFAULT" {
			config.Default = section.KeysHash()
			continue
		}
		config.Another = append(config.Another, section.KeysHash())
	}
	return config
}

func strConvInt(s string) int {
	res, err := strconv.Atoi(s)
	if err != nil {
		logrus.Fatal(err)
	}
	return res
}
