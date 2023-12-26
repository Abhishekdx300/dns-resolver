package dns

type ResourceRecord struct {
	body   Body
	TTL    uint32
	Length uint16
	Data   []byte
}
