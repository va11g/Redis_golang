package resp

import (
	"bufio"
	"coding/gxredis/logger"
	"context"
	"errors"
	"io"
	"strconv"
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

func ParseStream(ctx context.Context, reader io.Reader) <-chan *ParseRes { //外部调用，持续接收
	ch := make(chan *ParseRes)
	go parse(ctx, reader, ch)
	return ch
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

		if !state.multiLine {
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
						state.multiLine = false
						state.bulkLen = 0

					}
				}
			default:
				res, err = parseSingleLine(msg)
			}
		} else {
			res, err = parseMultiLine(msg)
		}
	}
}

func parseSingleLine(msg []byte) (RedisData, error) {
	msgType := msg[0]
	msgData := string(msg[1 : len(msg)-2])
	var res RedisData
	if len(msg) < 3 {
		return nil, errors.New("msg too short, possibly due to a http connection to redis port. ")
	}
	switch msgType {
	case '+':
		res = MakeStringData(msgData)
	case '-':
		res = MakeErrorData(msgData)
	case ':':
		data, err := strconv.ParseInt(msgData, 10, 64)
		if err != nil {
			logger.Error("Cant phrase int64 from " + msgData + " where error: " + string(msg))
			return nil, err
		}
		res = MakeIntData(data)
	default:
		res = MakePlainData(msgData)
	}
	return res, nil
}

func parseMultiLine(msg []byte) (RedisData, error) {
	if len(msg) < 2 {
		return nil, errors.New("protocol error: invalid bulk string")
	}
	msgData := msg[:len(msg)-2]
	res := MakeBulkData(msgData)
	return res, nil
}

func parseBulkHeader(msg []byte, state *readState) error {
	bulkLen, err = strconv.ParseInt()
	return
}

func parseArrayHeader(msg []byte, state *readState) error {
	return
}

func readLine(reader *bufio.Reader, state *readState) ([]byte, error) {
	var msg []byte
	var err error
	if state.multiLine && state.bulkLen >= 0 {

	}
	return
}
