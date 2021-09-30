package canonjson_test

import (
	"canon-json"
	"canon-json/utils"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"reflect"
	"strings"
	"testing"
	"testing/quick"
)


func read(path string) {
	// open file
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	// remember to close the file at the end of the program
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {

		}
	}(f)

	// read csv values using csv.Reader
	csvReader := csv.NewReader(f)
	for {
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		// got tired here
		// easy to use this to make tests use for loop from data file instead of being hardcoded
		fmt.Printf("%+v\n", rec)
	}
}

/**
Number Formatting
MUST represent all integer numbers (those with a zero-valued fractional part)

without a leading minus sign when the value is zero, and
without a decimal point, and
without an exponent, and
without insignificant leading zeroes (as already required of all JSON numbers)

MUST represent all non-integer numbers in exponential notation

including a nonzero single-digit significand integer part, and
including a nonempty significand fractional part, and
including no trailing zeroes in the significand fractional part (other than as part of a “.0” required to satisfy the preceding point), and
including a capital “E”, and
including no plus sign in the exponent, and
including no insignificant leading zeroes in the exponent
*/
func Test_NumberFormat(t *testing.T) {
	expect := func(f float64, val string) {
		var buf strings.Builder
		_, err := utils.NumOut(&buf, f)
		if err != nil {
			t.Error(err)
			return
		}
		if buf.String() != val {
			t.Errorf("%f became %s, not %s", f, buf.String(), val)
		}
	}

	expect(-150, "-150")
	expect(-150.1, "-1.501E2")
	expect(-15, "-15")
	expect(-15.1, "-1.51E1")
	expect(-1.5, "-1.5E0")
	expect(-1, "-1")
	expect(-0.15, "-1.5E-1")
	expect(-0.0, "0")
	expect(0, "0")
	expect(0.001, "1E-3")
	expect(0.01, "1E-2")
	expect(0.1, "1E-1")
	expect(0.11, "1.1E-1")
	expect(1, "1")
	expect(1.5, "1.5E0")
	expect(10, "10")
	expect(10.5, "1.05E1")
	expect(100, "100")
	expect(100.5, "1.005E2")
	expect(10000000000, "10000000000")
	expect(10000000000.1, "1.00000000001E10")
	expect(1000000000000000, "1000000000000000")
	expect(10000000000000000, "1E16")
	expect(1-(1<<53), "-9007199254740991")
	expect(-(1 << 53), "-9.007199254740992E15")
	expect((1<<53)-1, "9007199254740991")
	expect(1<<53, "9.007199254740992E15")
}

func Test_Canonify(t *testing.T) {
	expect := func(in, out string) {
		var buf strings.Builder
		_, err := utils.Canonify(&buf, strings.NewReader(in))
		if err != nil {
			t.Error(err)
			return
		}
		if buf.String() != out {
			t.Errorf("%v was turned into %v, should have been %v", in, buf.String(), out)
		}
	}
	expect("null", "null")
	expect(`[   null]`, `[null]`)
	expect(`{"b": {"hello": "world"}, "a": 10}`, `{"a":10,"b":{"hello":"world"}}`)
	expect(`[100.5, 1.5, 1.0, 1.0E0, 1.0e+0, 0.1e1]`, `[1.005E2,1.5E0,1,1,1,1]`)
	expect(`{"ab": [1, 2],"aba":1}{"x":"y"}`, `{"ab":[1,2],"aba":1}{"x":"y"}`)
	expect(`null null true false 1e1 1.5`, `null null true false 10 1.5E0`)
	expect(``, ``)
	expect(``, ``)
	expect(`99.9 [1,2] 3 [] null {} false true {} false`, `9.99E1[1,2]3[]null{}false true{}false`)
	expect(`-0.0 -0e0`, `0 0`)
	expect(`"hello\nworld"`, `"hello\nworld"`)
	expect(`"\u002F\u002f\//"`, `"////"`)
	expect(`"\u000F"`, `"\u000f"`)
	expect(`"\u000D"`, `"\r"`)
}

func Test_Spacing(t *testing.T) {
	expect := func(expected string, vs ...interface{}) {
		var buf strings.Builder
		enc := canonjson.CreateEncoder(&buf)
		for _, v := range vs {
			err := enc.Encode(v)
			if err != nil {
				t.Error(err)
				return
			}
		}
		if buf.String() != expected {
			t.Errorf("%v was turned into %v, should have been %v", vs, buf.String(), expected)
		}
	}

	expect("null null", nil, nil)
	expect("[1,2,3][4,5,6]", []int{1, 2, 3}, []int{4, 5, 6})
	expect(`{"a":"b","c":"d","e":"f","g":"h"}`, map[string]string{"c": "d", "g": "h", "e": "f", "a": "b"})
	expect(`0 0`, 0.0, 0.0)
	expect(`1E-1 1.1E0`, 0.1, 1.1)
}

func Test_NonSpaced(t *testing.T) {
	expect := func(expected string, vs ...interface{}) {
		var buf strings.Builder
		enc := canonjson.CreateEncoder(&buf)
		enc.SetSpacing(false)
		for _, v := range vs {
			err := enc.Encode(v)
			if err != nil {
				t.Error(err)
				return
			}
		}
		if buf.String() != expected {
			t.Errorf("%v became %v, not %v", vs, buf.String(), expected)
		}
	}

	expect("nullnull", nil, nil)
	expect("[1,2,3][4,5,6]", []int{1, 2, 3}, []int{4, 5, 6})
	expect(`{"a":"b","c":"d","e":"f","g":"h"}`, map[string]string{"c": "d", "g": "h", "e": "f", "a": "b"})
	expect(`00`, 0.0, 0.0)
	expect(`1E-11.1E0`, 0.1, 1.1)
}

func Test_Marshaling(t *testing.T) {
	expect := func(v interface{}, expected string) {
		bs, err := canonjson.Marshal(v)
		if err != nil {
			t.Error(err)
			return
		}
		if string(bs) != expected {
			t.Errorf("%v was turned into %v, should have been %v", v, string(bs), expected)
		}
	}

	expect(nil, "null")
	expect([]int{1, 2, 3}, "[1,2,3]")
	expect(map[string]string{"c": "d", "g": "h", "e": "f", "a": "b"}, `{"a":"b","c":"d","e":"f","g":"h"}`)
	expect(0.0, `0`)
	expect(0.1, `1E-1`)
	expect(1.1, `1.1E0`)
}

func Test_StringEquality(t *testing.T) {
	f := func(x string) bool {
		bs, err := canonjson.Marshal(x)
		if err != nil {
			return false
		}
		var v string
		err = json.Unmarshal(bs, &v)
		if err != nil {
			return false
		}
		return v == x
	}
	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}


/**
MUST represent all non-integer numbers in exponential notation

    including a nonzero single-digit significand integer part, and
    including a nonempty significand fractional part, and
    including no trailing zeroes in the significand fractional part (other than as part of a “.0” required to satisfy the preceding point), and
    including a capital “E”, and
    including no plus sign in the exponent, and
    including no insignificant leading zeroes in the exponent

 */
type int54 int64

func (i int54) Generate(rand *rand.Rand, size int) reflect.Value {
	return reflect.ValueOf(int54(rand.Int63n((1<<54)-1) - (1 << 53) + 1))
}

func Test_NumEquality(t *testing.T) {
	f := func(x int54) bool {
		bs, err := canonjson.Marshal(float64(x))
		if err != nil {
			return false
		}
		var v int54
		err = json.Unmarshal(bs, &v)
		if err != nil {
			return false
		}
		return v == x
	}
	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func Test_FloatEquality(t *testing.T) {
	f := func(x float64) bool {
		bs, err := canonjson.Marshal(x)
		if err != nil {
			return false
		}
		var v float64
		err = json.Unmarshal(bs, &v)
		if err != nil {
			return false
		}
		return v == x
	}
	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}

}
