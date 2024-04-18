package embedded

import (
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"gitlab.com/pygolo/py"
)

type PyRuntime struct {
	pyCtx py.Py
}

func NewPythonRuntime() (*PyRuntime, error) {
	_py, err := py.GoEmbed()
	if err != nil {
		return nil, err
	}
	return &PyRuntime{
		pyCtx: _py,
	}, nil
}

func (r *PyRuntime) Close() error {
	return r.pyCtx.Close()
}

func (r *PyRuntime) NewCallable(impl CodeSnippet) (*PyCallable, error) {
	defLine, ok := lo.Find(impl, func(item string) bool {
		return strings.Contains(item, "def")
	})
	if !ok {
		return nil, errors.New("function definition not found")
	}
	opParentheseIdx := strings.Index(defLine, "(")
	clParentheseIdx := strings.Index(defLine, ")")
	funcName := defLine[4:opParentheseIdx]
	signature := strings.Split(defLine[opParentheseIdx+1:clParentheseIdx], ",")
	signature = lo.Map(signature, func(s string, _ int) string { return strings.Trim(s, " ") })

	fileName := uuid.New().String() + ".py"
	callableObj, err := r.pyCtx.CompileString(impl.String(), fileName, py.File_input)
	if err != nil {
		return nil, err
	}
	return &PyCallable{
		parentRuntime: &r.pyCtx,
		fileName:      fileName,
		funcName:      funcName,
		signature:     signature,
		pyObj:         callableObj,
	}, nil
}

func (r *PyRuntime) NewDict(m map[any]any) (PyValue, error) {
	if m == nil {
		m = map[any]any{}
	}
	return r.pyObjConstructorImpl(m)
}

func (r *PyRuntime) NewList(l []any) (PyValue, error) {
	if l == nil {
		l = []any{}
	}
	return r.pyObjConstructorImpl(l)
}

func (r *PyRuntime) NewPrimitive(v any) (PyValue, error) {
	return r.pyObjConstructorImpl(v)
}

func (r *PyRuntime) pyObjConstructorImpl(obj any) (result PyValue, err error) {
	var pyObj py.Object
	if obj == nil {
		pyObj = py.None
	} else {
		pyObj, err = r.pyCtx.GoToObject(obj)
		if err != nil {
			return PyValue{}, err
		}
	}
	return PyValue{
		parentRuntime: &r.pyCtx,
		pyObj:         pyObj,
	}, nil
}

type PyCallable struct {
	parentRuntime *py.Py
	fileName      string
	funcName      string
	signature     []string
	pyObj         py.Object
}

func (c PyCallable) Call(args ...PyValue) (PyValue, error) {
	if len(args) != len(c.signature) {
		return PyValue{}, errors.New("arguments amount in function call doesn't match it's signature")
	}
	globalsDict, err := c.parentRuntime.Dict_New()
	if err != nil {
		return PyValue{}, err
	}
	defer c.parentRuntime.DecRef(globalsDict)
	localsDict, err := c.parentRuntime.Dict_New()
	if err != nil {
		return PyValue{}, err
	}
	defer c.parentRuntime.DecRef(localsDict)
	if _, err := c.parentRuntime.Eval_EvalCode(c.pyObj, globalsDict, localsDict); err != nil {
		return PyValue{}, err
	}
	code, err := c.parentRuntime.CompileString(
		fmt.Sprintf("%s(%s)", c.funcName, strings.Join(c.signature, ",")),
		c.fileName,
		py.Eval_input,
	)
	if err != nil {
		return PyValue{}, err
	}
	for i, a := range args {
		if err := c.parentRuntime.Dict_SetItem(localsDict, c.signature[i], a.pyObj); err != nil {
			return PyValue{}, err
		}
	}
	result, err := c.parentRuntime.Eval_EvalCode(code, globalsDict, localsDict)
	if err != nil {
		return PyValue{}, err
	}
	return PyValue{
		pyObj:         result,
		parentRuntime: c.parentRuntime,
	}, nil
}

type PyValue struct {
	parentRuntime *py.Py
	pyObj         py.Object
}

func (v PyValue) IsNone() bool {
	return v.pyObj == py.None
}

func (v PyValue) ToMap() (result map[any]any, err error) {
	if err = v.parentRuntime.GoFromObject(v.pyObj, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (v PyValue) ToSlice() (result []any, err error) {
	if err = v.parentRuntime.GoFromObject(v.pyObj, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (v PyValue) ToAny() (result any, err error) {
	if err = v.parentRuntime.GoFromObject(v.pyObj, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (v PyValue) Decref() {
	v.parentRuntime.DecRef(v.pyObj)
}
