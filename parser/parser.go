package parser

import (
	"bytes"
	"encoding/binary"
	"io"
	"strconv"
	"strings"
)

const (
	TypeA   uint16 = 1 // TypeA is the infamous A record type
	TypeNS  uint16 = 2 // TypesNS is Nameserver type
	ClassIn        = 1
)

type DNSRecord struct {
	// Name is the domain name
	Name []byte
	// Type is the record type,ex: A, AAAA, MX
	Type  uint16
	Class uint16
	TTL   uint32
	Data  []byte
}

func parseHeader(reader *bytes.Reader) DNSHeader {
	var header DNSHeader

	err := binary.Read(reader, binary.BigEndian, &header)
	if err != nil {
		panic(err)
	}

	return header
}

func parseQuestion(reader *bytes.Reader) DNSQuestion {

	// parse question type and class
	var typeAndClass struct {
		Type  uint16
		Class uint16
	}

	name := decodeName(reader)

	err := binary.Read(reader, binary.BigEndian, &typeAndClass)
	if err != nil {
		panic(err)
	}

	return DNSQuestion{
		Name:  name,
		Type:  typeAndClass.Type,
		Class: typeAndClass.Class,
	}
}

func decodeName(reader *bytes.Reader) []byte {
	var (
		length byte
		name   []byte
	)

	// parse domain name.
	for {
		length, _ = reader.ReadByte()
		if length == 0 {
			break
		}

		// the two most significant bits of length are set, indicates a compressed name
		if length&0b1100_0000 != 0 {
			// length indicates a compressed name, call a function to decode it.
			name = append(name, decodeCompressedName(length, reader)...)
			// add a period ('.) to the domain name after adding this domain name part
			name = append(name, 46)
			break
		} else {
			// no compressed name read specified number of bytes
			// find a way to read multiple bytes at once, possible reader.ReadAt()
			part := make([]byte, length)
			_, err := io.ReadFull(reader, part)
			if err != nil {
				return nil
			}
			name = append(name, part...)
			// add period('.') to the domain name after adding this domain name part
			name = append(name, 46)
		}
	}
	// strip the extra period('.') from the end
	name = name[0 : len(name)-1]

	return name
}

func decodeCompressedName(length byte, reader *bytes.Reader) []byte {
	// Read the next byte from the reader.
	b, _ := reader.ReadByte()

	// create a 2 byte slice containing the last 6 bits of 'length' and the entire byte 'b'.
	// This slice represents the pointer in compressed DNS format.
	pointerBytes := []byte{length & 0b0011_1111, b}

	// convert the 2byte slice to a 16 bit  unisgned int in big-endian order
	// This gives the offset to the compressed name in the DNS message.
	pointer := binary.BigEndian.Uint16(pointerBytes)

	// get the current position of the reader
	// This is done to later restore the reader to its original position.
	currentPos, _ := reader.Seek(0, io.SeekCurrent)

	// Moves the reader to the position indicated by the pointer,
	// effectively jumping to the location of the compressed name.
	_, err := reader.Seek(int64(pointer), io.SeekStart)
	if err != nil {
		return nil
	}

	// recursively decode the domain name from the new position
	result := decodeName(reader)

	// restore the original position of the reader
	_, err = reader.Seek(currentPos, io.SeekStart)
	if err != nil {
		return nil
	}

	return result
}

func parseRecord(reader *bytes.Reader) DNSRecord {
	name := decodeName(reader)

	var recordData struct {
		Type    uint16
		Class   uint16
		TTL     uint32
		DataLen uint16
	}

	err := binary.Read(reader, binary.BigEndian, &recordData)
	if err != nil {
		panic(err)
	}

	// A DNS record contains various fields, such as the name, type, class,
	// TTL (Time to Live), and data. The data part can vary based on the type of DNS record.
	var record = DNSRecord{
		Name:  name,
		Type:  recordData.Type,
		Class: recordData.Class,
		TTL:   recordData.TTL,
	}

	// Depending on the type of DNS record, parse and set the data field.
	switch recordData.Type {
	//case TypeA:
	case TypeNS: // record type is name server we keep on decoding the name
		data := decodeName(reader)
		record.Data = data
	default:
		var data = make([]byte, recordData.DataLen)

		_, err = io.ReadFull(reader, data)
		if err != nil {
			panic(err)
		}

		record.Data = data
	}

	return record
}

// ParseIP takes a byte slice data representing an IPv4 address and
// converts it into a string representation
func ParseIP(data []byte) string {
	// Initialize a strings.Builder to efficiently build the IP address string.
	var ip strings.Builder

	// Iterate over the first three bytes of the data.
	for _, b := range data[0:3] {
		// each byte is an IP segment, convert it to string and write to ip
		ip.WriteString(strconv.Itoa(int(b)))
		ip.WriteString(".")
	}

	ip.WriteString(strconv.Itoa(int(data[3])))
	return ip.String()
}
