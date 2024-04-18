/*
 * Copyright 2022, Pygolo Project contributors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package py

// #include "pygolo.h"
//
// int pyLong_CheckExact(PyObject *o)
// {
//     return PyLong_CheckExact(o);
// }
import "C"
import (
	"fmt"
	"math"
	"reflect"
	"regexp"
)

// Long_Type wraps the Python PyLong_Type type object.
//
// C API: https://docs.python.org/3/c-api/long.html#c.PyLong_Type
var Long_Type = TypeObject{&C.PyLong_Type}

// Long_CheckExact returns true if o is of type long, subtypes excluded.
//
// C API: https://docs.python.org/3/c-api/long.html#c.PyLong_CheckExact
func (Py Py) Long_CheckExact(o Object) bool {
	return C.pyLong_CheckExact(o.o) != 0
}

// Long_FromLong returns a new long object from int v.
//
// C API: https://docs.python.org/3/c-api/long.html#c.PyLong_FromLong
func (Py Py) Long_FromLong(v int) (Object, error) {
	return Py.wrap(C.PyLong_FromLong(C.long(v)))
}

// Long_FromLongLong returns a new long object from int64 v.
//
// C API: https://docs.python.org/3/c-api/long.html#c.PyLong_FromLongLong
func (Py Py) Long_FromLongLong(v int64) (Object, error) {
	return Py.wrap(C.PyLong_FromLongLong(C.longlong(v)))
}

// Long_FromUnsignedLong returns a new long object from uint v.
//
// C API: https://docs.python.org/3/c-api/long.html#c.PyLong_FromUnsignedLong
func (Py Py) Long_FromUnsignedLong(v uint) (Object, error) {
	return Py.wrap(C.PyLong_FromUnsignedLong(C.ulong(v)))
}

// Long_FromUnsignedLongLong returns a new long object from uint64 v.
//
// C API: https://docs.python.org/3/c-api/long.html#c.PyLong_FromUnsignedLongLong
func (Py Py) Long_FromUnsignedLongLong(v uint64) (Object, error) {
	return Py.wrap(C.PyLong_FromUnsignedLongLong(C.ulonglong(v)))
}

// Long_AsLong returns an int representation of o.
//
// C API: https://docs.python.org/3/c-api/long.html#c.PyLong_AsLong
func (Py Py) Long_AsLong(o Object) (int, error) {
	v := C.PyLong_AsLong(o.o)
	if v == -1 && Py.Err_Occurred() != (Object{}) {
		return 0, Py.GoCatchError()
	}
	return int(v), nil
}

// Long_AsLongLong returns an int64 representation of o.
//
// C API: https://docs.python.org/3/c-api/long.html#c.PyLong_AsLongLong
func (Py Py) Long_AsLongLong(o Object) (int64, error) {
	v := C.PyLong_AsLongLong(o.o)
	if v == -1 && Py.Err_Occurred() != (Object{}) {
		return 0, Py.GoCatchError()
	}
	return int64(v), nil
}

// Long_AsUnsignedLong returns a uint representation of o.
//
// C API: https://docs.python.org/3/c-api/long.html#c.PyLong_AsUnsignedLong
func (Py Py) Long_AsUnsignedLong(o Object) (uint, error) {
	v := C.PyLong_AsUnsignedLong(o.o)
	if C.long(v) == -1 && Py.Err_Occurred() != (Object{}) {
		return 0, Py.GoCatchError()
	}
	return uint(v), nil
}

// Long_AsUnsignedLongLong returns a uint64 representation of o.
//
// C API: https://docs.python.org/3/c-api/long.html#c.PyLong_AsUnsignedLongLong
func (Py Py) Long_AsUnsignedLongLong(o Object) (uint64, error) {
	v := C.PyLong_AsUnsignedLongLong(o.o)
	if C.longlong(v) == -1 && Py.Err_Occurred() != (Object{}) {
		return 0, Py.GoCatchError()
	}
	return uint64(v), nil
}

// longToObject converts a Go integer value to a Python long.
func longToObject(Py Py, a interface{}) (o Object, e error) {
	switch v := a.(type) {
	case int:
		o, e = Py.Long_FromLong(v)
	case int8:
		o, e = Py.Long_FromLong(int(v))
	case int16:
		o, e = Py.Long_FromLong(int(v))
	case int32:
		o, e = Py.Long_FromLong(int(v))
	case int64:
		o, e = Py.Long_FromLongLong(v)
	case uint:
		o, e = Py.Long_FromUnsignedLong(v)
	case uint8:
		o, e = Py.Long_FromUnsignedLong(uint(v))
	case uint16:
		o, e = Py.Long_FromUnsignedLong(uint(v))
	case uint32:
		o, e = Py.Long_FromUnsignedLong(uint(v))
	case uint64:
		o, e = Py.Long_FromUnsignedLongLong(v)
	default:
		e = fmt.Errorf("not an integer: %v", a)
	}
	return
}

// longFromObject converts a Python long to a Go integer value.
func longFromObject(Py Py, o Object, a interface{}) (e error) {
	if !Py.Long_CheckExact(o) {
		return Py.GoErrorConvFromObject(o, a)
	}
	switch target := a.(type) {
	case *int:
		var v int
		if v, e = Py.Long_AsLong(o); e == nil {
			*target = v
			return
		}
	case *int8:
		var v int
		if v, e = Py.Long_AsLong(o); e == nil && math.MinInt8 <= v && v <= math.MaxInt8 {
			*target = int8(v)
			return
		}
	case *int16:
		var v int
		if v, e = Py.Long_AsLong(o); e == nil && math.MinInt16 <= v && v <= math.MaxInt16 {
			*target = int16(v)
			return
		}
	case *int32:
		var v int
		if v, e = Py.Long_AsLong(o); e == nil && math.MinInt32 <= v && v <= math.MaxInt32 {
			*target = int32(v)
			return
		}
	case *int64:
		var v int64
		if v, e = Py.Long_AsLongLong(o); e == nil {
			*target = v
			return
		}
	case *uint:
		var v uint
		if v, e = Py.Long_AsUnsignedLong(o); e == nil {
			*target = v
			return
		}
	case *uint8:
		var v uint
		if v, e = Py.Long_AsUnsignedLong(o); e == nil && v <= math.MaxUint8 {
			*target = uint8(v)
			return
		}
	case *uint16:
		var v uint
		if v, e = Py.Long_AsUnsignedLong(o); e == nil && v <= math.MaxUint16 {
			*target = uint16(v)
			return
		}
	case *uint32:
		var v uint
		if v, e = Py.Long_AsUnsignedLong(o); e == nil && v <= math.MaxUint32 {
			*target = uint32(v)
			return
		}
	case *uint64:
		var v uint64
		if v, e = Py.Long_AsUnsignedLongLong(o); e == nil {
			*target = v
			return
		}
	case *interface{}:
		var vi int
		if vi, e = Py.Long_AsLong(o); e == nil {
			*target = vi
			return
		}
		var vu uint
		if vu, e = Py.Long_AsUnsignedLong(o); e == nil {
			*target = vu
			return
		}
	default:
		e = Py.GoErrorConvFromObject(o, a)
	}
	if e == nil {
		//lint:ignore ST1005 'Python' is a proper name
		e = fmt.Errorf("Python int too large to convert to Go %s", reflect.TypeOf(a).Elem())
	} else if matched, _ := regexp.MatchString("OverflowError: can't convert negative .*", e.Error()); matched {
		//lint:ignore ST1005 'Python' is a proper name
		e = fmt.Errorf("Python int is negative, cannot convert to Go %s", reflect.TypeOf(a).Elem())
	} else if matched, _ := regexp.MatchString("OverflowError: .*", e.Error()); matched {
		//lint:ignore ST1005 'Python' is a proper name
		e = fmt.Errorf("Python int too large to convert to Go %s", reflect.TypeOf(a).Elem())
	}
	return
}

func init() {
	for _, kind := range []reflect.Kind{
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
	} {
		c := GoConvConf{
			Kind:       kind,
			ToObject:   longToObject,
			FromObject: longFromObject,
		}
		if err := c.Register(); err != nil {
			panic(err)
		}
	}
	c := GoConvConf{
		TypeObject: Long_Type,
		FromObject: longFromObject,
	}
	if err := c.Register(); err != nil {
		panic(err)
	}
}
