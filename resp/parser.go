package resp

import (
	"bufio"
	"coding/gxredis/logger"
	"context"
	"io"
)

type ParseRes struct {
	Data RedisData
	Err  error
}

type readState struct {
	bulkLen   int64
	arrayLen  int
	multLine  bool
	arrayData *ArrayData
	inArray   bool
}

func parse(ctx context.Context, reader io.Reader, ch chan *ParseRes) {
	bufReader := bufio.NewReaderSize(reader, 2048)
	state := new(readState)
	for {
		var res RedisData
		msg, err := readLine(bufReader, state)
		if err != nil {
			if err == io.EOF { //读取结束
				ch <- &ParseRes{
					Err: io.EOF,
				}
				close(ch)
				return
			} else {
				if ctx.Err() != nil { //客户端关闭
					close(ch)
					return
				}
				logger.Error(err)
				ch <- &ParseRes{
					Err: err,
				}
				*state = readState{}
			}
			continue
		}

		if !state.multLine {
			switch msg[0] {
			case '*':
				err := parseArrayHeader(msg, state)
				if err != nil {
					logger.Error(err)
					ch <- &ParseRes{
						Err: err,
					}
					*state = readState{}
				} else {
					if state.arrayLen == -1 { //NULL
						ch <- &ParseRes{
							Data: MakeArrayData(nil),
						}
					} else if state.arrayLen == 0 {
						ch <- &ParseRes{
							Data: MakeArrayData([]RedisData{}),
						}
						*state = readState{}
					}
				}
			case '$':
				err := parseArrayHeader(msg, state)
				if err != nil {
					logger.Error(err)
					ch <- &ParseRes{
						Err: err,
					}
					*state = readState{}
				} else {
					if state.arrayLen == -1 { //NULL
						state.multLine = false
						state.bulkLen = 0

					}
				}
			default:
				res, err = parseSingleLine(msg)
			}
		} else {
			res, err = ParseMultLine(msg)
		}
	}
}

func ParseStream(ctx context.Context, reader io.Reader) <-chan *ParseRes {
	ch := make(chan *ParseRes)
	go parse(ctx, reader, ch)
	return ch
}
