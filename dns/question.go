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

func (body *Body) Encode() ([]byte, error) {
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
