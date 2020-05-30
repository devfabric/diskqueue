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
	MaxBytesPerFile: 2062336000,
	MinMsgSize:      4,
	MaxMsgSize:      1 << 10,
	SyncEvery:       100,
	SyncTimeout:     2,
	LogLevel:        1,
}

type DiskQueueConfig struct {
	Name            string `json:"name"`            //创建文件和元素据文件的名字
	DataPath        string `json:"dataPath"`        //数据文件保存路径
	MaxBytesPerFile int64  `json:"maxBytesPerFile"` //单个文件最大存储
	MinMsgSize      int32  `json:"minMsgSize"`      //消息最小值
	MaxMsgSize      int32  `json:"maxMsgSize"`      //消息最大值
	SyncEvery       int64  `json:"syncEvery"`       //写入多少条数据后同步写入多少条数据后同步
	SyncTimeout     int64  `json:"syncTimeout"`     //同步时间间隔
	LogLevel        int    `json:"logLevel"`        //日志函数
}

//代码解析
//https://www.cxc233.com/blog/646decf1.html

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
