package config

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

var QueueConfig = &DiskQueueConfig{
	Name:            "disk_queue",
	DataPath:        filepath.Join("./", "data"),
	MaxBytesPerFile: 1024,
	MinMsgSize:      4,
	MaxMsgSize:      1 << 10,
	SyncEvery:       2500,
	SyncTimeout:     2,
	LogLevel:        1,
}

type DiskQueueConfig struct {
	Name            string `json:"name"`
	DataPath        string `json:"dataPath"`
	MaxBytesPerFile int64  `json:"maxBytesPerFile"`
	MinMsgSize      int32  `json:"minMsgSize"`
	MaxMsgSize      int32  `json:"maxMsgSize"`
	SyncEvery       int64  `json:"syncEvery"`
	SyncTimeout     int64  `json:"syncTimeout"`
	LogLevel        int    `json:"logLevel"`
}

func CheckFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func LoadQueueConfig(dir string) (*DiskQueueConfig, error) {
	path := filepath.Join(dir, "./configs/diskqueue.toml")
	filePath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	config := new(DiskQueueConfig)
	if CheckFileIsExist(filePath) { //文件存在
		if _, err := toml.DecodeFile(filePath, config); err != nil {
			return nil, err
		} else {
			QueueConfig = config
		}
	} else {
		configBuf := new(bytes.Buffer)
		if err := toml.NewEncoder(configBuf).Encode(QueueConfig); err != nil {
			return nil, err
		}
		err := ioutil.WriteFile(filePath, configBuf.Bytes(), 0666)
		if err != nil {
			return nil, err
		}
	}
	return config, nil
}
