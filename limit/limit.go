// Package limit implements a plugin that limits the maximum number of records returned.
package limit

import (
	"context"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/pkg/nonwriter"

	"github.com/miekg/dns"
)

// limit implements limit plugin.
type Limit struct {
	Next  plugin.Handler
	Limit int
}

// ServeDNS implements the plugin.Handler interface.
func (lim Limit) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {

	// Use a nonwriter to capture the response.
	nw := nonwriter.New(w)

	rcode, err := plugin.NextOrFailure(lim.Name(), lim.Next, ctx, nw, r)
	if err != nil {
		// Simply return if there was an error.
		return rcode, err
	}

	lim.limit(nw.Msg)

	// Then write it to the client.
	w.WriteMsg(nw.Msg)
	return rcode, err
}

func (lim *Limit) limit(msg *dns.Msg) {
	// Examine the response and truncate, if required.
	if msg != nil && len(msg.Answer) > lim.Limit {
		msg.Answer = msg.Answer[0:lim.Limit]
	}
}

// Name implements the Handler interface.
func (lim Limit) Name() string { return "limit" }
