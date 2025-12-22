package middleware

import expressish "github.com/vinitngr/express-ish"

func Recovery(c *expressish.Ctx, next expressish.Next) {
	defer func() {
		if r := recover(); r != nil {
			c.Text(500, "internal server error")
		}
	}()
	next()
}