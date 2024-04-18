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
// PyTypeObject* py_Type(PyObject *o)
// {
//     return o ? Py_TYPE(o) : NULL;
// }
//
// PyObject* py_None()
// {
//     return Py_None;
// }
//
// Py_ssize_t py_RefCnt(PyObject *o)
// {
//     return Py_REFCNT(o);
// }
import "C"
import (
	"fmt"
	"unsafe"
)

// None_Type wraps the Python _PyNone_Type type object.
var None_Type = TypeObject{&C._PyNone_Type}

// None wraps the Python Py_None object.
//
// C API: https://docs.python.org/3/c-api/none.html#c.Py_None
var None = Object{C.py_None()}

// Object wraps a pointer to a C API PyObject.
//
// C API: https://docs.python.org/3/c-api/structures.html#c.PyObject
type Object struct {
	o *C.PyObject
}

// TypeObject wraps a pointer to a C API PyTypeObject.
//
// C API: https://docs.python.org/3/c-api/type.html#c.PyTypeObject
type TypeObject struct {
	t *C.PyTypeObject
}

// Name returns the name of the type object.
func (t TypeObject) Name() string {
	return C.GoString(t.t.tp_name)
}

// Object returns the type object as a plain object.
func (t TypeObject) Object() Object {
	return Object{(*C.PyObject)(unsafe.Pointer(t.t))}
}

// GoArgs is the positional arguments slice type in object calling.
type GoArgs []interface{}

// GoKwArgs is the named arguments map type in object calling.
type GoKwArgs map[string]interface{}

// Type returns the type of an object.
//
// If o is Object{} then TypeObject{} is returned.
//
// C API: https://docs.python.org/3/c-api/structures.html#c.Py_TYPE
func (Py Py) Type(o Object) TypeObject {
	return TypeObject{C.py_Type(o.o)}
}

// NewRef increments the reference count of object o and returns it.
//
// C API: https://docs.python.org/3/c-api/refcounting.html#c.Py_XNewRef
func (Py Py) NewRef(o Object) Object {
	C.Py_IncRef(o.o)
	return o
}

// IncRef increments the reference count of objects oo.
//
// C API: https://docs.python.org/3/c-api/refcounting.html#c.Py_IncRef
func (Py Py) IncRef(oo ...Object) {
	for _, o := range oo {
		C.Py_IncRef(o.o)
	}
}

// DecRef decrements the reference count of objects oo.
//
// Invoking DecRef on a zero-Object does not actually do anything and it's
// a supported behavior.
//
// C API: https://docs.python.org/3/c-api/refcounting.html#c.Py_DecRef
func (Py Py) DecRef(oo ...Object) {
	for _, o := range oo {
		C.Py_DecRef(o.o)
	}
}

// RefCnt returns the reference count of object o.
//
// C API: https://docs.python.org/3/c-api/structures.html#c.Py_REFCNT
func (Py Py) RefCnt(o Object) uint {
	return uint(C.py_RefCnt(o.o))
}

// Object_Call calls a callable object o with positional arguments given
// by args and named arguments given by kwargs.
//
// If no positional arguments are needed, args can be nil.
// If no named arguments are needed, kwargs can be nil.
//
// C API: https://docs.python.org/3/c-api/call.html#c.PyObject_Call
func (Py Py) Object_Call(o Object, args interface{}, kwargs interface{}) (Object, error) {
	var o_args, o_kwargs Object
	var err error

	if args == nil {
		args = GoArgs{}
	} else if o_args, ok := args.(Object); ok && o_args == (Object{}) {
		args = GoArgs{}
	}
	switch args := args.(type) {
	case GoArgs:
		o_args, err = Py.Tuple_Pack(args...)
		defer Py.DecRef(o_args)
		if err != nil {
			return Object{}, err
		}
	case Object:
		o_args = args
	default:
		return Object{}, fmt.Errorf("args must be either of type GoArgs or Object")
	}
	if kwargs == nil {
		return Py.wrap(C.PyObject_Call(o.o, o_args.o, nil))
	} else if o_kwargs, ok := kwargs.(Object); ok && o_kwargs == (Object{}) {
		return Py.wrap(C.PyObject_Call(o.o, o_args.o, nil))
	}
	switch kwargs := kwargs.(type) {
	case GoKwArgs:
		o_kwargs, err = Py.GoToObject(kwargs)
		defer Py.DecRef(o_kwargs)
		if err != nil {
			return Object{}, err
		}
	case Object:
		o_kwargs = kwargs
	default:
		return Object{}, fmt.Errorf("kwargs must be either of type GoKwArgs or Object")
	}
	return Py.wrap(C.PyObject_Call(o.o, o_args.o, o_kwargs.o))
}

// Object_CallFunction calls a callable object o with a variable number
// of arguments.
//
// C API: https://docs.python.org/3/c-api/call.html#c.PyObject_CallFunction
func (Py Py) Object_CallFunction(o Object, args ...interface{}) (Object, error) {
	o_args, err := Py.Tuple_Pack(args...)
	defer Py.DecRef(o_args)
	if err != nil {
		return Object{}, err
	}
	return Py.wrap(C.PyObject_CallObject(o.o, o_args.o))
}

// Object_CallMethod calls the method named name of object o with a
// variable number of arguments.
//
// C API: https://docs.python.org/3/c-api/call.html#c.PyObject_CallMethod
func (Py Py) Object_CallMethod(o Object, name string, args ...interface{}) (Object, error) {
	o_method, err := Py.Object_GetAttr(o, name)
	defer Py.DecRef(o_method)
	if err != nil {
		return Object{}, err
	}
	return Py.Object_CallFunction(o_method, args...)
}

// Object_Length returns the length of object o.
//
// C API: https://docs.python.org/3/c-api/object.html#c.PyObject_Length
func (Py Py) Object_Length(o Object) (int, error) {
	ret := C.PyObject_Length(o.o)
	if ret == -1 {
		return 0, Py.GoCatchError()
	}
	return int(ret), nil
}

// Object_Str computes a string representation of object o.
//
// C API: https://docs.python.org/3/c-api/object.html#c.PyObject_Str
func (Py Py) Object_Str(o Object) (Object, error) {
	return Py.wrap(C.PyObject_Str(o.o))
}

// Object_IsTrue returns true if the object is considered to be true,
// false otherwise.
//
// C API: https://docs.python.org/3/c-api/object.html#c.PyObject_IsTrue
func (Py Py) Object_IsTrue(o Object) (bool, error) {
	ret := C.PyObject_IsTrue(o.o)
	if ret == -1 {
		return false, Py.GoCatchError()
	}
	return ret == 1, nil
}

// Object_GetAttr retrieves an attribute named attr_name from object o.
//
// C API: https://docs.python.org/3/c-api/object.html#c.PyObject_GetAttr
func (Py Py) Object_GetAttr(o Object, attr_name interface{}) (Object, error) {
	o_attr_name, err := Py.GoToObject(attr_name)
	defer Py.DecRef(o_attr_name)
	if err != nil {
		return Object{}, err
	}
	return Py.wrap(C.PyObject_GetAttr(o.o, o_attr_name.o))
}

// Object_SetAttr sets the attribute named attr_name to value attr_value.
//
// C API: https://docs.python.org/3/c-api/object.html#c.PyObject_SetAttr
func (Py Py) Object_SetAttr(o Object, attr_name interface{}, attr_value interface{}) error {
	o_attr_name, err := Py.GoToObject(attr_name)
	defer Py.DecRef(o_attr_name)
	if err != nil {
		return err
	}
	o_attr_value, err := Py.GoToObject(attr_value)
	defer Py.DecRef(o_attr_value)
	if err != nil {
		return err
	}
	if Py.pglCheck(o_attr_value) {
		err = Py.Object_SetAttr(o_attr_value, "__name__", o_attr_name)
		if err != nil {
			return err
		}
	}
	if C.PyObject_SetAttr(o.o, o_attr_name.o, o_attr_value.o) == -1 {
		return Py.GoCatchError()
	}
	return nil
}
