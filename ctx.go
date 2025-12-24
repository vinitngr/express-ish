package expressish

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Ctx struct {
	w http.ResponseWriter
	r *http.Request

	params map[string]string
	store  map[string]any
	query  url.Values
	body   map[string]any

	wroteHeader bool
}

func NewCtx(w http.ResponseWriter, r *http.Request, params map[string]string) *Ctx {
	return &Ctx{w: w, r: r, params: params}
}

func (c *Ctx) Request() *http.Request {
	return c.r
}

func (c *Ctx) Param(key string) string {
	return c.params[key]
}

func (c *Ctx) Query(key string) string {
	if c.query == nil {
		c.query = c.r.URL.Query()
	}
	return c.query.Get(key)
}

func (c *Ctx) QueryInt(key string, def int) int {
	v, err := strconv.Atoi(c.Query(key))
	if err != nil {
		return def
	}
	return v
}

func (c *Ctx) QueryBool(key string, def bool) bool {
	v, err := strconv.ParseBool(c.Query(key))
	if err != nil {
		return def
	}
	return v
}

func (c *Ctx) Set(key string, val any) {
	if c.store == nil {
		c.store = make(map[string]any)
	}
	c.store[key] = val
}

func (c *Ctx) Get(key string) (any, bool) {
	if c.store == nil {
		return nil, false
	}
	v, ok := c.store[key]
	return v, ok
}


func (c *Ctx) Header(key, val string) {
	c.w.Header().Set(key, val)
}

func (c *Ctx) ReqHeader(key string) string {
	return c.r.Header.Get(key)
}

func (c *Ctx) Type(ct string) {
	c.Header("Content-Type", ct)
}

func (c *Ctx) Status(code int) *Ctx {
	if !c.wroteHeader {
		c.w.WriteHeader(code)
		c.wroteHeader = true
	}
	return c
}
func (c *Ctx) Text(code int, msg string) {
	c.Type("text/plain; charset=utf-8")
	c.Status(code)
	_, _ = c.w.Write([]byte(msg))
}

func (c *Ctx) JSON(code int, v any) {
	c.Type("application/json")
	c.Status(code)
	_ = json.NewEncoder(c.w).Encode(v)
}

func (c *Ctx) Send(v any) {
	switch t := v.(type) {
	case string:
		c.Text(http.StatusOK, t)
	default:
		c.JSON(http.StatusOK, t)
	}
}

func (c *Ctx) Redirect(url string, code ...int) {
	status := http.StatusMovedPermanently
	if len(code) > 0 {
		status = code[0]
	}
	http.Redirect(c.w, c.r, url, status)
}

func (c *Ctx) Body(dst any) error {
	if c.r.Body == nil {
		return nil
	}

	ct := c.r.Header.Get("Content-Type")

	switch {
	case strings.HasPrefix(ct, "application/json"):
		dec := json.NewDecoder(c.r.Body)
		dec.DisallowUnknownFields()
		return dec.Decode(dst)

	case strings.HasPrefix(ct, "application/x-www-form-urlencoded"):
		if err := c.r.ParseForm(); err != nil {
			return err
		}

		b, err := json.Marshal(c.r.PostForm)
		if err != nil {
			return err
		}
		return json.Unmarshal(b, dst)
	}

	return fmt.Errorf("unsupported content-type: %s", ct)
}

func (c *Ctx) Cookie(cookie *http.Cookie) *Ctx {
	http.SetCookie(c.w, cookie)
	return c
}