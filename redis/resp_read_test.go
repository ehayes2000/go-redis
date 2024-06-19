package redis

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

var tests = []struct {
	r    string
	want Value
	err  error
}{
	{
		r: "+OK\r\n",
		want: Value{
			typ: STRING,
			str: "OK",
		},
	},
	{
		r: "$5\r\nhello\r\n",
		want: Value{typ: BULK_STRING,
			bulk: "hello",
		},
	},
	{
		r: "$11\r\nhello world\r\n",
		want: Value{typ: BULK_STRING,
			bulk: "hello world",
		},
	},
	{
		r: "-E myerr\r\n",
		want: Value{typ: ERROR,
			bulk: "E",
			str:  "myerr",
		},
	},
	{
		r: "-ERR unknown command 'asdf'\r\n",
		want: Value{typ: ERROR,
			bulk: "ERR",
			str:  "unknown command 'asdf'",
		},
	},
	{
		r: "-WRONGTYPE Operation against a key holding the wrong kind of value\r\n",
		want: Value{typ: ERROR,
			bulk: "WRONGTYPE",
			str:  "Operation against a key holding the wrong kind of value",
		},
	},
	{
		r: "*0\r\n",
		want: Value{
			typ: ARRAY,
		},
	},
	{
		r: "*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n",
		want: Value{
			typ: ARRAY,
			array: []Value{
				{
					typ:  BULK_STRING,
					bulk: "hello",
				},
				{
					typ:  BULK_STRING,
					bulk: "world",
				},
			},
		},
	},
}

func TestRespRead(t *testing.T) {
	for _, test := range tests {
		testname := fmt.Sprintf("[%s]", test.r)
		t.Run(testname, func(t *testing.T) {
			resp := NewResp(strings.NewReader(test.r))
			res, e := resp.Read()
			if e != nil && test.err == nil {
				t.Errorf("unexpected error %v\n\n", e)
			} else if e == nil && test.err != nil {
				t.Errorf("missing expected error %v\nn", res)
			} else if !reflect.DeepEqual(res, test.want) {
				t.Errorf("result mismatch \nwant %v\n\ngot %v\n\n", test.want, res)
			}
		})
	}

}

func TestRespWrite(t *testing.T) {
	for _, test := range tests {
		testname := fmt.Sprintf("[%s]", test.r)
		t.Run(testname, func(t *testing.T) {
			b := test.want.Marshal()
			if !reflect.DeepEqual(b, []byte(test.r)) {
				t.Errorf("Unexpected serialization\nwant[%v]\ngot [%v]", []byte(test.r), b)
			}
		})
	}

}
