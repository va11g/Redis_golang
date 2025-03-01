package memdb

import (
	"context"
	"net"
	"resp/resp"
	"strings"
)

type cmdBytes = [][]byte
type cmdExecutor func(ctx context.Context, m *MemDb, cmd cmdBytes, conn net.Conn) resp.RedisData

var CmdTable = make(map[string]*cmdExecutor)

// 注册命令
func RegisterCommand(cmdName string, executor cmdExecutor) {
	CmdTable[cmdName] = &executor
}

// 识别命令
func MakeCommandBytes(input string) cmdBytes {
	cmdStrs := strings.Split(input, " ")
	cmds := make(cmdBytes, 0)
	for _, c := range cmdStrs {
		cmds = append(cmds, []byte(c))
	}
	return cmds
}
