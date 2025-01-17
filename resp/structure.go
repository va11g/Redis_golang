package resp

import (
	"fmt"
	"strconv"
	"strings"
)

var (
	CRLF = "\r\n"
	NIL  = []byte("$-1\r\n")
)

type RedisData interface {
	ToBytes() []byte //返回可供解析的标准格式
	ByteData() []byte
	String() string
}

type StringData struct { //+
	data string
}

type BulkData struct { //$
	data []byte
}

type IntData struct { //:
	data int64
}

type ErrorData struct { //-
	data string
}

type ArrayData struct { //*
	data []RedisData
}

type PlainData struct {
	data string
}

// ////////////////////////////////////////////////////////////////////////////////////////////////////StringData
func MakeStringData(data string) *StringData {
	return &StringData{
		data: data,
	}
}

func (r *StringData) Data() string {
	return r.data
}

func (r *StringData) ToBytes() []byte {
	return []byte("+" + r.data + CRLF)
}

func (r *StringData) ByteData() []byte {
	return []byte(r.data)
}

func (r *StringData) String() string {
	return r.data
}

// ////////////////////////////////////////////////////////////////////////////////////////////////////BulkData
func MakeBulkData(data []byte) *BulkData {
	return &BulkData{
		data: data,
	}
}

func (r *BulkData) Data() []byte {
	return r.data
}

func (r *BulkData) ToBytes() []byte {
	if r.data == nil {
		return NIL
	}
	return []byte("$" + strconv.Itoa(len(r.data)) + CRLF + string(r.data) + CRLF)
}

func (r *BulkData) ByteData() []byte {
	return r.data
}

func (r *BulkData) String() string {
	return string(r.data)
}

// ////////////////////////////////////////////////////////////////////////////////////////////////////IntData
func MakeIntData(data int64) *IntData {
	return &IntData{
		data: data,
	}
}

func (r *IntData) Data() int64 {
	return r.data
}

func (r *IntData) ToBytes() []byte {
	if r.data == 0 {
		return NIL
	}
	return []byte(":" + strconv.FormatInt(r.data, 10) + CRLF)
}

func (r *IntData) ByteData() []byte {
	return []byte(strconv.FormatInt(r.data, 10))
}

func (r *IntData) String() string {
	return strconv.FormatInt(r.data, 10)
}

// ////////////////////////////////////////////////////////////////////////////////////////////////////ErrorData
func MakeErrorData(data ...string) *ErrorData {
	errMsg := ""
	for _, v := range data {
		errMsg += v
	}
	return &ErrorData{
		data: errMsg,
	}
}

func (r *ErrorData) Data() string {
	return r.data
}

func (r *ErrorData) ToBytes() []byte {
	return []byte("-" + r.data + CRLF)
}

func (r *ErrorData) ByteData() []byte {
	return []byte(r.data)
}

func (r *ErrorData) String() string {
	return r.data
}

func MakeWrongNumberArgs(name string) *ErrorData {
	return &ErrorData{data: fmt.Sprintf("Ero wrong number of arguments for %s command", name)}
}

func MakeWrongType() *ErrorData {
	return &ErrorData{data: "Wrong type operation aginst a key holding the wrong kind of value"}
}

// ////////////////////////////////////////////////////////////////////////////////////////////////////ArrayData
func MakeArrayData(data []RedisData) *ArrayData {
	return &ArrayData{
		data: data,
	}
}

func MakeEmptyArrayData() *ArrayData {
	return &ArrayData{
		data: []RedisData{},
	}
}

func (r *ArrayData) Data() []RedisData {
	return r.data
}

func (r *ArrayData) ToBytes() []byte {
	if r.data == nil {
		return []byte("*-1\r\n")
	}
	res := []byte("*" + strconv.Itoa(len(r.data)) + CRLF)
	for _, v := range r.data {
		res = append(res, v.ToBytes()...)
	}
	return res
}

func (r *ArrayData) ByteData() []byte {
	res := make([]byte, 0)
	for _, v := range r.data {
		res = append(res, v.ByteData()...)
	}
	return res
}

func (r *ArrayData) ByteDataArr() [][]byte {
	res := make([][]byte, 0)
	for _, v := range r.data {
		res = append(res, v.ByteData())
	}
	return res
}

func (r *ArrayData) StringArr() []string {
	res := make([]string, 0, len(r.data))
	for _, v := range r.data {
		res = append(res, v.String())
	}
	return res
}

func (r *ArrayData) String() string {
	return strings.Join(r.StringArr(), " ")
}

// ////////////////////////////////////////////////////////////////////////////////////////////////////PlainData
func MakePlainData(data string) *PlainData {
	return &PlainData{
		data: data,
	}
}

func (r *PlainData) Data() string {
	return r.data
}

func (r *PlainData) ToBytes() []byte {
	return []byte(r.data + CRLF)
}

func (r *PlainData) ByteData() []byte {
	return []byte(r.data)
}

func (r *PlainData) String() string {
	return r.data
}
