package middleware

import (
	"bufio"
	"os"
	"strings"

	expressish "github.com/vinitngr/express-ish"
)

func DotEnv(paths ...string) expressish.Middleware {
	path := ".env"
	if len(paths) > 0 && paths[0] != "" {
		path = paths[0]
	}

	env := map[string]string{}

	file, err := os.Open(path)
	if err == nil {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			kv := strings.SplitN(line, "=", 2)
			if len(kv) == 2 {
				key := strings.TrimSpace(kv[0])
				val := strings.TrimSpace(kv[1])
				env[key] = val
				_ = os.Setenv(key, val)
			}
		}
		file.Close()
	}

	return func(c *expressish.Ctx, next expressish.Next) {
		for k, v := range env {
			c.Set(k, v)
		}
		next()
	}
}
