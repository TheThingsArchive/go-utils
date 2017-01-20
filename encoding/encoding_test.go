package encoding

import (
	"math"
	"strconv"
	"testing"

	s "github.com/smartystreets/assertions"
)

const (
	testTag = "test"

	Int     = "int"
	Int8    = "int8"
	Int16   = "int16"
	Int32   = "int32"
	Int64   = "int64"
	Uint    = "uint"
	Uint8   = "uint8"
	Uint16  = "uint16"
	Uint32  = "uint32"
	Uint64  = "uint64"
	Float32 = "float32"
	Float64 = "float64"
	String  = "string"
	Bool    = "bool"
)

type embStruct struct{}
type embInterface interface{}
type testStruct struct {
	embStruct
	embInterface

	Int     int     `test:"int"`
	Int8    int8    `test:"int8"`
	Int16   int16   `test:"int16"`
	Int32   int32   `test:"int32"`
	Int64   int64   `test:"int64"`
	Uint    uint    `test:"uint"`
	Uint8   uint    `test:"uint8"`
	Uint16  uint16  `test:"uint16"`
	Uint32  uint32  `test:"uint32"`
	Uint64  uint64  `test:"uint64"`
	Float32 float32 `test:"float32"`
	Float64 float64 `test:"float64"`
	String  string  `test:"string"`
	Bool    bool    `test:"bool"`

	Struct    struct{}
	Interface interface{}
}

func TestDecode(t *testing.T) {
	for name, arg := range map[string]interface{}{
		"by value":   testStruct{},
		"by pointer": &testStruct{},
	} {
		t.Run(name, func(t *testing.T) {

			a := s.New(t)

			m := map[string]string{
				Int:     strconv.Itoa(42),
				Int8:    strconv.FormatInt(int64(math.MaxInt8), 10),
				Int16:   strconv.FormatInt(int64(math.MaxInt16), 10),
				Int32:   strconv.FormatInt(int64(math.MaxInt32), 10),
				Int64:   strconv.FormatInt(int64(math.MaxInt64), 10),
				Uint:    strconv.FormatUint(42, 10),
				Uint8:   strconv.FormatUint(uint64(math.MaxUint8), 10),
				Uint16:  strconv.FormatUint(uint64(math.MaxUint16), 10),
				Uint32:  strconv.FormatUint(uint64(math.MaxUint32), 10),
				Uint64:  strconv.FormatUint(uint64(math.MaxUint64), 10),
				Float32: strconv.FormatFloat(float64(math.MaxFloat32), 'e', 7, 32),
				Float64: strconv.FormatFloat(float64(math.MaxFloat64), 'e', 16, 64),
				String:  "string",
				Bool:    "bool",
			}

			ret, err := Decode(testTag, arg, m)
			a.So(err, s.ShouldBeNil)

			var v testStruct
			if vPtr, ok := ret.(*testStruct); ok {
				a.So(vPtr, s.ShouldNotBeNil)
				v = *vPtr
			} else {
				v, ok = ret.(testStruct)
				a.So(ok, s.ShouldBeTrue)
			}

			a.So(v.Int, s.ShouldEqual, func() int { val, _ := strconv.ParseInt(m[Int], 10, 0); return int(val) }())
			a.So(v.Int8, s.ShouldEqual, func() int8 { val, _ := strconv.ParseInt(m[Int8], 10, 8); return int8(val) }())
			a.So(v.Int16, s.ShouldEqual, func() int16 { val, _ := strconv.ParseInt(m[Int16], 10, 16); return int16(val) }())
			a.So(v.Int32, s.ShouldEqual, func() int32 { val, _ := strconv.ParseInt(m[Int32], 10, 32); return int32(val) }())
			a.So(v.Int64, s.ShouldEqual, func() int64 { val, _ := strconv.ParseInt(m[Int64], 10, 64); return int64(val) }())
			a.So(v.Uint, s.ShouldEqual, func() uint { val, _ := strconv.ParseUint(m[Uint], 10, 0); return uint(val) }())
			a.So(v.Uint8, s.ShouldEqual, func() uint8 { val, _ := strconv.ParseUint(m[Uint8], 10, 8); return uint8(val) }())
			a.So(v.Uint16, s.ShouldEqual, func() uint16 { val, _ := strconv.ParseUint(m[Uint16], 10, 16); return uint16(val) }())
			a.So(v.Uint32, s.ShouldEqual, func() uint32 { val, _ := strconv.ParseUint(m[Uint32], 10, 32); return uint32(val) }())
			a.So(v.Uint64, s.ShouldEqual, func() uint64 { val, _ := strconv.ParseUint(m[Uint64], 10, 64); return uint64(val) }())
			a.So(v.Float32, s.ShouldEqual, func() float32 { val, _ := strconv.ParseFloat(m[Float32], 32); return float32(val) }())
			a.So(v.Float64, s.ShouldEqual, func() float64 { val, _ := strconv.ParseFloat(m[Float64], 64); return float64(val) }())
			a.So(v.Bool, s.ShouldEqual, func() bool { val, _ := strconv.ParseBool(m[Bool]); return bool(val) }())
			a.So(v.String, s.ShouldEqual, m[String])
		})
	}
}

func TestEncode(t *testing.T) {
	v := testStruct{
		Int:     42,
		Int8:    math.MaxInt8,
		Int16:   math.MaxInt16,
		Int32:   math.MaxInt32,
		Int64:   math.MaxInt64,
		Uint:    42,
		Uint8:   math.MaxUint8,
		Uint16:  math.MaxUint16,
		Uint32:  math.MaxUint32,
		Uint64:  math.MaxUint64,
		Float32: math.MaxFloat32,
		Float64: math.MaxFloat64,
	}

	for name, arg := range map[string]interface{}{
		"by value":   v,
		"by pointer": &v,
	} {
		t.Run(name, func(t *testing.T) {
			a := s.New(t)
			enc, err := Encode(testTag, arg)
			a.So(err, s.ShouldBeNil)

			a.So(enc[Int], s.ShouldEqual, strconv.FormatInt(int64(v.Int), 10))
			a.So(enc[Int8], s.ShouldEqual, strconv.FormatInt(int64(v.Int8), 10))
			a.So(enc[Int16], s.ShouldEqual, strconv.FormatInt(int64(v.Int16), 10))
			a.So(enc[Int32], s.ShouldEqual, strconv.FormatInt(int64(v.Int32), 10))
			a.So(enc[Int64], s.ShouldEqual, strconv.FormatInt(int64(v.Int64), 10))
			a.So(enc[Uint], s.ShouldEqual, strconv.FormatUint(uint64(v.Uint), 10))
			a.So(enc[Uint8], s.ShouldEqual, strconv.FormatUint(uint64(v.Uint8), 10))
			a.So(enc[Uint16], s.ShouldEqual, strconv.FormatUint(uint64(v.Uint16), 10))
			a.So(enc[Uint32], s.ShouldEqual, strconv.FormatUint(uint64(v.Uint32), 10))
			a.So(enc[Uint64], s.ShouldEqual, strconv.FormatUint(uint64(v.Uint64), 10))
			a.So(enc[Float32], s.ShouldEqual, strconv.FormatFloat(float64(v.Float32), 'e', 7, 32))
			a.So(enc[Float64], s.ShouldEqual, strconv.FormatFloat(float64(v.Float64), 'e', 16, 64))
			a.So(enc[Bool], s.ShouldEqual, strconv.FormatBool(v.Bool))
			a.So(enc[String], s.ShouldEqual, v.String)
		})
	}
}
