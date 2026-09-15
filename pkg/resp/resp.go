package resp

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

const (
	TypeSimpleString = '+'
	TypeError        = '-'
	TypeInteger      = ':'
	TypeBulkString   = '$'
	TypeArray        = '*'
)

type Value struct {
	Type  byte
	Str   string
	Num   int64
	Bulk  []byte
	Array []Value
}

type Reader struct {
	reader *bufio.Reader
}

func NewReader(rd io.Reader) *Reader {
	return &Reader{reader: bufio.NewReader(rd)}
}

func (r *Reader) ReadValue() (Value, error) {
	b, err := r.reader.ReadByte()
	if err != nil {
		return Value{}, err
	}

	switch b {
	case TypeArray:
		return r.readArray()
	case TypeBulkString:
		return r.readBulk()
	case TypeSimpleString:
		line, err := r.readLine()
		return Value{Type: TypeSimpleString, Str: line}, err
	default:
		// Fallback simple line
		line, err := r.readLine()
		return Value{Type: TypeSimpleString, Str: string(b) + line}, err
	}
}

func (r *Reader) readLine() (string, error) {
	line, err := r.reader.ReadBytes('\n')
	if err != nil {
		return "", err
	}
	n := len(line)
	if n >= 2 && line[n-2] == '\r' {
		return string(line[:n-2]), nil
	}
	return string(line[:n-1]), nil
}

func (r *Reader) readArray() (Value, error) {
	line, err := r.readLine()
	if err != nil {
		return Value{}, err
	}
	count, err := strconv.Atoi(line)
	if err != nil {
		return Value{}, err
	}

	val := Value{Type: TypeArray, Array: make([]Value, count)}
	for i := 0; i < count; i++ {
		elem, err := r.ReadValue()
		if err != nil {
			return Value{}, err
		}
		val.Array[i] = elem
	}
	return val, nil
}

func (r *Reader) readBulk() (Value, error) {
	line, err := r.readLine()
	if err != nil {
		return Value{}, err
	}
	lenBytes, err := strconv.Atoi(line)
	if err != nil {
		return Value{}, err
	}
	if lenBytes == -1 {
		return Value{Type: TypeBulkString, Bulk: nil}, nil
	}

	buf := make([]byte, lenBytes+2) // data + \r\n
	_, err = io.ReadFull(r.reader, buf)
	if err != nil {
		return Value{}, err
	}
	return Value{Type: TypeBulkString, Bulk: buf[:lenBytes]}, nil
}

type Writer struct {
	writer io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{writer: w}
}

func (w *Writer) WriteSimpleString(s string) error {
	_, err := fmt.Fprintf(w.writer, "+%s\r\n", s)
	return err
}

func (w *Writer) WriteError(msg string) error {
	_, err := fmt.Fprintf(w.writer, "-ERR %s\r\n", msg)
	return err
}

func (w *Writer) WriteInteger(i int64) error {
	_, err := fmt.Fprintf(w.writer, ":%d\r\n", i)
	return err
}

func (w *Writer) WriteBulk(b []byte) error {
	if b == nil {
		_, err := w.writer.Write([]byte("$-1\r\n"))
		return err
	}
	_, err := fmt.Fprintf(w.writer, "$%d\r\n%s\r\n", len(b), b)
	return err
}
