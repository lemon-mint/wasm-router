// +build wasm

package wasmrouter

import (
	"log"
	"net/url"
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
	ctrlKey := e.Get("ctrlKey")
	if ctrlKey.IsUndefined() || ctrlKey.Bool() {
		return js.Undefined()
	}

	target := e.Get("target")
	if !target.IsUndefined() {
		name := target.Get("tagName")
		if !name.IsUndefined() {
			if name.String() == "A" {
				href := target.Get("href")
				if !href.IsUndefined() {
					log.Println("Link Click:", href.String())
					location := js.Global().Get("window").Get("location").Get("href").String()
					locationURL, err := url.Parse(location)
					if err != nil {
						log.Println("location:", err)
						return js.Undefined()
					}
					linkURL, err := url.Parse(href.String())
					if err != nil {
						log.Println("link Error:", err)
						return js.Undefined()
					}
					if linkURL.Host == locationURL.Host {
						e.Call("preventDefault")
						Link("Link", href.String())
						log.Println("Link: preventDefault")
					}
				}
			}
		}
	}
	return js.Undefined()
}
