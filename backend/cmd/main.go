package main

import (
	"fmt"
	"net"
	"sync"

	"github.com/miekg/dns"
)

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

	// Get the question from the message
	for _, q := range r.Question {
		fmt.Printf("Received query for %s type %s and class %s\n", q.Name, dns.TypeToString[q.Qtype], dns.ClassToString[q.Qclass])
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
	response.Answer = append(response.Answer, &dns.A{
		Hdr: dns.RR_Header{
			Name:   "example.com.",
			Rrtype: dns.TypeA,
			Class:  dns.ClassINET,
			Ttl:    300,
		},
		A: net.ParseIP("10.10.10.10"),
	})

	w.WriteMsg(response)
}

func main() {
	fmt.Println("Hello, World!")

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
