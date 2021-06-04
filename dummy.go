// +build !wasm

package wasmrouter

func AttachRouter(r *Router) {
	r.run()
}

func Link(title, path string) {

}

func (r *Router) getPath() string {
	if r.useForcePATH {
		return r.forcePath
	}
	return "/"
}
