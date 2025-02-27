package server

import (
	"context"
	"io"
	"net"
	"resp/config"
	"resp/logger"
	"resp/resp"
)

type Manager struct {
	CurrentDB *memdb.memdb
	DBs       []*memdb.memdb
}

type MemStorageStats struct {
	initialState, firstIndex, lastIndex, entries, term, snapshot int
}

func NewManger(cfg *config.Config) *Manager {
	DBs := make([]*memdb.memdb, cfg.Databases)
	for i := 0; i < cfg.Databases; i++ {
		DBs[i] = memdb.NewMemDb()
	}
	return &Manager{
		CurrentDB: DBs[0],
		DBs:       DBs,
	}
}

func (m *Manager) Handle(ctx context.Context, conn net.Conn) {
	defer func() {
		err := conn.Close()
		if err != nil {
			logger.Error(err)
		}
	}()

	ch := resp.ParseStream(ctx, conn)

	for {
		select {
		case parsedRes := <-ch:
			if parsedRes.Err != nil {
				if parsedRes.Err == io.EOF {
					logger.Info("Close connection ", conn.RemoteAddr().String())
				} else {
					logger.Panic("Handle connection ", conn.RemoteAddr().String(), " panic: ", parsedRes.Err.Error())
				}
				return
			}

		}
	}
}
