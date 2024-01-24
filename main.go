package main

import (
	"fmt"
	"github.com/satyajitnayk/dns-resolver/parser"
	"github.com/satyajitnayk/dns-resolver/resolver"
)

func main() {
	domainName := "google.com"
	recordType := parser.TypeA
	ipData := resolver.Resolve(domainName, recordType)
	fmt.Printf("Resolved IP for %s is %s\n", domainName, parser.ParseIP(ipData))
}
