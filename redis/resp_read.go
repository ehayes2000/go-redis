package redis

import (
	"fmt"
	"strconv"
	"unicode"
)

/*
	redis protocol
	request []string
	request[i] = command and arguments that the server should execute

	Parts are seperatedb by CRLF (\r\n)



	The first byte of data determines type
	then is the size of data
	RESP data is simple, bulk, or aggregate


	ex:
		"+OK\r\n"
		"-ERR unknown command 'asdf'
		"-WRONGTYPE Operation against a key holding the wrong kind of value"
		"$5\r\nhello\r\n"
*/

func (r *Resp) Read() (Value, error) {
	b, e := r.reader.ReadByte()
	if e != nil {
		return Value{}, e
	}
	switch b {
	case STRING:
		return r.parseRespString()
	case ERROR:
		return r.parseRespErr()
	case INTEGER:
		return r.parseRespInt()
	case BULK_STRING:
		return r.parseRespBulk()
	case ARRAY:
		return r.parseRespArray()
	}
	return Value{}, nil
}

func (r *Resp) parseRespInt() (Value, error) {
	b, e := r.reader.Peek(1)
	var minus = false
	if e != nil {
		return Value{}, e
	}
	if b[0] == byte('-') {
		minus = true
	}
	if b[0] == byte('-') || b[0] == byte('+') {
		r.reader.Discard(1)
	}
	num, e := r.reader.ReadBytes('\n')
	if e != nil {
		return Value{}, e
	}
	intNum, e := strconv.ParseInt(string(num[:len(num)-2]), 10, 32)
	if e != nil {
		return Value{}, e
	}
	if minus {
		return Value{
			typ: INTEGER,
			num: -int(intNum),
		}, nil
	}
	r.reader.Discard(2)
	return Value{
		typ: INTEGER,
		num: int(intNum),
	}, nil

}

func (r *Resp) parseRespErr() (Value, error) {
	var errType []byte
	for {
		b, e := r.reader.ReadByte()
		if e != nil {
			return Value{}, e
		}
		if unicode.IsSpace(rune(b)) {
			break
		}
		errType = append(errType, b)
	}
	msg, e := r.reader.ReadBytes('\r')
	if e != nil {
		return Value{}, e
	}
	fmt.Printf("\n\nMSG ERR: %v\n\n", string(msg))
	return Value{
		typ:  ERROR,
		bulk: string(errType),
		str:  string(msg[:len(msg)-1]),
	}, nil
}

func (r *Resp) parseRespString() (Value, error) {
	val, e := r.reader.ReadBytes('\n')
	if e != nil {
		return Value{}, e
	}
	r.reader.Discard(2)
	return Value{
		typ: STRING,
		str: string(val[:len(val)-2]),
	}, nil
}

func (r *Resp) parseRespBulk() (Value, error) {
	size, _ := r.reader.ReadBytes('\n')
	strSize, e := strconv.Atoi(string(size[:len(size)-2]))
	if e != nil {
		return Value{}, error(fmt.Errorf("could not parse bulk string size %e", e))
	}
	bulk := make([]byte, strSize)
	_, e = r.reader.Read(bulk)
	if e != nil {
		return Value{}, error(fmt.Errorf("could not read bulk string %e", e))
	}
	r.reader.Discard(2)
	return Value{
		typ:  BULK_STRING,
		bulk: string(bulk),
	}, nil
}

func (r *Resp) parseRespArray() (Value, error) {
	b, e := r.reader.ReadBytes('\n')
	if e != nil {
		return Value{}, e
	}
	nElems, e := strconv.Atoi(string(b[:len(b)-2]))
	if e != nil {
		return Value{}, e
	}
	var items []Value
	for i := 0; i < nElems; i++ {

		fmt.Println(i)
		v, e := r.Read()

		fmt.Printf("HEHEXD %d - %v\n", i, v)
		if e != nil {
			return Value{}, e
		}
		items = append(items, v)
	}
	return Value{
		typ:   ARRAY,
		array: items,
	}, nil
}
