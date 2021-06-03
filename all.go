package wasmrouter

type RouterHandler func(path, url, host, query, hash string) error

type Router struct {
	routes []*route

	NotFoundHandler RouterHandler
}

type route struct {
	startsWith  string
	fullMatch   string
	isFullMatch bool
	ignore      bool
	h           RouterHandler
}
