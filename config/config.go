package config

import (
	"encoding/json"
	"os"
)

type database struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
	Host     string `json:"host"`
}

type Config struct {
	Db database `json:"database"`
}

func ParseConfig() (c *Config, err error) {
	f, err := os.Open("./config/config.json")
	if err != nil {
		return
	}
	c = new(Config)
	err = json.NewDecoder(f).Decode(c)
	return
}
