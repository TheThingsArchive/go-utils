// Copyright © 2017 The Things Network
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

package encoding

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"testing"
	"time"

	s "github.com/smartystreets/assertions"
)

func marshalToString(v interface{}) string {
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

	timeVar time.Time = time.Date(2017, 01, 02, 03, 04, 05, 06, time.UTC)

	intSliceVar    = []int{4, 2}
	boolSliceVar   = []bool{true, false}
	stringSliceVar = []string{"te", "st"}

	intArrayVar    = [2]int{4, 2}
	boolArrayVar   = [2]bool{true, false}
	stringArrayVar = [2]string{"te", "st"}

	stringStringMapVar = map[string]string{"te": "st"}

	inclStructVar = inclStruct{Int: intVar, SubInclStruct: subInclStruct{Int: intVar}}

	SquashedFieldVar  = "squashed field"
	squashedStructVar = SquashedStruct{
		SquashedField: SquashedFieldVar,
	}
)

var (
	timeVarString      = "2017-01-02T03:04:05.000000006Z"
	localTimeVarString = "2017-01-02T05:04:05.000000006+02:00"

	intSliceVarString    = marshalToString(intSliceVar)
	boolSliceVarString   = marshalToString(boolSliceVar)
	stringSliceVarString = marshalToString(stringSliceVar)

	intArrayVarString    = marshalToString(intArrayVar)
	boolArrayVarString   = marshalToString(boolArrayVar)
	stringArrayVarString = marshalToString(stringArrayVar)

	stringStringMapVarString = marshalToString(stringStringMapVar)
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
	Time    = "time"

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
	TimePtr    = "timePtr"

	IntSlice    = "intSlice"
	BoolSlice   = "boolSlice"
	StringSlice = "stringSlice"

	IntArray    = "intArray"
	BoolArray   = "boolArray"
	StringArray = "stringArray"

	StringStringMap = "stringStringMap"

	InclStruct         = "inclStruct"
	InclStructEmpty    = "inclStructEmpty"
	InclStructPtr      = "inclStructPtr"
	InclStructPtrEmpty = "inclStructPtrEmpty"

	InclStructEmb = "inclStructEmb"

	SubInclStruct = "subInclStruct"

	EmbStructTag    = "embStruct"
	EmbInterfaceTag = "embInterface"

	SquashedField = "squashedField"
)

type inclStruct struct {
	Int           int           `test:"int"`
	SubInclStruct subInclStruct `test:"subInclStruct,include"`
}

type subInclStruct struct {
	Int int `test:"int"`
}

type EmbStruct struct {
	Int int `test:"int"`
}
type EmbInterface interface {
	Test()
}

type SquashedStruct struct {
	SquashedField string `test:"squashedField"`
}
type testStruct struct {
	Int     int       `test:"int"`
	Int8    int8      `test:"int8"`
	Int16   int16     `test:"int16"`
	Int32   int32     `test:"int32"`
	Int64   int64     `test:"int64"`
	Uint    uint      `test:"uint"`
	Uint8   uint8     `test:"uint8"`
	Uint16  uint16    `test:"uint16"`
	Uint32  uint32    `test:"uint32"`
	Uint64  uint64    `test:"uint64"`
	Float32 float32   `test:"float32"`
	Float64 float64   `test:"float64"`
	String  string    `test:"string"`
	Bool    bool      `test:"bool"`
	Time    time.Time `test:"time"`

	IntPtr     *int       `test:"intPtr"`
	Int8Ptr    *int8      `test:"int8Ptr"`
	Int16Ptr   *int16     `test:"int16Ptr"`
	Int32Ptr   *int32     `test:"int32Ptr"`
	Int64Ptr   *int64     `test:"int64Ptr"`
	UintPtr    *uint      `test:"uintPtr"`
	Uint8Ptr   *uint8     `test:"uint8Ptr"`
	Uint16Ptr  *uint16    `test:"uint16Ptr"`
	Uint32Ptr  *uint32    `test:"uint32Ptr"`
	Uint64Ptr  *uint64    `test:"uint64Ptr"`
	Float32Ptr *float32   `test:"float32Ptr"`
	Float64Ptr *float64   `test:"float64Ptr"`
	StringPtr  *string    `test:"stringPtr"`
	BoolPtr    *bool      `test:"boolPtr"`
	TimePtr    *time.Time `test:"timePtr"`

	IntSlice    []int    `test:"intSlice"`
	BoolSlice   []bool   `test:"boolSlice"`
	StringSlice []string `test:"stringSlice"`

	IntArray    [2]int    `test:"intArray"`
	BoolArray   [2]bool   `test:"boolArray"`
	StringArray [2]string `test:"stringArray"`

	StringStringMap map[string]string `test:"stringStringMap"`

	EmbStruct    `test:"embStruct,include"`
	EmbInterface `test:"embInterface,include"`

	SquashedStruct `test:",squash"`

	InclStruct         inclStruct  `test:"inclStruct,include"`
	InclStructEmpty    inclStruct  `test:"inclStructEmpty,include,omitempty"`
	InclStructPtr      *inclStruct `test:"inclStructPtr,include"`
	InclStructPtrEmpty *inclStruct `test:"inclStructPtrEmpty,include,omitempty"`

	inclStruct `test:"inclStructEmb,include"`
}

func TestFromStringStringMap(t *testing.T) {
	for name, arg := range map[string]interface{}{
		"by value":   testStruct{},
		"by pointer": &testStruct{},
	} {
		t.Run(name, func(t *testing.T) {
			a := s.New(t)

			if argPtr, ok := arg.(*testStruct); ok {
				old := *argPtr
				defer a.So(*argPtr, s.ShouldResemble, old)
			}

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
				Time:    localTimeVarString,

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
				TimePtr:    timeVarString,

				IntSlice:    intSliceVarString,
				BoolSlice:   boolSliceVarString,
				StringSlice: stringSliceVarString,

				IntArray:    intArrayVarString,
				BoolArray:   boolArrayVarString,
				StringArray: stringArrayVarString,

				StringStringMap: stringStringMapVarString,

				EmbStructTag + "." + Int: strconv.Itoa(intVar),

				InclStruct + "." + Int:                       strconv.Itoa(intVar),
				InclStruct + "." + SubInclStruct + "." + Int: strconv.Itoa(intVar),

				InclStructPtr + "." + Int:                       strconv.Itoa(intVar),
				InclStructPtr + "." + SubInclStruct + "." + Int: strconv.Itoa(intVar),

				InclStructEmb + "." + Int:                       strconv.Itoa(intVar),
				InclStructEmb + "." + SubInclStruct + "." + Int: strconv.Itoa(intVar),

				SquashedField: SquashedFieldVar,
			}

			ret, err := FromStringStringMap(testTag, arg, m)
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
			a.So(v.Time, s.ShouldResemble, timeVar)

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
			if a.So(v.TimePtr, s.ShouldNotBeNil) {
				a.So(*v.TimePtr, s.ShouldResemble, timeVar)
			}

			a.So(v.IntSlice, s.ShouldResemble, intSliceVar)
			a.So(v.BoolSlice, s.ShouldResemble, boolSliceVar)
			a.So(v.StringSlice, s.ShouldResemble, stringSliceVar)

			a.So(v.IntArray, s.ShouldResemble, intArrayVar)
			a.So(v.BoolArray, s.ShouldResemble, boolArrayVar)
			a.So(v.StringArray, s.ShouldResemble, stringArrayVar)

			a.So(v.StringStringMap, s.ShouldResemble, stringStringMapVar)

			a.So(v.EmbStruct.Int, s.ShouldEqual, func() int { val, _ := strconv.ParseInt(m[EmbStructTag+"."+Int], 10, 0); return int(val) }())

			a.So(v.SquashedStruct.SquashedField, s.ShouldEqual, SquashedFieldVar)

			a.So(v.InclStruct.Int, s.ShouldEqual, func() int { val, _ := strconv.ParseInt(m[InclStruct+"."+Int], 10, 0); return int(val) }())
			a.So(v.InclStruct.SubInclStruct.Int, s.ShouldEqual, func() int {
				val, _ := strconv.ParseInt(m[InclStruct+"."+SubInclStruct+"."+Int], 10, 0)
				return int(val)
			}())

			if a.So(v.InclStructPtr, s.ShouldNotBeNil) {
				a.So(v.InclStructPtr.Int, s.ShouldEqual, func() int { val, _ := strconv.ParseInt(m[InclStructPtr+"."+Int], 10, 0); return int(val) }())
				a.So(v.InclStructPtr.SubInclStruct.Int, s.ShouldEqual, func() int {
					val, _ := strconv.ParseInt(m[InclStructPtr+"."+SubInclStruct+"."+Int], 10, 0)
					return int(val)
				}())
			}

			a.So(v.inclStruct, s.ShouldBeZeroValue)
			a.So(v.InclStructEmpty, s.ShouldBeZeroValue)
			a.So(v.InclStructPtrEmpty, s.ShouldBeNil)
		})
	}
}

var testStructVar = testStruct{
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
	Time:    timeVar,

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
	TimePtr:    &timeVar,

	IntSlice:    intSliceVar,
	BoolSlice:   boolSliceVar,
	StringSlice: stringSliceVar,

	IntArray:    intArrayVar,
	BoolArray:   boolArrayVar,
	StringArray: stringArrayVar,

	StringStringMap: stringStringMapVar,

	InclStruct:    inclStructVar,
	InclStructPtr: &inclStructVar,

	inclStruct: inclStructVar,

	SquashedStruct: squashedStructVar,
}

func TestToStringStringMap(t *testing.T) {
	v := testStructVar

	for name, arg := range map[string]interface{}{
		"by value":   v,
		"by pointer": &v,
	} {
		t.Run(name, func(t *testing.T) {
			a := s.New(t)

			if argPtr, ok := arg.(*testStruct); ok {
				old := *argPtr
				defer a.So(*argPtr, s.ShouldResemble, old)
			}

			enc, err := ToStringStringMap(testTag, arg)
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
			a.So(enc[Time], s.ShouldEqual, timeVarString)

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
			a.So(enc[TimePtr], s.ShouldEqual, timeVarString)

			a.So(enc[IntSlice], s.ShouldEqual, marshalToString(v.IntSlice))
			a.So(enc[BoolSlice], s.ShouldEqual, marshalToString(v.BoolSlice))
			a.So(enc[StringSlice], s.ShouldEqual, marshalToString(v.StringSlice))

			a.So(enc[IntArray], s.ShouldEqual, marshalToString(v.IntArray))
			a.So(enc[BoolArray], s.ShouldEqual, marshalToString(v.BoolArray))
			a.So(enc[StringArray], s.ShouldEqual, marshalToString(v.StringArray))

			a.So(enc[StringStringMap], s.ShouldEqual, marshalToString(v.StringStringMap))

			a.So(enc[EmbStructTag+".int"], s.ShouldEqual, strconv.FormatInt(int64(v.EmbStruct.Int), 10))

			a.So(enc[SquashedField], s.ShouldEqual, SquashedFieldVar)

			a.So(enc[InclStruct+".int"], s.ShouldEqual, strconv.FormatInt(int64(v.InclStruct.Int), 10))
			a.So(enc[InclStruct+"."+SubInclStruct+".int"], s.ShouldEqual, strconv.FormatInt(int64(v.InclStruct.SubInclStruct.Int), 10))

			a.So(enc[InclStructPtr+".int"], s.ShouldEqual, strconv.FormatInt(int64(v.InclStructPtr.Int), 10))
			a.So(enc[InclStructPtr+"."+SubInclStruct+".int"], s.ShouldEqual, strconv.FormatInt(int64(v.InclStructPtr.SubInclStruct.Int), 10))

			a.So(enc, s.ShouldNotContainKey, InclStructEmb+".int")
			a.So(enc, s.ShouldNotContainKey, InclStructEmb+"."+SubInclStruct+".int")
			a.So(enc, s.ShouldNotContainKey, InclStructEmpty+".int")
			a.So(enc, s.ShouldNotContainKey, InclStructEmpty+"."+SubInclStruct+".int")
			a.So(enc, s.ShouldNotContainKey, InclStructPtrEmpty+".int")
			a.So(enc, s.ShouldNotContainKey, InclStructPtrEmpty+"."+SubInclStruct+".int")

			if arg, ok := arg.(testStruct); ok {
				arg.IntPtr = nil
				arg.Int8Ptr = nil
				arg.Int16Ptr = nil
				arg.Int32Ptr = nil
				arg.Int64Ptr = nil
				arg.UintPtr = nil
				arg.Uint8Ptr = nil
				arg.Uint16Ptr = nil
				arg.Uint32Ptr = nil
				arg.Uint64Ptr = nil
				arg.Float32Ptr = nil
				arg.Float64Ptr = nil
				arg.BoolPtr = nil
				arg.StringPtr = nil
				arg.TimePtr = nil

				enc, err := ToStringStringMap(testTag, arg)
				a.So(err, s.ShouldBeNil)

				a.So(enc[IntPtr], s.ShouldEqual, "")
				a.So(enc[Int8Ptr], s.ShouldEqual, "")
				a.So(enc[Int16Ptr], s.ShouldEqual, "")
				a.So(enc[Int32Ptr], s.ShouldEqual, "")
				a.So(enc[Int64Ptr], s.ShouldEqual, "")
				a.So(enc[UintPtr], s.ShouldEqual, "")
				a.So(enc[Uint8Ptr], s.ShouldEqual, "")
				a.So(enc[Uint16Ptr], s.ShouldEqual, "")
				a.So(enc[Uint32Ptr], s.ShouldEqual, "")
				a.So(enc[Uint64Ptr], s.ShouldEqual, "")
				a.So(enc[Float32Ptr], s.ShouldEqual, "")
				a.So(enc[Float64Ptr], s.ShouldEqual, "")
				a.So(enc[BoolPtr], s.ShouldEqual, "")
				a.So(enc[StringPtr], s.ShouldEqual, "")
				a.So(enc[TimePtr], s.ShouldEqual, "")

				arg.Time = time.Time{}
				enc, err = ToStringStringMap(testTag, arg)
				a.So(err, s.ShouldBeNil)
				a.So(enc[Time], s.ShouldEqual, "")

				arg.Time = time.Unix(0, 0)
				enc, err = ToStringStringMap(testTag, arg)
				a.So(err, s.ShouldBeNil)
				a.So(enc[Time], s.ShouldEqual, "")
			}
		})
	}
}

func TestToStringInterfaceMap(t *testing.T) {
	v := testStructVar

	for name, arg := range map[string]interface{}{
		"by value":   v,
		"by pointer": &v,
	} {
		t.Run(name, func(t *testing.T) {
			a := s.New(t)

			if argPtr, ok := arg.(*testStruct); ok {
				old := *argPtr
				defer a.So(*argPtr, s.ShouldResemble, old)
			}

			enc, err := ToStringInterfaceMap(testTag, arg)
			a.So(err, s.ShouldBeNil)

			a.So(enc[Int], s.ShouldEqual, v.Int)
			a.So(enc[Int8], s.ShouldEqual, v.Int8)
			a.So(enc[Int16], s.ShouldEqual, v.Int16)
			a.So(enc[Int32], s.ShouldEqual, v.Int32)
			a.So(enc[Int64], s.ShouldEqual, v.Int64)
			a.So(enc[Uint], s.ShouldEqual, v.Uint)
			a.So(enc[Uint8], s.ShouldEqual, v.Uint8)
			a.So(enc[Uint16], s.ShouldEqual, v.Uint16)
			a.So(enc[Uint32], s.ShouldEqual, v.Uint32)
			a.So(enc[Uint64], s.ShouldEqual, v.Uint64)
			a.So(enc[Float32], s.ShouldEqual, v.Float32)
			a.So(enc[Float64], s.ShouldEqual, v.Float64)
			a.So(enc[Bool], s.ShouldEqual, v.Bool)
			a.So(enc[String], s.ShouldEqual, v.String)
			a.So(enc[Time], s.ShouldResemble, v.Time)

			a.So(enc[IntPtr], s.ShouldEqual, *v.IntPtr)
			a.So(enc[Int8Ptr], s.ShouldEqual, *v.Int8Ptr)
			a.So(enc[Int16Ptr], s.ShouldEqual, *v.Int16Ptr)
			a.So(enc[Int32Ptr], s.ShouldEqual, *v.Int32Ptr)
			a.So(enc[Int64Ptr], s.ShouldEqual, *v.Int64Ptr)
			a.So(enc[UintPtr], s.ShouldEqual, *v.UintPtr)
			a.So(enc[Uint8Ptr], s.ShouldEqual, *v.Uint8Ptr)
			a.So(enc[Uint16Ptr], s.ShouldEqual, *v.Uint16Ptr)
			a.So(enc[Uint32Ptr], s.ShouldEqual, *v.Uint32Ptr)
			a.So(enc[Uint64Ptr], s.ShouldEqual, *v.Uint64Ptr)
			a.So(enc[Float32Ptr], s.ShouldEqual, *v.Float32Ptr)
			a.So(enc[Float64Ptr], s.ShouldEqual, *v.Float64Ptr)
			a.So(enc[BoolPtr], s.ShouldEqual, *v.BoolPtr)
			a.So(enc[StringPtr], s.ShouldEqual, *v.StringPtr)
			a.So(enc[TimePtr], s.ShouldResemble, *v.TimePtr)

			a.So(enc[IntSlice], s.ShouldResemble, v.IntSlice)
			a.So(enc[BoolSlice], s.ShouldResemble, v.BoolSlice)
			a.So(enc[StringSlice], s.ShouldResemble, v.StringSlice)

			a.So(enc[IntArray], s.ShouldResemble, v.IntArray)
			a.So(enc[BoolArray], s.ShouldResemble, v.BoolArray)
			a.So(enc[StringArray], s.ShouldResemble, v.StringArray)

			a.So(enc[StringStringMap], s.ShouldResemble, v.StringStringMap)

			a.So(enc[EmbStructTag+".int"], s.ShouldEqual, v.EmbStruct.Int)

			a.So(enc[SquashedField], s.ShouldEqual, SquashedFieldVar)

			a.So(enc[InclStruct+".int"], s.ShouldEqual, v.InclStruct.Int)
			a.So(enc[InclStruct+"."+SubInclStruct+".int"], s.ShouldEqual, v.InclStruct.SubInclStruct.Int)

			a.So(enc[InclStructPtr+".int"], s.ShouldEqual, v.InclStructPtr.Int)
			a.So(enc[InclStructPtr+"."+SubInclStruct+".int"], s.ShouldEqual, v.InclStructPtr.SubInclStruct.Int)

			a.So(enc, s.ShouldNotContainKey, InclStructEmb+".int")
			a.So(enc, s.ShouldNotContainKey, InclStructEmb+"."+SubInclStruct+".int")
			a.So(enc, s.ShouldNotContainKey, InclStructEmpty+".int")
			a.So(enc, s.ShouldNotContainKey, InclStructEmpty+"."+SubInclStruct+".int")
			a.So(enc, s.ShouldNotContainKey, InclStructPtrEmpty+".int")
			a.So(enc, s.ShouldNotContainKey, InclStructPtrEmpty+"."+SubInclStruct+".int")

			if arg, ok := arg.(testStruct); ok {
				arg.IntPtr = nil
				arg.Int8Ptr = nil
				arg.Int16Ptr = nil
				arg.Int32Ptr = nil
				arg.Int64Ptr = nil
				arg.UintPtr = nil
				arg.Uint8Ptr = nil
				arg.Uint16Ptr = nil
				arg.Uint32Ptr = nil
				arg.Uint64Ptr = nil
				arg.Float32Ptr = nil
				arg.Float64Ptr = nil
				arg.BoolPtr = nil
				arg.StringPtr = nil

				enc, err := ToStringInterfaceMap(testTag, arg)
				a.So(err, s.ShouldBeNil)

				a.So(enc[IntPtr], s.ShouldEqual, nil)
				a.So(enc[Int8Ptr], s.ShouldEqual, nil)
				a.So(enc[Int16Ptr], s.ShouldEqual, nil)
				a.So(enc[Int32Ptr], s.ShouldEqual, nil)
				a.So(enc[Int64Ptr], s.ShouldEqual, nil)
				a.So(enc[UintPtr], s.ShouldEqual, nil)
				a.So(enc[Uint8Ptr], s.ShouldEqual, nil)
				a.So(enc[Uint16Ptr], s.ShouldEqual, nil)
				a.So(enc[Uint32Ptr], s.ShouldEqual, nil)
				a.So(enc[Uint64Ptr], s.ShouldEqual, nil)
				a.So(enc[Float32Ptr], s.ShouldEqual, nil)
				a.So(enc[Float64Ptr], s.ShouldEqual, nil)
				a.So(enc[BoolPtr], s.ShouldEqual, nil)
				a.So(enc[StringPtr], s.ShouldEqual, nil)
			}

		})
	}
}

func TestNonUniqueKeys(t *testing.T) {
	a := s.New(t)

	invalidStruct1 := struct {
		Conflicting string         `test:"squashedField"`
		Squashed    SquashedStruct `test:",squash"`
	}{
		Conflicting: "hello",
		Squashed: SquashedStruct{
			SquashedField: "there",
		},
	}

	invalidStruct2 := struct {
		Squashed    SquashedStruct `test:",squash"`
		Conflicting string         `test:"squashedField"`
	}{
		Conflicting: "hello",
		Squashed: SquashedStruct{
			SquashedField: "there",
		},
	}

	for _, sub := range []interface{}{invalidStruct1, invalidStruct2} {
		var err error
		catch := func() {
			if r := recover(); r != nil {
				var ok bool
				err, ok = r.(error)
				if !ok {
					err = fmt.Errorf("pkg: %v", r)
				}
			}
		}

		expected := errors.New("field names not unique (squashedField)")

		func() {
			defer catch()
			ToStringStringMap("test", sub)
		}()
		a.So(err, s.ShouldResemble, expected)

		func() {
			defer catch()
			ToStringInterfaceMap("test", sub)
		}()
		a.So(err, s.ShouldResemble, expected)
	}
}

func TestCastType(t *testing.T) {
	a := s.New(t)

	strct := struct {
		Int     uint        `test:"Int,cast=int"`
		Int8    uint8       `test:"Int8,cast=int8"`
		Int16   uint16      `test:"Int16,cast=int16"`
		Int32   uint32      `test:"Int32,cast=int32"`
		Int64   uint64      `test:"Int64,cast=int64"`
		Uint    int         `test:"Uint,cast=uint"`
		Uint8   int8        `test:"Uint8,cast=uint8"`
		Uint16  int16       `test:"Uint16,cast=uint16"`
		Uint32  int32       `test:"Uint32,cast=uint32"`
		Uint64  int64       `test:"Uint64,cast=uint64"`
		Float32 float64     `test:"Float32,cast=float32"`
		Float64 float32     `test:"Float64,cast=float64"`
		String  interface{} `test:"String,cast=string"`
	}{
		Int:     42,
		Int8:    42,
		Int16:   42,
		Int32:   42,
		Int64:   42,
		Uint:    42,
		Uint8:   42,
		Uint16:  42,
		Uint32:  42,
		Uint64:  42,
		Float32: 42,
		Float64: 42,
		String: struct {
			a int64
			b interface{}
			c bool
			d func()
		}{42, struct{ z int }{42}, true, func() {}},
	}

	enc, err := ToStringInterfaceMap("test", strct)
	if a.So(err, s.ShouldBeNil) && a.So(enc, s.ShouldHaveLength, reflect.ValueOf(strct).NumField()) {
		for name, v := range map[string]interface{}{
			"Int":     int(strct.Int),
			"Int8":    int8(strct.Int8),
			"Int16":   int16(strct.Int16),
			"Int32":   int32(strct.Int32),
			"Int64":   int64(strct.Int64),
			"Uint":    uint(strct.Uint),
			"Uint8":   uint8(strct.Uint8),
			"Uint16":  uint16(strct.Uint16),
			"Uint32":  uint32(strct.Uint32),
			"Uint64":  uint64(strct.Uint64),
			"Float32": float32(strct.Float32),
			"Float64": float64(strct.Float64),
			"String":  fmt.Sprint(strct.String),
		} {
			a.So(enc, s.ShouldContainKey, name)
			a.So(enc[name], s.ShouldHaveSameTypeAs, v)
			a.So(enc[name], s.ShouldEqual, v)
		}
	}

	invalid := struct {
		A int `test:"a,cast=unknown"`
	}{42}
	a.So(func() { ToStringInterfaceMap("test", invalid) }, s.ShouldPanicWith, "Wrong cast type specified: unknown")
}
