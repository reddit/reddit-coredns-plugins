package limit

import (
	"fmt"
	"net"
	"reflect"
	"testing"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/whoami"

	"github.com/miekg/dns"
)

func TestLimit(t *testing.T) {
	lim := Limit{
		Limit: 10,
	}

	inputAnswer := make([]dns.RR, 254)
	// Generate a records for 0.0.0.0/16. Can generate up to 64516 records with this function.
	for i := 0; i < 254; i++ {
		inputAnswer[i] = &dns.A{
			A: net.ParseIP(fmt.Sprintf("192.168.%v.%v", i/254, i%254)),
		}
	}

	tests := []struct {
		next         plugin.Handler
		qname        string
		limit        int
		inputAnswer  []dns.RR
		outputAnswer []dns.RR
		expectedErr  error
	}{
		// This plugin is responsible for limiting the number of records in outgoing queries.
		// If the number of inputs is < limit, the full set should be returned.
		{
			next:         whoami.Whoami{},
			qname:        ".",
			inputAnswer:  inputAnswer[0:3],
			outputAnswer: inputAnswer[0:3],
			expectedErr:  nil,
		},
		// If the number of inputs is > limit, the first n (limit) should be returned.
		{
			next:         whoami.Whoami{},
			qname:        ".",
			inputAnswer:  inputAnswer,
			outputAnswer: inputAnswer[0:10],
			expectedErr:  nil,
		},
	}

	for i, tc := range tests {
		req := new(dns.Msg)
		req.SetQuestion(dns.Fqdn(tc.qname), dns.TypeA)
		req.Answer = tc.inputAnswer

		lim.limit(req)

		if !reflect.DeepEqual(req.Answer, tc.outputAnswer) {
			t.Errorf("Test %d: Expected answer %v, but got %v", i, tc.outputAnswer, req.Answer)
		}
	}
}
