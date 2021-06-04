package wasmrouter

import "strings"

type RouterHandler func(path string) error

type Router struct {
	routes []*route

	NotFoundHandler RouterHandler

	useForcePATH bool
	forcePath    string
}

type route struct {
	startsWith  string
	fullMatch   string
	isFullMatch bool
	ignore      bool
	h           RouterHandler
}

func (r *Router) SetForcePATH(path string) {
	r.forcePath = path
	r.useForcePATH = true
}

func (r *Router) run() {
	pathname := r.getPath()
	IsFound := false
	for i := range r.routes {
		if !r.routes[i].ignore {
			if r.routes[i].isFullMatch {
				if r.routes[i].fullMatch == pathname {
					IsFound = true
					r.routes[i].h(pathname)
				}
			} else {
				if strings.HasPrefix(pathname, r.routes[i].startsWith) {
					IsFound = true
					r.routes[i].h(pathname)
				}
			}
		} else {
			return
		}
	}
	if !IsFound {
		r.NotFoundHandler(pathname)
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
