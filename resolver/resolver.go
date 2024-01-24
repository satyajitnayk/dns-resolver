package resolver

import (
	"bytes"
	"fmt"
	"math/rand"
	"net"

	"github.com/satyajitnayk/dns-resolver/parser"
)

// root server by Verisign, Inc. https://www.iana.org/domains/root/servers
const rootNameServer = "198.41.0.4"

// build a DNS Query
func BuildQuery(domainName string, recordType uint16) []byte {
	encodedDomainName := parser.EncodeDomainName(domainName)
	id := uint16(rand.Intn(65535))
	//recursionDesired := uint16(1 << 8)
	header := parser.DNSHeader{
		ID:          id,
		Flags:       0, // Set Flags to 0, since we don't need recursion.
		NumQuestion: 1,
	}
	question := parser.DNSQuestion{
		Name:  encodedDomainName,
		Type:  recordType,
		Class: parser.ClassIn,
	}

	var query bytes.Buffer
	query.Write(header.ToBytes())
	query.Write(question.ToBytes())

	return query.Bytes()
}

func SendQuery(ip, domainName string, recordType uint16) parser.DNSPacket {
	query := BuildQuery(domainName, recordType)

	// create a UDP socket
	conn, err := net.Dial("udp", fmt.Sprintf("%s:53", ip))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// send our query to 8.8.8.8, port 53. Port 53 is the DNS port.
	_, err = conn.Write(query)
	if err != nil {
		panic(err)
	}

	// read the response. UDP DNS responses are usually less than 512 bytes
	// so reading 1024 bytes is enough
	response := make([]byte, 1024)
	_, err = conn.Read(response)
	if err != nil {
		panic(err)
	}

	// need to parse or extract information from binary data, and a Reader provides
	// a convenient way to sequentially read data from a byte slice or any other
	// data source that implements the io.Reader interface
	responseHeader := bytes.NewReader(response)

	// process the response
	packet := parser.ParseDNSPacket(responseHeader)

	return packet
}

func Resolve(domainName string, recordTyep uint16) []byte {
	nameServer := rootNameServer
	for {
		fmt.Printf("Querying %s for %s\n", nameServer, domainName)
		responsePacket := SendQuery(nameServer, domainName, recordTyep)
		if ip := responsePacket.GetAnswer(); ip != nil {
			return ip
		}

		if nsIP := responsePacket.GetNameServerIP(); nsIP != nil {
			nameServer = parser.ParseIP(nsIP)
			continue
		}

		if nsDomain := responsePacket.GetNameServer(); nsDomain != "" {
			nameServer = parser.ParseIP(Resolve(nsDomain, parser.TypeA))
		}
	}
}
