package main

import (
	"encoding/json"
	"os"
)

// PGConn : sub struct for parsing config.json
type PGConn struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DbName   string `json:"dbname"`
}

// Config : struct for parsing config.json
type Config struct {
	RootDir  string `json:"root_dir"`
	DataDir  string `json:"data_dir"`
	Database string `json:"database"`
	PGConn   PGConn `json:"pgconn"`
}

// ParseConfig return config info
func ParseConfig() Config {
	f, err := os.Open("./config.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var c Config
	dec := json.NewDecoder(f)
	dec.Decode(&c)
	return c
}
