/**
@author: Krisna Pranav, GodzillaFrameworkDevelopers
@filename: js/js.go

Copyright [2021 - 2023] [Krisna Pranav, GodzillaFrameworkDeveloeprs]

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v8

import (
	"os"

	"github.com/livebud/bud/package/js"
	"go.kuoruan.net/v8go-polyfills/console"
	"go.kuoruan.net/v8go-polyfills/fetch"
	"go.kuoruan.net/v8go-polyfills/timers"
	"go.kuoruan.net/v8go-polyfills/url"
	"rogchap.com/v8go"
)

type Value = v8go.Value
type Error = v8go.JSError

var _ js.VM = (*VM)(nil)

type VM struct {
	isolate *v8go.Isolate
	context *v8go.Context
}

func (vm *VM) Eval(path, expr string) (string, error) {
	value, err := vm.context.RunScript(expr, path)
	if err != nil {
		return "", err
	}

	if value.IsPromise() {
		prom, err := value.AsPromise()
		if err != nil {
			return "", err
		}

		for prom.State() == v8go.Pending {
			continue
		}
		return prom.Result().String(), nil
	}
	return value.String(), nil
}

func Eval(path, code string) (string, error) {
	vm, err := Load()
	if err != nil {
		return "", err
	}
	return vm.Eval(path, code)
}

func load() (*v8go.Isolate, *v8go.Context, error) {
	isolate := v8go.NewIsolate()
	global := v8go.NewObjectTemplate(isolate)
	if err := fetch.InjectTo(isolate, global); err != nil {
		isolate.TerminateExecution()
		isolate.Dispose()
		return nil, nil, err
	}

	if err := timers.InjectTo(isolate, global); err != nil {
		isolate.TerminateExecution()
		isolate.Dispose()
		return nil, nil, err
	}

	context := v8go.NewContext(isolate, global)

	if err := url.InjectTo(context); err != nil {
		context.Close()
		isolate.TerminateExecution()
		isolate.Dispose()
		return nil, nil, err
	}

	if err := console.InjectMultipleTo(context,
		console.NewConsole(console.WithOutput(os.Stderr), console.WithMethodName("error")),
		console.NewConsole(console.WithOutput(os.Stderr), console.WithMethodName("warn")),
		console.NewConsole(console.WithOutput(os.Stdout), console.WithMethodName("log")),
	); err != nil {
		context.Close()
		isolate.TerminateExecution()
		isolate.Dispose()
		return nil, nil, err
	}
	return isolate, context, nil
}

func Load() (*VM, error) {
	isolate, context, err := load()
	if err != nil {
		return nil, err
	}
	return &VM{
		isolate: isolate,
		context: context,
	}, nil
}

func Compile(path, code string) (*VM, error) {
	isolate, context, err := load()
	if err != nil {
		return nil, err
	}
	script, err := isolate.CompileUnboundScript(code, path, v8go.CompileOptions{})
	if err != nil {
		return nil, err
	}

	if _, err := script.Run(context); err != nil {
		return nil, err
	}
	return &VM{
		isolate: isolate,
		context: context,
	}, nil
}

func (vm *VM) Script(path, code string) error {
	script, err := vm.isolate.CompileUnboundScript(code, path, v8go.CompileOptions{})
	if err != nil {
		return err
	}

	if _, err := script.Run(vm.context); err != nil {
		return err
	}
	return nil
}

func (vm *VM) Close() error {
	vm.context.Close()
	vm.isolate.TerminateExecution()
	vm.isolate.Dispose()
	return nil
}
