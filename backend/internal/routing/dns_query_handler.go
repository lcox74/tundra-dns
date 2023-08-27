package routing

import (
	"fmt"
	"sync"

	"github.com/lcox74/tundra-dns/backend/internal/database"
	"github.com/miekg/dns"
	"github.com/redis/go-redis/v9"
)

func LaunchDNSQueryHandler(rdb *redis.Client) {
	// Attach a function to handle DNS requests
	server := &dns.Server{Addr: ":53", Net: "udp"}
	dns.HandleFunc(".", func(w dns.ResponseWriter, m *dns.Msg) {
		handleDNSRequest(rdb, w, m)
	})

	// Run the DNS Server
	fmt.Println("Starting DNS Server")
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}

// Testing out go memory pools to see if it helps with performance. Very much
// not needed for this project, but I want to learn how to use them.
var responsePool = sync.Pool{
	New: func() interface{} {
		return &dns.Msg{}
	},
}

func handleDNSRequest(rdb *redis.Client, w dns.ResponseWriter, r *dns.Msg) {
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

	// Disable DNSSEC
	response.MsgHdr.RecursionDesired = false
	response.MsgHdr.RecursionAvailable = false

	fmt.Println("Got DNS Request from: ", w.RemoteAddr().String())

	// Get the question from the message
	for _, q := range r.Question {
		fmt.Println("Got DNS Question: ", dns.TypeToString[q.Qtype], q.Name)

		// Get the record from the cache
		record, err := database.FetchRecordCache(rdb, dns.TypeToString[q.Qtype], q.Name)
		if err != nil {
			// Most likely not found, so the record doesn't exist
			response.SetRcode(r, dns.RcodeNameError)
			fmt.Println("Failed to fetch record from cache: ", err)
			break
		}

		// Get the response from the record
		rr := record.GetResponse()
		if rr == nil {
			fmt.Println("Failed to generate response record: ", err)
			response.SetRcode(r, dns.RcodeNameError)
			break
		}

		// Add the answer to the response
		response.Answer = append(response.Answer, rr)
	}

	// Send the response
	fmt.Println("Sending DNS Response")
	err := w.WriteMsg(response)
	if err != nil {
		fmt.Println("Failed to send response: ", err)
	}
}
