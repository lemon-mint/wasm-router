// +build wasm

package wasmrouter

import (
	"log"
	"syscall/js"
)

func AttachRouter(r *Router, SSR bool) {
	window := js.Global().Get("window")
	js.Global().Set("g_rf_onpopstate", js.FuncOf(rf(r)))
	window.Set("onpopstate", js.Global().Get("g_rf_onpopstate"))
	js.Global().Set("g_link", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		Link("", args[0].String())
		return js.Undefined()
	}))
	if !SSR {
		js.Global().Call("g_rf_onpopstate")
	}
	window.Call("addEventListener", "click", js.FuncOf(onClick))
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

func onClick(this js.Value, args []js.Value) interface{} {
	e := args[0]
	target := e.Get("target")
	if !target.IsUndefined() {
		name := target.Get("tagName")
		if !name.IsUndefined() {
			if name.String() == "A" {
				href := target.Get("href")
				if !href.IsUndefined() {
					log.Println("Link Click:", href.String())
				}
			}
		}
	}
	return js.Undefined()
}
