// +build !wasm

package wasmrouter

func AttachRouter(r *Router) {
	r.run()
}

func Link(title, path string) {
	
}
