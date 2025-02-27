package server

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"resp/config"
	"resp/logger"
	"strconv"
	"sync"
	"syscall"
)

func Start(cfg *config.Config) error {
	listener, err := net.Listen("tcp", cfg.Host+":"+strconv.Itoa(cfg.Port))
	if err != nil {
		logger.Panic(err)
		return err
	}
	var isTerminating bool
	wg := sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		logger.Info("Shutting down gracefully")
		err := listener.Close()
		if err != nil {
			logger.Error(err)
		}
		cancel()
		wg.Wait()
		logger.Info("Bye")
	}()

	fmt.Println("Hello")
	logger.Info("Server Listen At ", cfg.Host, ":", cfg.Port)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)
	clients := make(chan net.Conn)
	mgr := NewManger(cfg)

	go func() {
		for {
			conn, err = listener.Accept()
			if isTerminating {
				close(clients)
				return
			}
			if err != nil {
				logger.Error(err)
				return
			}
			clients <- conn
		}
	}()
	return nil
}
