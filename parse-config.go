package main

import (
	"encoding/json"
	"os"
)

// Data : sub struct for parsing config.json
type Data struct {
	DirName    string `json:"dir_name"`
	SourceFile string `json:"source_file"`
	ImportFile string `json:"import_file"`
	ExportFile string `json:"export_file"`
	VectisFile string `json:"vectis_file"`
}

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
	Data     Data   `json:"data"`
	Database string `json:"database"`
	PGConn   PGConn `json:"pgconn"`
	SqliteDb string `json:"sqlite_db"`
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
