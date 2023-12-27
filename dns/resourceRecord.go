package dns

import "encoding/binary"

type ResourceRecord struct {
	body   Body
	TTL    uint32
	Length uint16
	Data   []byte
}

func decodeResource(buffer []byte, offset int) (*ResourceRecord, int) {

	name, offset := decodeDomainName(buffer, offset)
	newbody := Body{
		Name:  name,
		Type:  binary.BigEndian.Uint16(buffer[offset : offset+2]),
		Class: binary.BigEndian.Uint16(buffer[offset+2 : offset+4]),
	}
	ttl := binary.BigEndian.Uint32(buffer[offset+4 : offset+8])
	length := binary.BigEndian.Uint16(buffer[offset+8 : offset+10])
	var data []byte
	if newbody.Type == 2 && newbody.Class == 1 {
		newDomainName, _ := decodeDomainName(buffer, offset+10)
		data = []byte(newDomainName)
	} else {
		data = buffer[10+offset : 10+offset+int(length)]
	}

	newResourceBody := ResourceRecord{
		body:   newbody,
		TTL:    ttl,
		Length: length,
		Data:   data,
	}
	return &newResourceBody, offset + 10 + int(length)
}
