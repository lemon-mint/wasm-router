// +build !wasm

package wasmrouter

func AttachRouter(r *Router) {
}

func (r *Router) Full(path string, h RouterHandler) {

}

func (r *Router) After(path string, h RouterHandler) {

}
