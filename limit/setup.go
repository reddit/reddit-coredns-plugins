package limit

import (
	"github.com/coredns/caddy"
	"strconv"

	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"

	clog "github.com/coredns/coredns/plugin/pkg/log"
)

const pluginName = "limit"

var log = clog.NewWithPlugin(pluginName)

func init() { plugin.Register(pluginName, setup) }

func setup(c *caddy.Controller) error {
	limit, err := parse(c)
	if err != nil {
		return plugin.Error(pluginName, err)
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		return Limit{Next: next, Limit: limit}
	})

	return nil
}

func parse(c *caddy.Controller) (int, error) {
	for c.Next() {
		args := c.RemainingArgs()
		switch len(args) {
		case 1:
			// Specified value is needed to verify
			limit, err := strconv.Atoi(args[0])
			if err != nil {
				return -1, plugin.Error(pluginName, c.ArgErr())
			}
			if limit < 1 {
				return -1, plugin.Error(pluginName, c.ArgErr())
			}
			return limit, nil
		default:
			// Only 1 argument is acceptable
			return -1, plugin.Error(pluginName, c.ArgErr())
		}
	}
	return -1, plugin.Error(pluginName, c.ArgErr())
}
