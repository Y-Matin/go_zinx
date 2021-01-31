package utils

import (
	"encoding/json"
	"io/ioutil"
	"zinx/src/ziface"
)

/**
存储zinx的所有全局参数，封装为一个全局对象中，暴露给其他模块使用
*/

type GlobalConfig struct {

	/**
	server
	*/
	TcpServer  ziface.Iserver //  当前zinx全局的server对象
	ServerName string         //  当前服务器的名称
	Ip         string         //  当前服务器ip
	Port       int            //  当前服务器端口
	IPVersion  string         //  当前的使用的协议

	/**
	zinx
	*/
	Version        string // 当前zinx的版本号
	MaxConn        int    // 当前服务器主机的最大连接数
	MaxPackageSize uint32 // 当前zinx框架数据包的最大值
}

var Config *GlobalConfig

/**
获取配置对象
*/
func GetGlobalConfig() (config *GlobalConfig) {
	return config
}

/**
初始化init
*/
func init() {
	Config = &GlobalConfig{
		TcpServer:      nil,
		ServerName:     "ZinxServer",
		Ip:             "0.0.0.0",
		Port:           899,
		IPVersion:      "tcp4",
		Version:        "v0.4",
		MaxConn:        1000,
		MaxPackageSize: 4096,
	}
	reloadConfig()
}

/**
reload config from json file
*/

func reloadConfig() {
	file, err := ioutil.ReadFile("config/global.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(file, &Config)
	if err != nil {
		panic(err)
	}
}
