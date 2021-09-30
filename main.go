package canonjson

import (
	"bytes"
	"canon-json/utils"
	"encoding/json"
	"io"
)

// Marshal - Returns the canon json encoding of an input
func Marshal(input interface{}) ([]byte, error) {
	var buf bytes.Buffer
	err := CreateEncoder(&buf).Encode(input)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Encoder - Writes canon JSON values to an output stream.
type Encoder struct {
	canon   utils.Inputs
	out     io.Writer
	spacing bool
}

// CreateEncoder - Creates an encoder for a writer
func CreateEncoder(w io.Writer) *Encoder {
	return &Encoder{
		out:     w,
		spacing: true,
	}
}

// Encode - Write marshaled Json of input to the stream.
// - Separates multiple inputs with spaces (Enabled by default)
// - Uses encoding/json.Marshal to convert Go Values to Json.
func (enc *Encoder) Encode(input interface{}) error {
	marshaled, err := json.Marshal(input)
	if err != nil {
		return err
	}
	enc.canon.Decoder = json.NewDecoder(bytes.NewBuffer(marshaled))
	_, err = enc.canon.Process(enc.out)
	enc.canon.Spaced = enc.canon.Spaced && enc.spacing
	return err
}

// SetSpacing - Enable or Disable spacing between multiple inputs
func (enc *Encoder) SetSpacing(on bool) {
	enc.spacing = on
}
