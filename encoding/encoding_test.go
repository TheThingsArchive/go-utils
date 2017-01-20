package encoding

import (
	"encoding/json"
	"math"
	"strconv"
	"testing"

	s "github.com/smartystreets/assertions"
)

func sliceToString(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}

var (
	intVar     int     = 42
	int8Var    int8    = math.MaxInt8
	int16Var   int16   = math.MaxInt16
	int32Var   int32   = math.MaxInt32
	int64Var   int64   = math.MaxInt64
	uintVar    uint    = 42
	uint8Var   uint8   = math.MaxUint8
	uint16Var  uint16  = math.MaxUint16
	uint32Var  uint32  = math.MaxUint32
	uint64Var  uint64  = math.MaxUint64
	float32Var float32 = math.MaxFloat32
	float64Var float64 = math.MaxFloat64
	boolVar    bool    = true
	stringVar  string  = "test"

	intSliceVar    = []int{4, 2}
	boolSliceVar   = []bool{true, false}
	stringSliceVar = []string{"te", "st"}

	intArrayVar    = [2]int{4, 2}
	boolArrayVar   = [2]bool{true, false}
	stringArrayVar = [2]string{"te", "st"}
)

var (
	intSliceVarString    = sliceToString(intSliceVar)
	boolSliceVarString   = sliceToString(boolSliceVar)
	stringSliceVarString = sliceToString(stringSliceVar)

	intArrayVarString    = sliceToString(intArrayVar)
	boolArrayVarString   = sliceToString(boolArrayVar)
	stringArrayVarString = sliceToString(stringArrayVar)
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

	IntPtr     = "intPtr"
	Int8Ptr    = "int8Ptr"
	Int16Ptr   = "int16Ptr"
	Int32Ptr   = "int32Ptr"
	Int64Ptr   = "int64Ptr"
	UintPtr    = "uintPtr"
	Uint8Ptr   = "uint8Ptr"
	Uint16Ptr  = "uint16Ptr"
	Uint32Ptr  = "uint32Ptr"
	Uint64Ptr  = "uint64Ptr"
	Float32Ptr = "float32Ptr"
	Float64Ptr = "float64Ptr"
	StringPtr  = "stringPtr"
	BoolPtr    = "boolPtr"

	IntSlice    = "intSlice"
	BoolSlice   = "boolSlice"
	StringSlice = "stringSlice"

	IntArray    = "intArray"
	BoolArray   = "boolArray"
	StringArray = "stringArray"

	Struct    = "struct"
	Interface = "interface"
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
	Uint8   uint8   `test:"uint8"`
	Uint16  uint16  `test:"uint16"`
	Uint32  uint32  `test:"uint32"`
	Uint64  uint64  `test:"uint64"`
	Float32 float32 `test:"float32"`
	Float64 float64 `test:"float64"`
	String  string  `test:"string"`
	Bool    bool    `test:"bool"`

	IntPtr     *int     `test:"intPtr"`
	Int8Ptr    *int8    `test:"int8Ptr"`
	Int16Ptr   *int16   `test:"int16Ptr"`
	Int32Ptr   *int32   `test:"int32Ptr"`
	Int64Ptr   *int64   `test:"int64Ptr"`
	UintPtr    *uint    `test:"uintPtr"`
	Uint8Ptr   *uint8   `test:"uint8Ptr"`
	Uint16Ptr  *uint16  `test:"uint16Ptr"`
	Uint32Ptr  *uint32  `test:"uint32Ptr"`
	Uint64Ptr  *uint64  `test:"uint64Ptr"`
	Float32Ptr *float32 `test:"float32Ptr"`
	Float64Ptr *float64 `test:"float64Ptr"`
	StringPtr  *string  `test:"stringPtr"`
	BoolPtr    *bool    `test:"boolPtr"`

	Struct    struct{}    `test:"struct"`
	Interface interface{} `test:"interface"`

	IntSlice    []int    `test:"intSlice"`
	BoolSlice   []bool   `test:"boolSlice"`
	StringSlice []string `test:"stringSlice"`

	IntArray    [2]int    `test:"intArray"`
	BoolArray   [2]bool   `test:"boolArray"`
	StringArray [2]string `test:"stringArray"`
}

func TestDecode(t *testing.T) {
	for name, arg := range map[string]interface{}{
		"by value":   testStruct{},
		"by pointer": &testStruct{},
	} {
		t.Run(name, func(t *testing.T) {

			a := s.New(t)

			m := map[string]string{
				Int:     strconv.Itoa(intVar),
				Int8:    strconv.FormatInt(int64(int8Var), 10),
				Int16:   strconv.FormatInt(int64(int16Var), 10),
				Int32:   strconv.FormatInt(int64(int32Var), 10),
				Int64:   strconv.FormatInt(int64(int64Var), 10),
				Uint:    strconv.FormatUint(uint64(uintVar), 10),
				Uint8:   strconv.FormatUint(uint64(uint8Var), 10),
				Uint16:  strconv.FormatUint(uint64(uint16Var), 10),
				Uint32:  strconv.FormatUint(uint64(uint32Var), 10),
				Uint64:  strconv.FormatUint(uint64(uint64Var), 10),
				Float32: strconv.FormatFloat(float64(float32Var), 'e', 7, 32),
				Float64: strconv.FormatFloat(float64(float64Var), 'e', 16, 64),
				String:  stringVar,
				Bool:    strconv.FormatBool(boolVar),

				IntPtr:     strconv.Itoa(intVar),
				Int8Ptr:    strconv.FormatInt(int64(int8Var), 10),
				Int16Ptr:   strconv.FormatInt(int64(int16Var), 10),
				Int32Ptr:   strconv.FormatInt(int64(int32Var), 10),
				Int64Ptr:   strconv.FormatInt(int64(int64Var), 10),
				UintPtr:    strconv.FormatUint(uint64(uintVar), 10),
				Uint8Ptr:   strconv.FormatUint(uint64(uint8Var), 10),
				Uint16Ptr:  strconv.FormatUint(uint64(uint16Var), 10),
				Uint32Ptr:  strconv.FormatUint(uint64(uint32Var), 10),
				Uint64Ptr:  strconv.FormatUint(uint64(uint64Var), 10),
				Float32Ptr: strconv.FormatFloat(float64(float32Var), 'e', 7, 32),
				Float64Ptr: strconv.FormatFloat(float64(float64Var), 'e', 16, 64),
				StringPtr:  stringVar,
				BoolPtr:    strconv.FormatBool(boolVar),

				IntSlice:    intSliceVarString,
				BoolSlice:   boolSliceVarString,
				StringSlice: stringSliceVarString,

				IntArray:    intArrayVarString,
				BoolArray:   boolArrayVarString,
				StringArray: stringArrayVarString,
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

			if a.So(v.IntPtr, s.ShouldNotBeNil) {
				a.So(*v.IntPtr, s.ShouldEqual, func() int { val, _ := strconv.ParseInt(m[Int], 10, 0); return int(val) }())
			}
			if a.So(v.Int8Ptr, s.ShouldNotBeNil) {
				a.So(*v.Int8Ptr, s.ShouldEqual, func() int8 { val, _ := strconv.ParseInt(m[Int8], 10, 8); return int8(val) }())
			}
			if a.So(v.Int16Ptr, s.ShouldNotBeNil) {
				a.So(*v.Int16Ptr, s.ShouldEqual, func() int16 { val, _ := strconv.ParseInt(m[Int16], 10, 16); return int16(val) }())
			}
			if a.So(v.Int32Ptr, s.ShouldNotBeNil) {
				a.So(*v.Int32Ptr, s.ShouldEqual, func() int32 { val, _ := strconv.ParseInt(m[Int32], 10, 32); return int32(val) }())
			}
			if a.So(v.Int64Ptr, s.ShouldNotBeNil) {
				a.So(*v.Int64Ptr, s.ShouldEqual, func() int64 { val, _ := strconv.ParseInt(m[Int64], 10, 64); return int64(val) }())
			}
			if a.So(v.UintPtr, s.ShouldNotBeNil) {
				a.So(*v.UintPtr, s.ShouldEqual, func() uint { val, _ := strconv.ParseUint(m[Uint], 10, 0); return uint(val) }())
			}
			if a.So(v.Uint8Ptr, s.ShouldNotBeNil) {
				a.So(*v.Uint8Ptr, s.ShouldEqual, func() uint8 { val, _ := strconv.ParseUint(m[Uint8], 10, 8); return uint8(val) }())
			}
			if a.So(v.Uint16Ptr, s.ShouldNotBeNil) {
				a.So(*v.Uint16Ptr, s.ShouldEqual, func() uint16 { val, _ := strconv.ParseUint(m[Uint16], 10, 16); return uint16(val) }())
			}
			if a.So(v.Uint32Ptr, s.ShouldNotBeNil) {
				a.So(*v.Uint32Ptr, s.ShouldEqual, func() uint32 { val, _ := strconv.ParseUint(m[Uint32], 10, 32); return uint32(val) }())
			}
			if a.So(v.Uint64Ptr, s.ShouldNotBeNil) {
				a.So(*v.Uint64Ptr, s.ShouldEqual, func() uint64 { val, _ := strconv.ParseUint(m[Uint64], 10, 64); return uint64(val) }())
			}
			if a.So(v.Float32Ptr, s.ShouldNotBeNil) {
				a.So(*v.Float32Ptr, s.ShouldEqual, func() float32 { val, _ := strconv.ParseFloat(m[Float32], 32); return float32(val) }())
			}
			if a.So(v.Float64Ptr, s.ShouldNotBeNil) {
				a.So(*v.Float64Ptr, s.ShouldEqual, func() float64 { val, _ := strconv.ParseFloat(m[Float64], 64); return float64(val) }())
			}
			if a.So(v.BoolPtr, s.ShouldNotBeNil) {
				a.So(*v.BoolPtr, s.ShouldEqual, func() bool { val, _ := strconv.ParseBool(m[Bool]); return bool(val) }())
			}
			if a.So(v.StringPtr, s.ShouldNotBeNil) {
				a.So(*v.StringPtr, s.ShouldEqual, m[String])
			}

			a.So(v.IntSlice, s.ShouldResemble, intSliceVar)
			a.So(v.BoolSlice, s.ShouldResemble, boolSliceVar)
			a.So(v.StringSlice, s.ShouldResemble, stringSliceVar)

			a.So(v.IntArray, s.ShouldResemble, intArrayVar)
			a.So(v.BoolArray, s.ShouldResemble, boolArrayVar)
			a.So(v.StringArray, s.ShouldResemble, stringArrayVar)
		})
	}
}

func TestEncode(t *testing.T) {
	v := testStruct{
		Int:     intVar,
		Int8:    int8Var,
		Int16:   int16Var,
		Int32:   int32Var,
		Int64:   int64Var,
		Uint:    uintVar,
		Uint8:   uint8Var,
		Uint16:  uint16Var,
		Uint32:  uint32Var,
		Uint64:  uint64Var,
		Float32: float32Var,
		Float64: float64Var,
		Bool:    boolVar,
		String:  stringVar,

		IntPtr:     &intVar,
		Int8Ptr:    &int8Var,
		Int16Ptr:   &int16Var,
		Int32Ptr:   &int32Var,
		Int64Ptr:   &int64Var,
		UintPtr:    &uintVar,
		Uint8Ptr:   &uint8Var,
		Uint16Ptr:  &uint16Var,
		Uint32Ptr:  &uint32Var,
		Uint64Ptr:  &uint64Var,
		Float32Ptr: &float32Var,
		Float64Ptr: &float64Var,
		BoolPtr:    &boolVar,
		StringPtr:  &stringVar,

		IntSlice:    intSliceVar,
		BoolSlice:   boolSliceVar,
		StringSlice: stringSliceVar,

		IntArray:    intArrayVar,
		BoolArray:   boolArrayVar,
		StringArray: stringArrayVar,
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

			a.So(enc[IntPtr], s.ShouldEqual, strconv.FormatInt(int64(v.Int), 10))
			a.So(enc[Int8Ptr], s.ShouldEqual, strconv.FormatInt(int64(v.Int8), 10))
			a.So(enc[Int16Ptr], s.ShouldEqual, strconv.FormatInt(int64(v.Int16), 10))
			a.So(enc[Int32Ptr], s.ShouldEqual, strconv.FormatInt(int64(v.Int32), 10))
			a.So(enc[Int64Ptr], s.ShouldEqual, strconv.FormatInt(int64(v.Int64), 10))
			a.So(enc[UintPtr], s.ShouldEqual, strconv.FormatUint(uint64(v.Uint), 10))
			a.So(enc[Uint8Ptr], s.ShouldEqual, strconv.FormatUint(uint64(v.Uint8), 10))
			a.So(enc[Uint16Ptr], s.ShouldEqual, strconv.FormatUint(uint64(v.Uint16), 10))
			a.So(enc[Uint32Ptr], s.ShouldEqual, strconv.FormatUint(uint64(v.Uint32), 10))
			a.So(enc[Uint64Ptr], s.ShouldEqual, strconv.FormatUint(uint64(v.Uint64), 10))
			a.So(enc[Float32Ptr], s.ShouldEqual, strconv.FormatFloat(float64(v.Float32), 'e', 7, 32))
			a.So(enc[Float64Ptr], s.ShouldEqual, strconv.FormatFloat(float64(v.Float64), 'e', 16, 64))
			a.So(enc[BoolPtr], s.ShouldEqual, strconv.FormatBool(v.Bool))
			a.So(enc[StringPtr], s.ShouldEqual, v.String)

			a.So(enc[IntSlice], s.ShouldEqual, sliceToString(v.IntSlice))
			a.So(enc[BoolSlice], s.ShouldEqual, sliceToString(v.BoolSlice))
			a.So(enc[StringSlice], s.ShouldEqual, sliceToString(v.StringSlice))

			a.So(enc[IntArray], s.ShouldEqual, sliceToString(v.IntArray))
			a.So(enc[BoolArray], s.ShouldEqual, sliceToString(v.BoolArray))
			a.So(enc[StringArray], s.ShouldEqual, sliceToString(v.StringArray))
		})
	}
}
