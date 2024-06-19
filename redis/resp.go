package redis

import (
	"bufio"
	"io"
)

const (
	STRING      = '+'
	ERROR       = '-'
	INTEGER     = ':'
	BULK_STRING = '$'
	ARRAY       = '*'
)

type Value struct {
	typ   rune
	str   string
	num   int
	bulk  string
	array []Value
}

type Resp struct {
	reader *bufio.Reader
}

func NewResp(rd io.Reader) *Resp {
	return &Resp{reader: bufio.NewReader(rd)}
}
