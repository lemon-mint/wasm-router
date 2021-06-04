// +build wasm

package wasmrouter

import (
	"strings"
	"syscall/js"
)

func AttachRouter(r *Router) {
	window := js.Global().Get("window")
	js.Global().Set("g_rf_onpopstate", js.FuncOf(rf(r)))
	window.Set("onpopstate", js.Global().Get("g_rf_onpopstate"))
	js.Global().Call("g_rf_onpopstate")
	js.Global().Set("g_link", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		Link("", args[0].String())
		return js.Undefined()
	}))
}

func Link(title, path string) {
	js.Global().Get("window").Get("history").Call("pushState", js.Undefined(), title, path)
	js.Global().Call("g_rf_onpopstate")
}

func rf(r *Router) func(this js.Value, args []js.Value) interface{} {
	return func(this js.Value, args []js.Value) interface{} {
		location := js.Global().Get("window").Get("location")
		pathname := location.Get("pathname").String()
		IsFound := false
		for i := range r.routes {
			if !r.routes[i].ignore {
				if r.routes[i].isFullMatch {
					if r.routes[i].fullMatch == pathname {
						IsFound = true
						r.routes[i].h(pathname, location.Get("href").String(), location.Get("host").String(), location.Get("search").String(), location.Get("hash").String())
					}
				} else {
					if strings.HasPrefix(pathname, r.routes[i].startsWith) {
						IsFound = true
						r.routes[i].h(pathname, location.Get("href").String(), location.Get("host").String(), location.Get("search").String(), location.Get("hash").String())
					}
				}
			} else {
				return 1
			}
		}
		if !IsFound {
			r.NotFoundHandler(pathname, location.Get("href").String(), location.Get("host").String(), location.Get("search").String(), location.Get("hash").String())
		}
		return 1
	}
}

func (r *Router) Full(path string, h RouterHandler) {
	r.routes = append(r.routes, &route{
		fullMatch:   path,
		isFullMatch: true,
		h:           h,
	})
}

func (r *Router) After(path string, h RouterHandler) {
	r.routes = append(r.routes, &route{
		startsWith:  path,
		isFullMatch: false,
		h:           h,
	})
}
