package redis

import "fmt"

func (v Value) Marshal() []byte {
	switch v.typ {
	case STRING:
		return v.marshalString()
	case ERROR:
		return v.marshalError()
	case INTEGER:
		return v.marshalInteger()
	case BULK_STRING:
		return v.marshalBulkString()
	case ARRAY:
		return v.marshalArray()
	}
	return nil
}

func (v Value) marshalString() []byte {
	b := []byte{STRING}
	b = append(b, []byte(v.str)...)
	return append(b, '\r', '\n')
}

func (v Value) marshalError() []byte {
	b := []byte{ERROR}
	b = append(b, []byte(v.bulk)...)
	b = append(b, ' ')
	b = append(b, []byte(v.str)...)
	return append(b, '\r', '\n')
}

func (v Value) marshalInteger() []byte {
	b := []byte{INTEGER}
	if v.num < 0 {
		b = append(b, '-')
	} else {
		b = append(b, '+')
	}
	b = append(b, []byte(fmt.Sprint(v.num))...)
	return append(b, '\r', '\n')
}

func (v Value) marshalBulkString() []byte {
	b := []byte{BULK_STRING}
	b = append(b, []byte(fmt.Sprint(len(v.bulk)))...)
	b = append(b, '\r', '\n')
	b = append(b, []byte(v.bulk)...)
	return append(b, '\r', '\n')
}

func (v Value) marshalArray() []byte {
	b := []byte{ARRAY}
	b = append(b, []byte(fmt.Sprint(len(v.array)))...)
	b = append(b, '\r', '\n')
	for _, v := range v.array {
		b = append(b, v.Marshal()...)
	}
	return b
}
