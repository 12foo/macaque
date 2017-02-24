package macaque

import (
	"fmt"
	"github.com/mholt/caddy"
	"github.com/mholt/caddy/caddyhttp/httpserver"
	"strconv"
)

func init() {
	caddy.RegisterPlugin("macaque", caddy.Plugin{
		ServerType: "http",
		Action:     setup,
	})
}

func macaqueParseConfig(c *caddy.Controller) (config MacaqueConfig, err error) {
	config = MacaqueConfig{
		SyncInterval: 10,
		PortalAPI:    "",
	}
	config.Policies = make(map[string]MacaquePolicy)
	for c.Next() {
		args := c.RemainingArgs()
		if len(args) != 1 {
			return config, c.Err("Expected a path parameter.")
		}
		config.Path = args[0]
		for c.NextBlock() {
			directive := c.Val()
			switch directive {
			case "policy":
				pol := MacaquePolicy{"", 0, 0, 0, 0}
				args = c.RemainingArgs()
				if len(args) != 4 && len(args) != 7 {
					return config, c.ArgErr()
				}
				pol.Name = args[0]
				for i := 1; i < len(args); i += 3 {
					switch args[i] {
					case "per_ip":
						pol.PerIPMax, err = strconv.Atoi(args[i+1])
						if err != nil {
							return
						}
						pol.PerIPInterval, err = strconv.Atoi(args[i+2])
						if err != nil {
							return
						}
					case "per_key":
						pol.PerKeyMax, err = strconv.Atoi(args[i+1])
						if err != nil {
							return
						}
						pol.PerKeyInterval, err = strconv.Atoi(args[i+2])
						if err != nil {
							return
						}
					default:
						return config, c.Err("Policy quotas must be specified per_ip or per_key.")
					}
				}
				config.Policies[pol.Name] = pol
			case "database":
				args = c.RemainingArgs()
				if len(args) != 1 {
					return config, c.ArgErr()
				}
				config.DBPath = args[0]
			case "sync_interval":
				args = c.RemainingArgs()
				if len(args) != 1 {
					return config, c.ArgErr()
				}
				config.SyncInterval, err = strconv.Atoi(args[0])
				if err != nil {
					return
				}
			case "portal_api":
				args = c.RemainingArgs()
				if len(args) != 1 {
					return config, c.ArgErr()
				}
				config.PortalAPI = args[0]
			default:
				return config, c.Errf("Unknown macaque directive: %s", directive)
			}
		}
	}
	if len(config.Policies) == 0 {
		config.Policies["default"] == MacaquePolicy{"default", 100, 5, 0, 0}
	}
	if config.DBPath == "" {
		return config, c.Err("You must specify a database file path.")
	}
	return config, nil
}

func setup(c *caddy.Controller) error {
	cfg := httpserver.GetConfig(c)

	mcfg, err := macaqueParseConfig(c)
	if err != nil {
		return err
	}

	mq := NewMacaque(mcfg)
	cfg.AddMiddleware(func(next httpserver.Handler) httpserver.Handler {
		mq.Next = next
		return mq
	})

	return nil
}
