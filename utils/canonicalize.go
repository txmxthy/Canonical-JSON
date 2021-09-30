package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

// Canonify - Fit input Json to Inputs requirements and write to output (in, out).
// - Multiple inputs separated by space
func Canonify(out io.Writer, in io.Reader) (int64, error) {
	canon := &Inputs{
		Decoder: json.NewDecoder(in),
	}
	var processed int64
	var err error
	for {
		var target int64
		target, err = canon.Process(out)
		if err != nil {
			break
		}
		processed += target
	}
	if err == io.EOF {
		err = nil
	}
	return processed, err
}

type Inputs struct {
	buffer  bytes.Buffer
	Decoder *json.Decoder
	Spaced  bool
}

// Process - Process Json with respective types' helper methods
func (canon *Inputs) Process(out io.Writer) (int64, error) {
	token, err := canon.Decoder.Token()
	if err != nil {
		return 0, err
	}
	switch t := token.(type) {
	case string:
		canon.Spaced = false
		return canon.StringOut(out, t)
	case float64:
		return canon.f64(out, t)
	case json.Delim:
		switch t {
		case '[':
			canon.Spaced = false
			return canon.arr(out)
		case '{':
			canon.Spaced = false
			return canon.obj(out)
		}
	case bool:
		return canon.boole(out, t, err)
	default:
		i, err, contents := canon.others(out, t)
		if contents {
			return i, err
		}
	}
	panic(fmt.Sprintf("unexpected value found in input %v", token))
}

// others - default json conversion processing method
func (canon *Inputs) others(out io.Writer, t json.Token) (int64, error, bool) {
	if t == nil {
		var contents int64
		if canon.Spaced {
			current, err := out.Write([]byte{' '})
			contents += int64(current)
			if err != nil {
				return contents, err, true
			}
		}
		canon.Spaced = true
		w, err := out.Write([]byte("null"))
		contents += int64(w)
		return contents, err, true
	}
	return 0, nil, false
}

// boole - boolean to json conversion helper method
func (canon *Inputs) boole(out io.Writer, token bool, err error) (int64, error) {
	var contents int64
	if canon.Spaced {
		var current, err = out.Write([]byte{' '})
		contents += int64(current)
		if err != nil {
			return contents, err
		}
	}
	var current int
	if !token {
		current, err = out.Write([]byte("false"))
	} else {
		current, err = out.Write([]byte("true"))
	}
	contents += int64(current)
	canon.Spaced = true
	return contents, err
}

// f64 - float64 to json conversion helper method
func (canon *Inputs) f64(out io.Writer, token float64) (int64, error) {
	var contents int64
	if canon.Spaced {
		current, err := out.Write([]byte{' '})
		contents += int64(current)
		if err != nil {
			return contents, err
		}
	}
	current, err := NumOut(out, token)
	contents += int64(current)
	canon.Spaced = true
	return contents, err
}

// arr - array to json conversion helper method
func (canon *Inputs) arr(out io.Writer) (int64, error) {
	var contents int64
	current, err := out.Write([]byte{'['})
	contents += int64(current)
	if err != nil {
		return contents, err
	}
	first := true
	for {
		if !canon.Decoder.More() {
			_, err := canon.Decoder.Token()
			if err != nil {
				return contents, err
			}
			w, err := out.Write([]byte{']'})
			contents += int64(w)
			return contents, err
		}
		if !first {
			current, err = out.Write([]byte{','})
			contents += int64(current)
			if err != nil {
				return contents, err
			}
		}
		first = false
		current64, err := canon.Process(out)
		canon.Spaced = false
		contents += current64
		if err != nil {
			return contents, err
		}
	}
}

// obj - object to json conversion helper method
func (canon *Inputs) obj(out io.Writer) (int64, error) {
	var contents tuples
	for {
		if !canon.Decoder.More() {
			_, err := canon.Decoder.Token()
			if err != nil {
				return 0, err
			}
			return canon.ObjectOut(out, contents)
		}
		var key string
		token, err := canon.Decoder.Token()
		if err != nil {
			return 0, err
		}
		switch t := token.(type) {
		case string:
			key = t
		default:
			return 0, fmt.Errorf("unexpected type %T (%v) expected string key", t, t)
		}
		buf := new(bytes.Buffer)
		_, err = canon.Process(buf)
		canon.Spaced = false
		if err != nil {
			return 0, err
		}
		contents = append(contents, tuple{key: key, val: buf.Bytes()})
	}
}
