package memdb

import (
	"context"
	"net"
	"resp/resp"
	"strings"
)

func setString(ctx context.Context, m *MemDb, cmd [][]byte, conn net.Conn) resp.RedisData {
	if len(cmd) < 3 {
		return resp.MakeErrorData("error: commands is invalid")
	}

	for i := 3; i < len(cmd); i++ {
		switch strings.ToLower(string(cmd[i]))
	}
}
