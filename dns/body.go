package dns

import (
	"encoding/binary"
	"fmt"
	"strings"
)

type Body struct {
	Name  string
	Type  uint16
	Class uint16
}

func (body *Body) encode() ([]byte, error) {
	splitted := strings.Split(body.Name, ".")
	// dots will be replaced by uint8 integers
	// so the effective length of the name will be every character+no. of parts
	strlen := 0
	for _, val := range splitted {
		strlen += len(val)
	}
	strlen += len(splitted) + 1 // len for no. of "." and +1 for extra zero at last

	size := strlen + 4 // 2 for type, 2 for class
	encoded := make([]byte, size)

	// now need to convert the string and append to encoded
	// dns.google.com --> 3dns6google3com0

	ind := 0
	for _, str := range splitted {
		if len(str) > 63 {
			return nil, fmt.Errorf("size cant fit in uint8")
		}
		encoded[ind] = uint8(len(str)) // first the length
		for i := 0; i < len(str); i++ {
			encoded[ind+i+1] = str[i]
		}
		ind += len(str) + 1 // for next parts
	}
	// lastly append 0
	encoded[ind] = 0

	// now append type and class
	binary.BigEndian.PutUint16(encoded[size-4:size-2], body.Type)
	binary.BigEndian.PutUint16(encoded[size-2:size], body.Class)

	return encoded, nil
}

func BodyBuilder(name string) Body {
	Body := Body{
		Name:  name,
		Type:  1,
		Class: 1,
	}
	return Body
}

func decodeBody(buffer []byte, offset int) (*Body, int, error) {
	var name string

	name, offset = decodeDomainName(buffer, offset)

	newQuestionBody := Body{
		Name:  name,
		Type:  binary.BigEndian.Uint16(buffer[offset+1 : offset+3]),
		Class: binary.BigEndian.Uint16(buffer[offset+3 : offset+5]),
	}
	return &newQuestionBody, offset + 5, nil
}

func decodeDomainName(buffer []byte, offset int) (string, int) {
	var name string
	for buffer[offset] != 0 {
		cnt := buffer[offset]

		// first 2 bit is set for a pointer in the response --> which is 192 in int/byte
		//? https://datatracker.ietf.org/doc/html/rfc1035#section-4.1.4
		if cnt == 192 {
			str, _ := decodeDomainName(buffer, int(buffer[offset+1]))
			offset += 2
			name += str
			return name, offset
		} else {
			for i := 1; i <= int(cnt); i++ {
				name = name + string(buffer[offset+i])
			}
			offset = offset + int(cnt) + 1
			name += "."
		}
	}
	// offset pointing at last zero
	// extra "." removed
	return name[:len(name)-1], offset
}
