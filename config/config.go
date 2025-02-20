package config

var Configure *Config

var (
	defaultHost           = "127.0.0.1"
	dafaultPort           = 6380
	dafaultLogDir         = "./"
	dafaultLogLevel       = "info"
	dafaultshardNum       = 1024
	dafaultChanBufferSize = 10
	configFile            = "./redis.conf"
)

type Config struct {
	ConfFile       string
	Host           string
	port           string
	LogDir         string
	LogLevel       string
	ShardNUm       int
	ChanBufferSize int
	DataBases      int
	Others         map[string]any
	IsCluster      bool
	PeerAddrs      string
	PeerIDs        string
	RaftAddr       string
	NodeID         int
	KVPort         int
	JoinCluster    bool
}

type CfgError struct {
	message string
}

func (cErr *CfgError) Error() string {
	return cErr.message
}
