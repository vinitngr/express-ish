package expressish

import "net/http"

type Ctx struct {
	W      http.ResponseWriter
	R      *http.Request
	params map[string]string
}

func (c *Ctx) Param(key string) string {
	return c.params[key]
}

func (c *Ctx) Text(code int, msg string) {
	c.W.WriteHeader(code)
	c.W.Write([]byte(msg))
}