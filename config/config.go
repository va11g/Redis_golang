package config

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

var Configure *Config

var (
	defaultHost           = "127.0.0.1"
	defaultPort           = 6380
	defaultLogDir         = "./"
	defaultLogLevel       = "info"
	defaultshardNum       = 1024
	defaultChanBufferSize = 10
	configFile            = "./redis.conf"
)

type Config struct {
	ConfFile          string
	Host              string
	Port              int // 用数字便于检验
	LogDir            string
	LogLevel          string
	ShardNum          int
	ChanBufferSize    int
	Databases         int
	Others            map[string]any
	ClusterConfigPath string
	IsCluster         bool
	PeerAddrs         string
	PeerIDs           string
	RaftAddr          string
	NodeID            int
	KVPort            int
	JoinCluster       bool
}

type CfgError struct {
	message string
}

func (cErr *CfgError) Error() string {
	return cErr.message
}

func flagInit(cfg *Config) {
	// 标志的值，标志的名称，标志的默认值，标志的使用说明
	flag.StringVar(&(cfg.ConfFile), "config", configFile, "Appoint a config file: such as /etc/redis.conf")
	flag.StringVar(&(cfg.Host), "host", defaultHost, "Bind host ip: default is 127.0.0.1")
	flag.IntVar(&(cfg.Port), "port", defaultPort, "Bind a listening port: default is 6379")
	flag.StringVar(&(cfg.LogDir), "logdir", defaultLogDir, "Create log directory: default is /tmp")
	flag.StringVar(&(cfg.LogLevel), "loglevel", defaultLogLevel, "Create log level: default is info")
	flag.IntVar(&(cfg.ChanBufferSize), "chanBufSize", defaultChanBufferSize, "set the buffer size of channels in PUB/SUB commands. ")
	// distribution flags
	flag.StringVar(&cfg.ClusterConfigPath, "ClusterConfigPath", "./cluster_config.json", "config file to start cluster mode")
	flag.BoolVar(&cfg.IsCluster, "IsCluster", false, "flag indicates running in cluster mode")
	flag.StringVar(&cfg.PeerAddrs, "PeerAddrs", "http://127.0.0.1:16380", "comma separated cluster peers")
	flag.IntVar(&cfg.NodeID, "NodeID", -1, "node ID")
	flag.IntVar(&cfg.KVPort, "KVPort", 6380, "key-value server port")
	flag.BoolVar(&cfg.JoinCluster, "Join", false, "join an existing cluster")
}

// 初始化并检查config
func SetUp() (*Config, error) {
	cfg := &Config{
		ConfFile:          configFile,
		Host:              defaultHost,
		Port:              defaultPort,
		LogDir:            defaultLogDir,
		LogLevel:          defaultLogLevel,
		ShardNum:          defaultshardNum,
		ChanBufferSize:    defaultChanBufferSize,
		Databases:         16,
		Others:            make(map[string]any),
		ClusterConfigPath: "",
		IsCluster:         false,
		PeerAddrs:         "",
		RaftAddr:          "",
		NodeID:            -1,
		KVPort:            0,
		JoinCluster:       false,
	}
	flagInit(cfg)
	flag.Parse()
	if cfg.ConfFile != "" {
		if err := cfg.Parse(cfg.ConfFile); err != nil {
			return nil, err
		}
	} else {
		if ip := net.ParseIP(cfg.Host); ip == nil {
			ipErr := &CfgError{
				message: fmt.Sprintf("Given ip address %s is invalid", cfg.Host),
			}
			return nil, ipErr
		}
		if cfg.Port <= 1024 || cfg.Port >= 65535 {
			portErr := &CfgError{
				message: fmt.Sprintf("Listening port should between 1024 and 65535, but %d is given.", cfg.Port),
			}
			return nil, portErr
		}
	}
	// cluster mode
	if cfg.IsCluster {
		if cfg.ClusterConfigPath == "" {
			return nil, errors.New("cluster mode need a cluster config file to start. ")
		}
		err := cfg.Parse(cfg.ClusterConfigPath)
		if err != nil {
			return nil, err
		}
	}
	Configure = cfg
	return cfg, nil
}

// 解析redis.conf文件
func (cfg *Config) Parse(cfgFile string) error {
	fl, err := os.Open(cfgFile)
	if err != nil {
		return err
	}
	defer func() error {
		err := fl.Close()
		if err != nil {
			return err
		}
		return nil
	}()

	reader := bufio.NewReader(fl)
	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return err
		}
		// 不要注释
		if len(line) > 0 && line[0] == '#' {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) >= 2 {
			cfgName := strings.ToLower(fields[0])
			switch cfgName {
			case "host":
				if ip := net.ParseIP(fields[1]); ip == nil {
					ipErr := &CfgError{
						message: fmt.Sprintf("Given ip address %s is invalid", cfg.Host),
					}
					return ipErr
				}
				cfg.Host = fields[1]
			case "port":
				port, err := strconv.Atoi(fields[1])
				if err != nil {
					return err
				}
				if port <= 1024 || port >= 65535 {
					portErr := &CfgError{
						message: fmt.Sprintf("Listening port should between 1024 and 65535, but %d is given.", port),
					}
					return portErr
				}
				cfg.Port = port
			case "logdir":
				cfg.LogDir = strings.ToLower(fields[1])
			case "loglevel":
				cfg.LogLevel = strings.ToLower(fields[1])
			case "shardnum":
				cfg.ShardNum, err = strconv.Atoi(fields[1])
				if err != nil {
					fmt.Println("ShardNum should be a number. Get: ", fields[1])
					panic(err)
				}
			case "databases":
				cfg.Databases, err = strconv.Atoi(fields[1])
				if err != nil {
					log.Fatal("Databases should be an integer. Get: ", fields[1])
				}
				if cfg.Databases <= 0 {
					log.Fatal("Databases should be an positive integer. Get: ", fields[1])
				}
			default:
				cfg.Others[cfgName] = fields[1]
			}
		}

		if err == io.EOF {
			break
		}

	}
	return nil
}
