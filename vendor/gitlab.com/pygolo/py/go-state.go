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
// #include "go-type.h"
import "C"

type state struct {
	GoFunction_Type TypeObject
}

func (Py Py) initTypes(m Object) error {
	t_fn, err := Py.newGoFunction_Type()
	defer Py.DecRef(t_fn.Object())
	if err != nil {
		return err
	}
	if m.o != nil {
		err = Py.moduleAddType(m, t_fn)
		if err != nil {
			return err
		}
	}
	state := Py.state.Get().(*state)
	state.GoFunction_Type = t_fn
	Py.IncRef(t_fn.Object())
	return nil
}

func (Py Py) deinitTypes() {
	state := Py.state.Get().(*state)
	for _, t := range []*TypeObject{
		&state.GoFunction_Type,
	} {
		GoConvConf{TypeObject: *t}.Unregister()
		Py.DecRef(t.Object())
		*t = TypeObject{}
	}
}

func (Py Py) getStateFromModule(m Object) *GoHandle {
	return (*GoHandle)(C.PyModule_GetState(m.o))
}

func (Py Py) getStateFromObject(o Object) *GoHandle {
	return (*GoHandle)(C.getTypeHandle(Py.Type(o).t))
}
