package utils

import (
	"io"
	"math"
	"sort"
	"strconv"
)

// strip remove leaders, padding, and exponents
func strip(dst io.Writer, buff []byte) (int, error) {
	var exponent int

	for byte_ := range buff {
		if buff[byte_] == 'E' {
			exponent = byte_
		}
	}
	var start = exponent + 1
	var offset = 0
	var plus = buff[exponent+1] == '+'

	if plus {
		offset++
		start = exponent + 2
	}
	padded := buff[exponent+2] == '0'
	if padded {
		offset++
		start = exponent + 3
	}
	for start < len(buff) {
		buff[start-offset] = buff[start]
		start++
	}
	buff = buff[:len(buff)-offset]
	return dst.Write(buff)
}

// NumOut - write out a number
func NumOut(dst io.Writer, f float64) (int, error) {
	if 1-(1<<53) <= f && f <= (1<<53)-1 {
		_, frac := math.Modf(f)
		if frac == 0.0 {
			return io.WriteString(dst, strconv.FormatInt(int64(f), 10))
		}
	}
	bufBytes := strconv.AppendFloat([]byte{}, f, 'E', -1, 64)
	return strip(dst, bufBytes)
}

// Sort with custom comparator
// https://yourbasic.org/golang/how-to-sort-in-go/
type tuple struct {
	key string
	val []byte
}

type tuples []tuple

func (t tuples) Len() int {
	return len(t)
}

func (t tuples) Less(i, j int) bool {
	return t[i].key < t[j].key
}

func (t tuples) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

// -^

// ObjectOut - write out an object
func (canon *Inputs) ObjectOut(out io.Writer, values tuples) (int64, error) {
	var contents int64
	w, err := out.Write([]byte{'{'})
	contents += int64(w)
	if err != nil {
		return contents, err
	}
	sort.Sort(values)
	first := true
	for _, value := range values {
		if !first {
			w, err = out.Write([]byte{','})
			contents += int64(w)
			if err != nil {
				return contents, err
			}
		}
		first = false
		w64, err := canon.StringOut(out, value.key)
		contents += w64
		if err != nil {
			return contents, err
		}
		_, err = out.Write([]byte{':'})
		if err != nil {
			return 0, err
		}
		w, err = out.Write(value.val)
		contents += int64(w)
		if err != nil {
			return contents, err
		}
	}
	w, err = out.Write([]byte{'}'})
	contents += int64(w)
	return contents, err
}

var Exceptions = []byte{
	0, 0, 0, 0, 0, 0, 0, 0, 'b', 't', 'n', 0, 'f', 'r', 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
}

/**
MUST represent all strings (including object member names) in their minimal-length UTF-8 encoding

    avoiding escape sequences for characters except those otherwise inexpressible in
    JSON (U+0022 QUOTATION MARK, U+005C REVERSE SOLIDUS, and ASCII control
    characters U+0000 through U+001F) or UTF-8 (U+D800 through U+DFFF), and avoiding
    escape sequences for combining characters, variation selectors, and other code
    points that affect preceding characters, and using two-character escape
    sequences where possible for characters that require escaping:
        \b U+0008 BACKSPACE
        \t U+0009 CHARACTER TABULATION (“tab”)
        \n U+000A LINE FEED (“newline”)
        \f U+000C FORM FEED
        \r U+000D CARRIAGE RETURN
        \" U+0022 QUOTATION MARK
        \\ U+005C REVERSE SOLIDUS (“backslash”), and
    using six-character \u00xx uppercase hexadecimal escape sequences for control
    characters that require escaping but lack a two-character sequence, and using
    six-character \uDxxx uppercase hexadecimal escape sequences for lone surrogates

*/

const HEX = "0123456789abcdef"

// StringOut - write string out
func (canon *Inputs) StringOut(out io.Writer, s string) (int64, error) {
	canon.buffer.Reset()
	canon.buffer.WriteByte('"')
	for _, v := range s {
		if v < 0x20 {
			ex := Exceptions[v]
			if ex != 0 {
				canon.buffer.WriteByte('\\')
				canon.buffer.WriteByte(ex)
			} else {
				canon.buffer.Write([]byte(`\u00`))
				canon.buffer.WriteByte(HEX[v>>4])
				canon.buffer.WriteByte(HEX[v&0x0f])
			}
		} else {
			if v == '\\' || v == '"' {
				canon.buffer.WriteByte('\\')
			}
			canon.buffer.WriteRune(v)
		}
	}
	canon.buffer.WriteByte('"')
	return canon.buffer.WriteTo(out)
}
