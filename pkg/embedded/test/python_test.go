package test

import (
	"fmt"
	"inttest-runtime/pkg/embedded"
	"testing"
)

func TestPythonFunction(t *testing.T) {
	pyRuntime, err := embedded.NewPythonRuntime()
	if err != nil {
		t.Fatal(err)
	}
	defer pyRuntime.Close()

	myFunc := embedded.CodeSnippet{
		"def mockApiFunction(url_params, query_params, headers, body):",
		"\tif url_params[\"id\"] == \"1\":",
		"\t\treturn {\"error\": \"element with id does not exist\"}",
		"\tif \"X-Access-Token\" not in headers.keys():",
		"\t\treturn {\"error\": \"Access token not specified\"}",
		"\treturn {\"my_response\": \"my response is correct\", \"nested_objects\": {\"are_valid\": \"too\"}}",
	}

	callable, err := pyRuntime.NewCallable(myFunc)
	if err != nil {
		t.Fatal(err)
	}

	urlParams := map[any]any{
		"id": 1,
	}
	queryParams := map[any]any{
		"something": "here",
	}
	headers := map[any]any{
		"X-Access-Token": "is here",
	}
	body := map[any]any{
		"valid_json_object": true,
	}
	urlParamsObj, err := pyRuntime.NewDict(urlParams)
	if err != nil {
		t.Fatal(err)
	}
	queryParamsObj, err := pyRuntime.NewDict(queryParams)
	if err != nil {
		t.Fatal(err)
	}
	headersObj, err := pyRuntime.NewDict(headers)
	if err != nil {
		t.Fatal(err)
	}
	bodyObj, err := pyRuntime.NewDict(body)
	if err != nil {
		t.Fatal(err)
	}
	result, err := callable.Call(urlParamsObj, queryParamsObj, headersObj, bodyObj)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(result.ToMap())
}

func BenchmarkCallableWithRuntimeInit(b *testing.B) {
	for i := 0; i < b.N; i++ {
		pyRuntime, err := embedded.NewPythonRuntime()
		if err != nil {
			b.Fatal(err)
		}
		defer pyRuntime.Close()

		myFunc := embedded.CodeSnippet{
			"def mockApiFunction(url_params, query_params, headers, body):",
			"\tif url_params[\"id\"] == \"1\":",
			"\t\treturn {\"error\": \"element with id does not exist\"}",
			"\tif \"X-Access-Token\" not in headers.keys():",
			"\t\treturn {\"error\": \"Access token not specified\"}",
			"\treturn {\"my_response\": \"my response is correct\", \"nested_objects\": {\"are_valid\": \"too\"}}",
		}

		callable, err := pyRuntime.NewCallable(myFunc)
		if err != nil {
			b.Fatal(err)
		}

		urlParams := map[any]any{
			"id": 1,
		}
		queryParams := map[any]any{
			"something": "here",
		}
		headers := map[any]any{
			"X-Access-Token": "is here",
		}
		body := map[any]any{
			"valid_json_object": true,
		}
		urlParamsObj, err := pyRuntime.NewDict(urlParams)
		if err != nil {
			b.Fatal(err)
		}
		queryParamsObj, err := pyRuntime.NewDict(queryParams)
		if err != nil {
			b.Fatal(err)
		}
		headersObj, err := pyRuntime.NewDict(headers)
		if err != nil {
			b.Fatal(err)
		}
		bodyObj, err := pyRuntime.NewDict(body)
		if err != nil {
			b.Fatal(err)
		}
		result, err := callable.Call(urlParamsObj, queryParamsObj, headersObj, bodyObj)
		if err != nil {
			b.Fatal(err)
		}
		fmt.Println(result.ToMap())
		pyRuntime.Close()
	}
}

func BenchmarkCallableNoRuntimeInit(b *testing.B) {
	pyRuntime, err := embedded.NewPythonRuntime()
	if err != nil {
		b.Fatal(err)
	}
	defer pyRuntime.Close()

	myFunc := embedded.CodeSnippet{
		"def mockApiFunction(url_params, query_params, headers, body):",
		"\tif url_params[\"id\"] == \"1\":",
		"\t\treturn {\"error\": \"element with id does not exist\"}",
		"\tif \"X-Access-Token\" not in headers.keys():",
		"\t\treturn {\"error\": \"Access token not specified\"}",
		"\treturn {\"my_response\": \"my response is correct\", \"nested_objects\": {\"are_valid\": \"too\"}}",
	}

	callable, err := pyRuntime.NewCallable(myFunc)
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		urlParams := map[any]any{
			"id": 1,
		}
		queryParams := map[any]any{
			"something": "here",
		}
		headers := map[any]any{
			"X-Access-Token": "is here",
		}
		body := map[any]any{
			"valid_json_object": true,
		}
		urlParamsObj, err := pyRuntime.NewDict(urlParams)
		if err != nil {
			b.Fatal(err)
		}
		queryParamsObj, err := pyRuntime.NewDict(queryParams)
		if err != nil {
			b.Fatal(err)
		}
		headersObj, err := pyRuntime.NewDict(headers)
		if err != nil {
			b.Fatal(err)
		}
		bodyObj, err := pyRuntime.NewDict(body)
		if err != nil {
			b.Fatal(err)
		}
		result, err := callable.Call(urlParamsObj, queryParamsObj, headersObj, bodyObj)
		if err != nil {
			b.Fatal(err)
		}
		fmt.Println(result.ToMap())
	}
}
