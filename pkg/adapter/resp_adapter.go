package adapter

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
)

const (
	TypeString  = '*'
	TypeError   = '-'
	TypeInteger = ':'
	TypeBulk    = '$'
	TypeArray   = '*'
)

// typ is used to determine the data type carried by the value.
// str holds the value of the string received from the simple strings.
// num holds the value of the integer received from the integers.
// bulk is used to store the string received from the bulk strings.
// array holds all the values received from the arrays.

type RespDataOutPut struct {
	Typ   string
	Str   string
	Num   int
	Bulk  string
	Value []RespDataOutPut
}

type RespInput struct {
	reader *bufio.Reader
}

func NewRespInput(rd io.Reader) *RespInput {
	return &RespInput{reader: bufio.NewReader(rd)}
}

func (r *RespInput) readLine() (line []byte, n int, err error) {
	for {
		b, err := r.reader.ReadByte()
		if err != nil {
			return nil, 0, fmt.Errorf("error to read line [%w]", err)
		}

		n += 1
		line = append(line, b)

		if n >= 2 && line[n-2] == '\r' {
			break
		}
	}

	return line[:n-2], n, nil
}

func (r *RespInput) readInteger() (x int, n int, err error) {
	line, n, err := r.readLine()

	if err != nil {
		return 0, 0, fmt.Errorf("error to read integer [%w]", err)
	}
	strLine := string(line)
	i64, err := strconv.ParseInt(strLine, 10, 64)

	if err != nil {
		return 0, 0, fmt.Errorf("error to read integer [%w]", err)
	}

	return int(i64), n, nil
}

func (r *RespInput) Read() (RespDataOutPut, error) {
	symbol, err := r.reader.ReadByte()

	if err != nil {
		return RespDataOutPut{}, fmt.Errorf("error to read symbol [%w]", err)
	}

	switch symbol {
	case TypeArray:
		return r.readArray()
	case TypeBulk:
		return r.readBulk()
	default:
		return RespDataOutPut{}, errors.New("unknow type")
	}

}

func (r *RespInput) readArray() (RespDataOutPut, error) {
	d := RespDataOutPut{}
	d.Typ = "array"

	len, _, err := r.readInteger()

	if err != nil {
		return d, err
	}

	d.Value = make([]RespDataOutPut, 0)

	for i := 0; i < len; i++ {
		val, err := r.Read()

		if err != nil {
			return d, err
		}

		d.Value = append(d.Value, val)
	}

	return d, nil
}

func (r *RespInput) readBulk() (RespDataOutPut, error) {
	d := RespDataOutPut{}

	d.Typ = "bulk"

	len, _, err := r.readInteger()
	if err != nil {
		return d, err
	}

	bulk := make([]byte, len)

	r.reader.Read(bulk)

	d.Bulk = string(bulk)

	// Read the trailing CRLF
	r.readLine()

	return d, nil
}

type WriterResp struct {
	writer io.Writer
}

func NewWriter(w io.Writer) *WriterResp {
	return &WriterResp{writer: w}
}

func (w *WriterResp) Write(d RespDataOutPut) error {
	var bytes = d.Marshal()
	_, err := w.writer.Write(bytes)
	if err != nil {
		return err
	}

	return nil
}

func (d RespDataOutPut) Marshal() []byte {
	switch d.Typ {
	case "array":
		return d.marshalArray()
	case "bulk":
		return d.marshalBulk()
	case "string":
		return d.marshalString()
	case "null":
		return d.marshalNull()
	case "error":
		return d.marshalError()
	default:
		return []byte{}
	}
}

func (d RespDataOutPut) marshalArray() []byte {
	var bytes []byte
	len := len(d.Value)

	bytes = append(bytes, TypeArray)
	bytes = append(bytes, strconv.Itoa(len)...)
	bytes = append(bytes, '\r', '\n')

	for i := 0; i < len; i++ {
		bytes = append(bytes, d.Value[i].Marshal()...)
	}

	return bytes
}
func (d RespDataOutPut) marshalBulk() []byte {
	var bytes []byte
	bytes = append(bytes, TypeBulk)
	bytes = append(bytes, strconv.Itoa(len(d.Bulk))...)
	bytes = append(bytes, '\r', '\n')
	bytes = append(bytes, d.Bulk...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}
func (d RespDataOutPut) marshalString() []byte {
	var bytes []byte
	bytes = append(bytes, TypeString)
	bytes = append(bytes, d.Str...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}
func (d RespDataOutPut) marshalError() []byte {
	var bytes []byte
	bytes = append(bytes, TypeError)
	bytes = append(bytes, d.Str...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

func (d RespDataOutPut) marshalNull() []byte {
	return []byte("$-1\r\n")
}
