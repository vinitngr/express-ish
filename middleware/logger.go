package middleware

import (
	"fmt"
	"time"

	expressish "github.com/vinitngr/express-ish"
)

func Logger(c *expressish.Ctx, next expressish.Next) {
	start := time.Now()

	req := c.Request()
	method := req.Method
	path := req.URL.Path

	next()

	fmt.Printf(
		"[%s] %s %s %v\n",
		time.Now().Format("15:04:05.000"),
		method,
		path,
		time.Since(start),
	)
}
