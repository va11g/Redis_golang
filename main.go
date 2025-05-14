package main

import (
	"os"
	"resp/logger"
	"resp/server"
)

func init() {
	// Register commands
	memdb.RegisterKeyCommands()
	memdb.RegisterStringCommands()
	memdb.RegisterListCommands()
	memdb.RegisterSetCommands()
	memdb.RegisterHashCommands()
	memdb.RegisterPubSubCommands()
	memdb.RegisterSortedSetCommands()
	memdb.RegisterStreamCommands()
	memdb.RegisterRaftCommand()
}

func main() {
	// setup config
	cfg, err := config.Setup()
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}
	// setup logger
	err = logger.SetUp(cfg)
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}
	// setup tcp server
	err = server.Start(cfg)
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}
}
