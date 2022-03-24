package main

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/miekg/dns"
)

// Question structure
//
// type Question struct {
// 	Name   string `dns:"cdomain-name"` // "cdomain-name" specifies encoding (and may be compressed)
// 	Qtype  uint16                     // type of the query; like `A` record type but encoded in integer
// 	Qclass uint16                     // class of the query; 'internet' class
// }

// Answer structure
//
// type RR_Header struct {
// 	Name     string `dns:"cdomain-name"`
// 	Rrtype   uint16
// 	Class    uint16
// 	Ttl      uint32
// 	Rdlength uint16 // Length of data after header.
// }

func resolve(name string) net.IP {
	// Root nameserver
	nameserver := net.ParseIP("8.8.8.8")
	for {
		reply := dnsQuery(name, nameserver)
		if ip := getAnswer(reply); ip != nil {
			// Best case: we get an answer to our query and we're done
			return ip
		} else if target := getCNAME(reply); target != "" {
			return resolve(target)
		} else if nsIP := getGlue(reply); nsIP != nil {
			// Second best: we get a "glue record" with the *IP address* of another nameserver to query
			nameserver = nsIP
		} else if domain := getNS(reply); domain != "" {
			// Third best: we get the *domain name* of another nameserver to query, which we can look up the IP for
			nameserver = resolve(domain)
		} else {
			// If there's no A record we just terminate
			fmt.Println("Couldn't resolve " + name + " , no A record found")
			os.Exit(1)
		}
	}
}

func getAnswer(reply *dns.Msg) net.IP {
	for _, record := range reply.Answer {
		if record.Header().Rrtype == dns.TypeA {
			fmt.Println("  ", record)
			return record.(*dns.A).A
		}
	}
	return nil
}

func getGlue(reply *dns.Msg) net.IP {
	for _, record := range reply.Extra {
		if record.Header().Rrtype == dns.TypeA {
			fmt.Println("  ", record)
			return record.(*dns.A).A
		}
	}
	return nil
}

func getCNAME(reply *dns.Msg) string {
	for _, record := range reply.Answer {
		if record.Header().Rrtype == dns.TypeCNAME {
			fmt.Println("  ", record)
			return record.(*dns.CNAME).Target
		}
	}
	return ""
}

func getNS(reply *dns.Msg) string {
	for _, record := range reply.Ns {
		if record.Header().Rrtype == dns.TypeNS {
			fmt.Println("  ", record)
			return record.(*dns.NS).Ns
		}
	}
	return ""
}

func dnsQuery(name string, server net.IP) *dns.Msg {
	msg := new(dns.Msg)
	msg.SetQuestion(name, dns.TypeA)
	c := new(dns.Client)
	reply, _, _ := c.Exchange(msg, server.String()+":53")
	return reply
}

func main() {
	name := os.Args[1]
	if !strings.HasSuffix(name, ".") {
		name = name + "."
	}
	fmt.Println("Result:", resolve(name))
}
