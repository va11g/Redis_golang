package resp

import (
	"bufio"
	"context"
	"io"
)

type ParseRes struct { //读取结果
	Data RedisData
	Err  error
}

type readState struct { //当前读取状态
	bulkLen   int64
	arrayLen  int
	multiLine bool
	arrayData *ArrayData
	inArray   bool
}

func ParseStream(ctx context.Context, reader io.Reader) <-chan *ParseRes {
	ch := make(chan *ParseRes)
	go parse(ctx, reader, ch)
	return ch
}

func parse(ctx context.Context, reader io.Reader, ch chan *ParseRes) {
	bufReader := bufio.NewReaderSize(reader, 2048)
	state := new(readState)
}

func readLine(reader *bufio.Reader, state *readState) ([]byte, error) {
	var msg []byte
	var err error
	if state.multiLine && state.bulkLen >= 0 {
		msg = make([]byte, state.bulkLen+2)
	}
}
