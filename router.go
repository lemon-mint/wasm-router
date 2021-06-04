// +build wasm

package wasmrouter

import (
	"syscall/js"
)

func AttachRouter(r *Router) {
	window := js.Global().Get("window")
	js.Global().Set("g_rf_onpopstate", js.FuncOf(rf(r)))
	window.Set("onpopstate", js.Global().Get("g_rf_onpopstate"))
	js.Global().Set("g_link", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		Link("", args[0].String())
		return js.Undefined()
	}))
	js.Global().Call("g_rf_onpopstate")
}

func Link(title, path string) {
	js.Global().Get("window").Get("history").Call("pushState", js.Undefined(), title, path)
	js.Global().Call("g_rf_onpopstate")
}

func rf(r *Router) func(this js.Value, args []js.Value) interface{} {
	return func(this js.Value, args []js.Value) interface{} {
		r.run()
		return 0
	}
}

func (r *Router) getPath() string {
	if r.useForcePATH {
		return r.forcePath
	}
	location := js.Global().Get("window").Get("location")
	pathname := location.Get("pathname").String()
	return pathname
}
