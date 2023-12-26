package main

import (
	"fmt"
	"github.com/Abhishekdx300/dns-resolver/dns"
)

func main() {

	// 2 byte any number  -- Big Endian
	id := uint16(22)
	// 1st case -- recursion bit set
	// flags := uint16(256) // 8th bit set from right to left --> RFC 1035 4.1.1
	// now setting to zero in the next case
	flags := uint16(0)

	queryString := "open.spotify.com"

	query := dns.DNSquery{
		Header: dns.HeaderBuilder(id, flags),
		Body:   dns.BodyBuilder(queryString),
	}

	encodedQuery, err := query.EncodeQuery()
	if err != nil {
		panic(err)
	}

	responseIP, err := dns.Deserialize(encodedQuery, queryString)
	if err != nil {
		panic(err)
	}
	fmt.Println(responseIP)

}
