package main

import (
	"fmt"

	"github.com/Abhishekdx300/dns-resolver/dns"
)

func main() {

	// queryString := "cses.fi"
	// queryString := "sync-write.vercel.app"

	var queryString string
	fmt.Println("Enter url:")
	fmt.Scanln(&queryString)

	resp := dns.Initialize(queryString, 22)
	fmt.Println()
	fmt.Printf("The resolved IP address is %s", resp)
}
