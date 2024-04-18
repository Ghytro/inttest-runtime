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
import "C"
import (
	"fmt"
	"reflect"
)

// GoConvToObject is the type of a Go→Python conversion handler.
type GoConvToObject func(Py, interface{}) (Object, error)

// GoConvFromObject is the type of a Python→Go conversion handler.
type GoConvFromObject func(Py, Object, interface{}) error

var fromTypeConverters = make(map[reflect.Type]GoConvToObject)
var fromKindConverters = make(map[reflect.Kind]GoConvToObject)
var toTypeConverters = make(map[reflect.Type]GoConvFromObject)
var toKindConverters = make(map[reflect.Kind]GoConvFromObject)
var toTypeObjectConverters = make(map[TypeObject]GoConvFromObject)

// GoConvConf configures the Go⟷Python conversion handling.
//
// When converting from Go to Python, the (source) Go value is examined
// to find the most suitable conversion function. The type is examined
// first, if no handler for that type is found then the kind is examined.
// If again no handler is found, the conversion fails.
//
// When converting from Python to Go, in addition to the rules above
// also the type of the (source) Python object is examined as last.
type GoConvConf struct {
	// TypeOf is any Go value which type is key in the conversion.
	//
	// The type of this value, not the value, is used as key for registering
	// the conversion handlers. If omitted, the handlers are not registered
	// in the conversion map based on value type.
	TypeOf interface{}

	// Kind identifies a Go value kind key in the conversion.
	//
	// If omitted, the handlers are not registered in the conversion map based
	// on value kind.
	reflect.Kind

	// TypeObject identifies the Python type key in the conversion.
	//
	// If omitted, the handler is not registered in the Python type conversion map.
	TypeObject

	// ToObject handles the Go→Python conversion
	ToObject GoConvToObject

	// FromObject handles the Python→Go conversion
	FromObject GoConvFromObject
}

// Register adds the FromObject and ToObject handlers to the appropriate
// conversion maps.
func (c GoConvConf) Register() error {
	if c.TypeOf == nil && c.Kind == reflect.Invalid && c.TypeObject == (TypeObject{}) {
		return fmt.Errorf("either TypeOf, Kind or TypeObject must be set")
	}
	if c.TypeObject != (TypeObject{}) && c.FromObject == nil {
		return fmt.Errorf("FromObject must be set")
	}
	if c.ToObject == nil && c.FromObject == nil {
		return fmt.Errorf("either ToObject or FromObject must be set")
	}
	Type := reflect.TypeOf(c.TypeOf)
	if Type != nil {
		_, found_from := fromTypeConverters[Type]
		_, found_to := toTypeConverters[Type]
		if found_from && c.ToObject != nil || found_to && c.FromObject != nil {
			return fmt.Errorf("Type handler is already registered: %s", Type)
		}
	}
	if c.Kind != reflect.Invalid {
		_, found_from := fromKindConverters[c.Kind]
		_, found_to := toKindConverters[c.Kind]
		if found_from && c.ToObject != nil || found_to && c.FromObject != nil {
			//lint:ignore ST1005 'Kind' is a proper name
			return fmt.Errorf("Kind handler is already registered: %s", c.Kind)
		}
	}
	if c.TypeObject != (TypeObject{}) {
		_, found_to := toTypeObjectConverters[c.TypeObject]
		if found_to && c.FromObject != nil {
			return fmt.Errorf("TypeObject handler is already registered: %v", c.TypeObject)
		}
	}
	if Type != nil && c.ToObject != nil {
		fromTypeConverters[Type] = c.ToObject
	}
	if Type != nil && c.FromObject != nil {
		toTypeConverters[Type] = c.FromObject
	}
	if c.Kind != reflect.Invalid && c.ToObject != nil {
		fromKindConverters[c.Kind] = c.ToObject
	}
	if c.Kind != reflect.Invalid && c.FromObject != nil {
		toKindConverters[c.Kind] = c.FromObject
	}
	if c.TypeObject != (TypeObject{}) && c.FromObject != nil {
		toTypeObjectConverters[c.TypeObject] = c.FromObject
	}
	return nil
}

// Unregister removes the FromObject and ToObject handlers from all the
// conversion maps.
func (c GoConvConf) Unregister() {
	Type := reflect.TypeOf(c.TypeOf)
	if c.TypeOf != nil {
		delete(fromTypeConverters, Type)
		delete(toTypeConverters, Type)
	}
	if c.Kind != reflect.Invalid {
		delete(fromKindConverters, c.Kind)
		delete(toKindConverters, c.Kind)
	}
	if c.TypeObject != (TypeObject{}) {
		delete(toTypeObjectConverters, c.TypeObject)
	}
}

// GoToObject converts a Go value to a Python object.
//
// A new object is created or a new reference to an existing one is returned.
func (Py Py) GoToObject(a interface{}) (Object, error) {
	if o, ok := a.(Object); ok {
		return Py.NewRef(o), nil
	}
	if a == nil {
		return Object{}, Py.GoErrorConvToObject(a, TypeObject{})
	}
	if from := fromTypeConverters[reflect.ValueOf(a).Type()]; from != nil {
		return from(Py, a)
	}
	if from := fromKindConverters[reflect.ValueOf(a).Kind()]; from != nil {
		return from(Py, a)
	}
	return Object{}, Py.GoErrorConvToObject(a, TypeObject{})
}

// GoFromObject converts a Python object to a Go value.
//
// A pointer to the Go value is passed so that its type can be examined
// for guiding the conversion. Either the type is suitable or the conversion
// fails.
//
// If `any` is used as Go type, then the conversion is driven by the
// Python object type and it's left to the caller to use type assertion
// for accessing the actual value.
func (Py Py) GoFromObject(o Object, a interface{}) error {
	if reflect.ValueOf(a).Kind() != reflect.Ptr {
		return fmt.Errorf("cannot store in %T, need a pointer", a)
	}
	if o == (Object{}) {
		return fmt.Errorf("cannot convert Python <nil> to Go %s", reflect.TypeOf(a).Elem())
	}
	if a, ok := a.(*Object); ok {
		*a = Py.NewRef(o)
		return nil
	}
	ptr := reflect.TypeOf(a).Elem()
	if to := toTypeConverters[ptr]; to != nil {
		return to(Py, o, a)
	}
	if k := ptr.Kind(); k != reflect.Interface {
		if to := toKindConverters[k]; to != nil {
			return to(Py, o, a)
		}
	}
	if to := toTypeObjectConverters[Py.Type(o)]; to != nil {
		return to(Py, o, a)
	}
	return Py.GoErrorConvFromObject(o, a)
}

// GoErrorConvToObject formats an error encountered by GoToObject.
func (Py Py) GoErrorConvToObject(a interface{}, t TypeObject) error {
	if t.t == nil {
		return fmt.Errorf("cannot convert Go %T to Python", a)
	}
	return fmt.Errorf("cannot convert Go %T to Python %s", a, t.Name())
}

// GoErrorConvFromObject formats an error encountered by GoFromObject.
func (Py Py) GoErrorConvFromObject(o Object, a interface{}) error {
	return fmt.Errorf("cannot convert Python %s to Go %s",
		Py.Type(o).Name(), reflect.TypeOf(a).Elem())
}
