package main

import (
	"fmt"
	"math/rand"

	"github.com/Abhishekdx300/dns-resolver/dns"
)

func random() uint16 {
	// Generate a random uint16
	randomUint16 := uint16(rand.Uint32())
	return randomUint16
}

func main() {
	fmt.Println("Main")

	/*
		var id uint16 = random()
		flags := random()
		// var flags uint16 = random() | (1 << 7)
		flags |= (1 << 7)
	*/
	//* for some reason these isnt working will check later

	// 2 byte any number  -- //!Big Endian - remember
	id := uint16(22)
	// 1st case
	// flags := uint16(256) // 8th bit set from right to left --> RFC 1035 4.1.1
	//* now setting to zero in the next case
	flags := uint16(0)

	query := dns.DNSquery{
		// Header: dns.Header{ID: 22, Flags: 0, BodyEntryCount: 1, AnswerRecordCount: 0, AuthorityRecordCount: 0, AdditionalRecordCount: 0},
		Header: dns.HeaderBuilder(id, flags),
		// Body:   dns.Body{Name: "dns.google.com", Type: 1, Class: 1}, ---> question
		Body: dns.BodyBuilder("dns.google.com"),
	}

	encodedQuery, err := query.EncodeQuery()
	if err != nil {
		panic(err)
	}

	dns.SendMessage(encodedQuery)

}
