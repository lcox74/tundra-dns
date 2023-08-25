package routing

import (
	"fmt"
	"sync"

	"github.com/miekg/dns"
)

func LaunchDNSQueryHandler() {
	// Attach a function to handle DNS requests
	server := &dns.Server{Addr: ":53", Net: "udp"}
	dns.HandleFunc(".", handleDNSRequest)

	// Run the DNS Server
	fmt.Println("Starting DNS Server")
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}

func handleDNSQuery(queryType string, fqdn string) *dns.RR {
	return nil
}

// Testing out go memory pools to see if it helps with performance. Very much
// not needed for this project, but I want to learn how to use them.
var responsePool = sync.Pool{
	New: func() interface{} {
		return &dns.Msg{}
	},
}

func handleDNSRequest(w dns.ResponseWriter, r *dns.Msg) {
	// Check if the message is a query
	if r.MsgHdr.Response {
		return
	}

	if len(r.Question) < 1 {
		return
	}

	// Get a response from the pool
	response := responsePool.Get().(*dns.Msg)
	defer responsePool.Put(response)

	// Clear the response data
	response.Answer = make([]dns.RR, 0)
	response.Ns = nil
	response.Extra = nil

	// Set the response header
	response.SetReply(r)

	// Get the question from the message
	for _, q := range r.Question {
		var rr *dns.RR
		if rr = handleDNSQuery(dns.TypeToString[q.Qtype], q.Name); rr == nil {
			response.SetRcode(r, dns.RcodeNameError)
			break
		}

		// Add the answer to the response
		response.Answer = append(response.Answer, *rr)
	}

	w.WriteMsg(response)
}
