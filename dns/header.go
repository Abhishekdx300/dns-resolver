package dns

import "encoding/binary"

type Header struct {
	ID                    uint16
	Flags                 uint16
	BodyEntryCount        uint16
	AnswerRecordCount     uint16
	AuthorityRecordCount  uint16
	AdditionalRecordCount uint16
}

func (hdr *Header) Encode() []byte {
	/*
		ID -- 2 Byte
		Flags -- 2 Byte
		Q. Cnt -- 2 Byte
		A. Cnt -- 2 Byte
		Auth. Cnt --- 2 Byte
		Addi. Cnt -- 2 Byte
		===12 Bytes===
	*/
	buffer := make([]byte, 12)

	binary.BigEndian.PutUint16(buffer[:2], hdr.ID)
	binary.BigEndian.PutUint16(buffer[2:4], hdr.Flags)
	binary.BigEndian.PutUint16(buffer[4:6], hdr.BodyEntryCount)
	binary.BigEndian.PutUint16(buffer[6:8], hdr.AnswerRecordCount)
	binary.BigEndian.PutUint16(buffer[8:10], hdr.AuthorityRecordCount)
	binary.BigEndian.PutUint16(buffer[10:12], hdr.AdditionalRecordCount)

	return buffer
}
