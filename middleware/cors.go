package middleware

import (
	expressish "github.com/vinitngr/express-ish"
)

func CORS(allowedOrigins []string) expressish.Middleware {
	allowAll := len(allowedOrigins) == 0
	allowed := make(map[string]struct{}, len(allowedOrigins))

	for _, o := range allowedOrigins {
		allowed[o] = struct{}{}
	}

	return func(c *expressish.Ctx, next expressish.Next) {
		origin := c.ReqHeader("Origin")

		if allowAll {
			c.Header("Access-Control-Allow-Origin", "*")
		} else if _, ok := allowed[origin]; ok {
			c.Header("Access-Control-Allow-Origin", origin)
		}

		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request().Method == "OPTIONS" {
			c.Status(204)
			return
		}

		next()
	}
}
