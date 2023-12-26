package dns

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
)

const DnsPort = "53"

// const serverIP = "8.8.8.8"

const serverIP = "198.41.0.4"

func HeaderBuilder(id, flags uint16) Header {
	header := Header{
		ID: id,
		// flags 7th bit (recursion desired bit) must be set for request
		// Flags:                 RandomNumGen() | (1 << 7),
		Flags:                 flags,
		BodyEntryCount:        1,
		AnswerRecordCount:     0, // remaining 3 are not present --> only header and Body
		AuthorityRecordCount:  0,
		AdditionalRecordCount: 0,
	}
	return header
}

func BodyBuilder(name string) Body {
	Body := Body{
		Name:  name,
		Type:  1,
		Class: 1,
	}
	return Body
}

type DNSmessage struct {
	Header           Header
	Body             Body
	AnswerRecord     []ResourceRecord
	AuthorityRecord  []ResourceRecord
	AdditionalRecord []ResourceRecord
}

type DNSquery struct {
	Header Header
	Body   Body
}

func (query *DNSquery) EncodeQuery() ([]byte, error) {
	header := query.Header.Encode()

	body, err := query.Body.Encode()

	if err != nil {
		return nil, err
	}

	encodedQuery := append(header, body...)
	return encodedQuery, nil
}

func decodeHeader(header []byte) *Header {
	newheader := Header{
		ID:                    binary.BigEndian.Uint16(header[0:2]),
		Flags:                 binary.BigEndian.Uint16(header[2:4]),
		BodyEntryCount:        binary.BigEndian.Uint16(header[4:6]),
		AnswerRecordCount:     binary.BigEndian.Uint16(header[6:8]),
		AuthorityRecordCount:  binary.BigEndian.Uint16(header[8:10]),
		AdditionalRecordCount: binary.BigEndian.Uint16(header[10:12]),
	}
	fmt.Println("header decoded")
	return &newheader
}

func decodeQuestionBody(buffer []byte, offset int) (*Body, int, error) {
	var name string

	name, offset = decodeDomainName(buffer, offset)

	newQuestionBody := Body{
		Name:  name,
		Type:  binary.BigEndian.Uint16(buffer[offset+1 : offset+3]),
		Class: binary.BigEndian.Uint16(buffer[offset+3 : offset+5]),
	}
	fmt.Println("question decoded")
	return &newQuestionBody, offset + 5, nil
}

func decodeDomainName(buffer []byte, offset int) (string, int) {
	var name string
	for buffer[offset] != 0 {
		cnt := buffer[offset]

		for i := 1; i <= int(cnt); i++ {
			name = name + string(buffer[offset+i])
		}
		offset = offset + int(cnt) + 1
		name += "."
	}
	// offset pointing at last zero
	// extra "." removed
	return name[:len(name)-1], offset
}

func decodeCompressedName(buffer []byte, offset int) (string, int) {
	var name string
	//* this is the case for pointer
	if buffer[offset] == 192 {
		name, _ = decodeDomainName(buffer, int(buffer[offset+1])) // offset+1 has the position of the domain name from the start
		offset = offset + 2
	} else {
		//! normal condition ??
		name, offset = decodeDomainName(buffer, offset)
	}
	return name, offset
}

func decodeResource(buffer []byte, offset int) (*ResourceRecord, int) {
	// first 2 bit is set for a pointer in the response , --> which is //*192 in int/byte
	//? https://datatracker.ietf.org/doc/html/rfc1035#section-4.1.4
	name, offset := decodeCompressedName(buffer, offset) // check above
	newbody := Body{
		Name:  name,
		Type:  binary.BigEndian.Uint16(buffer[offset : offset+2]),
		Class: binary.BigEndian.Uint16(buffer[offset+2 : offset+4]),
	}

	length := binary.BigEndian.Uint16(buffer[offset+8 : offset+10]) // <- cut from below

	newResourceBody := ResourceRecord{
		body:   newbody,
		TTL:    binary.BigEndian.Uint32(buffer[offset+4 : offset+8]),
		Length: length,                                    // cut from here ->
		Data:   buffer[offset+10 : offset+10+int(length)], // giving 4 because normally its 4 byte IP address
	}
	// fmt.Println("resource decoded")

	return &newResourceBody, offset + 10 + int(length)
}

func SendMessage(message []byte) {
	conn, err := net.Dial("udp", fmt.Sprintf("%s:53", serverIP))
	if err != nil {
		panic(err)
	}
	_, err = conn.Write(message)
	if err != nil {
		panic(err)
	}
	// now need to read the buffer from the connection
	buffer := make([]byte, 512)

	_, err = conn.Read(buffer)

	if err != nil {
		panic(err)
	}
	defer conn.Close()
	fmt.Println(buffer)
	deserialize(buffer)
}

func deserialize(buffer []byte) {
	// deserialize and print
	header := buffer[0:12]
	decodedHeader := decodeHeader(header)
	decodedQuestion, offset, err := decodeQuestionBody(buffer, 12) // after the header
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("unprocessed data :")
	fmt.Println(decodedHeader)
	fmt.Println(decodedQuestion)
	//todo no. of answers are given in the header section...--- will configure it later, for now doing it manually

	ansCnt := int(decodedHeader.AnswerRecordCount)
	authCnt := int(decodedHeader.AuthorityRecordCount)
	addCnt := int(decodedHeader.AdditionalRecordCount)
	var decodedAnswer *ResourceRecord
	var decodedAuthority *ResourceRecord
	var decodedAdditional *ResourceRecord

	for i := 0; i < ansCnt; i++ {
		fmt.Println("Answer decoded")
		decodedAnswer, offset = decodeResource(buffer, offset)
		fmt.Println(decodedAnswer)

	}

	for i := 0; i < authCnt; i++ {
		fmt.Println("Authority decoded")

		decodedAuthority, offset = decodeResource(buffer, offset)
		fmt.Println(decodedAuthority)
	}

	for i := 0; i < addCnt; i++ {
		fmt.Println("Additional decoded")

		decodedAdditional, offset = decodeResource(buffer, offset)
		fmt.Println(decodedAdditional)
	}

	fmt.Println("not needed")
	fmt.Println(offset)

}
