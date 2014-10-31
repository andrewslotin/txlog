package txlog

type Context map[string]interface{}

func NewContext() Context {
	return Context(make(map[string]interface{}))
}

func (dst Context) Update(src Context) {
	for k, v := range src {
		dst[k] = v
	}
}
