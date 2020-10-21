package service

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/sipt/GoJsoner"
	"github.com/tang-go/go-dog/lib/net"
	"github.com/tang-go/go-dog/pkg/config"
)

var (
	configpath string
)

func init() {
	flag.StringVar(&configpath, "c", "./config/config.json", "config配置路径")
}

//Config 配置
type Config struct {
	//服务名称
	ServerName string `json:"server_name"`
	//服务说明
	Explain string `json:"explain"`
	//使用端口号
	Port int `json:"port"`
	//GossipPort端口
	GossipPort int `json:"gossip_port"`
	//members Goosip成员
	Members []string `json:"members"`
	//Discovery 服务发现
	Discovery []string `json:"discovery"`
	//Redis地址
	Redis []string `json:"redis"`
	//Kafka地址
	Kafka []string `json:"kafka"`
	//Nats地址
	Nats []string `json:"nats"`
	//RocketMq地址
	RocketMq []string `json:"rocket_mq"`
	//nsq地址
	Nsq []string `json:"nsq"`
	//Jaeger 链路追踪地址
	Jaeger string `json:"jaeger"`
	//读数据库地址
	ReadMysql *config.MysqlCfg `json:"read_mysql"`
	//写数据库地址
	WriteMysql *config.MysqlCfg `json:"write_mysql"`
	//本机地址
	Host string `json:"host"`
	//运行日志等级 panic fatal error warn info debug trace
	Runmode string `json:"runmode"`
}

//GetServerName 获取服务名称
func (c *Config) GetServerName() string {
	return c.ServerName
}

//GetPort 获取端口
func (c *Config) GetPort() int {
	return c.Port
}

//GetGossipPort 获取gossipport端口
func (c *Config) GetGossipPort() int {
	return c.GossipPort
}

//GetMembers 获取gossip会有
func (c *Config) GetMembers() []string {
	return c.Members
}

//GetExplain 获取服务说明
func (c *Config) GetExplain() string {
	return c.Explain
}

//GetDiscovery 获取服务发现配置
func (c *Config) GetDiscovery() []string {
	return c.Discovery
}

//GetRedis 获取redis配置
func (c *Config) GetRedis() []string {
	return c.Redis
}

//GetKafka 获取kfaka地址
func (c *Config) GetKafka() []string {
	return c.Kafka
}

//GetNats 获取nats地址
func (c *Config) GetNats() []string {
	return c.Nats
}

//GetRocketMq 获取rocketmq地址
func (c *Config) GetRocketMq() []string {
	return c.RocketMq
}

//GetNsq 获取nsq地址
func (c *Config) GetNsq() []string {
	return c.Nsq
}

//GetReadMysql 获取ReadMysql地址
func (c *Config) GetReadMysql() *config.MysqlCfg {
	return c.ReadMysql
}

//GetWriteMysql 获取GetWriteMysql地址
func (c *Config) GetWriteMysql() *config.MysqlCfg {
	return c.WriteMysql
}

//GetHost 获取本机地址配置
func (c *Config) GetHost() string {
	return c.Host
}

//GetRunmode 获取runmode地址配置
func (c *Config) GetRunmode() string {
	return c.Runmode
}

//NewConfig 初始化Config
func NewConfig() *Config {
	c := new(Config)
	//从文件读取json文件并且解析
	flag.Parse()
	s := os.Getenv("config")
	if s == "" {
		gameConfigData, err := ioutil.ReadFile(configpath)
		if err != nil {
			panic(err.Error())
		}
		gameConfigResult, err := GoJsoner.Discard(string(gameConfigData))
		if err != nil {
			panic(err.Error())
		}
		err = json.Unmarshal([]byte(gameConfigResult), c)
		if err != nil {
			panic(err.Error())
		}
	} else {
		gameConfigResult, err := GoJsoner.Discard(s)
		if err != nil {
			panic(err.Error())
		}
		err = json.Unmarshal([]byte(gameConfigResult), c)
		if err != nil {
			panic(err.Error())
		}
	}

	host := os.Getenv("HOST")
	if host != "" {
		c.Host = host
	} else {
		if c.Host == "" {
			host, err := net.ExternalIP()
			if err != nil {
				panic(err.Error())
			}
			c.Host = host.String()
		}
	}
	//先看环境变量是否有端口号
	port := os.Getenv("PORT")
	if port != "" {
		p, err := strconv.Atoi(port)
		if err != nil {
			c.Port = p
		}
	}
	if c.Port <= 0 {
		//获取随机端口
		p, err := net.GetFreePort()
		if err != nil {
			panic(err.Error())
		}
		c.Port = p
	}
	//先看环境变量是否有端口号
	gossipport := os.Getenv("GOSSIPPORT")
	if gossipport != "" {
		p, err := strconv.Atoi(gossipport)
		if err != nil {
			c.GossipPort = p
		}
	}
	if c.GossipPort <= 0 {
		//获取随机端口
		p, err := net.GetFreePort()
		if err != nil {
			panic(err.Error())
		}
		c.GossipPort = p
	}
	//先看环境变量是否有端口号
	members := os.Getenv("MERBERS")
	if members != "" {
		c.Members = strings.Split(members, ",")
	}
	//设置运行模式
	runmode := os.Getenv("RUNMODE")
	if runmode != "" {
		c.Runmode = runmode
	}
	fmt.Println("************************************************")
	fmt.Println("*                                              *")
	fmt.Println("*             	   Cfg  Init                    *")
	fmt.Println("*                                              *")
	fmt.Println("************************************************")
	fmt.Println("### Port:         ", c.Port)
	fmt.Println("### GossipPort:   ", c.GossipPort)
	fmt.Println("### Members:      ", c.Members)
	fmt.Println("### Host:         ", c.Host)
	fmt.Println("### RunMode:      ", c.Runmode)
	return c
}
