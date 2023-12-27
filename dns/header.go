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

func (hdr *Header) encode() []byte {
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

	binary.BigEndian.PutUint16(buffer[0:2], hdr.ID)
	binary.BigEndian.PutUint16(buffer[2:4], hdr.Flags)
	binary.BigEndian.PutUint16(buffer[4:6], hdr.BodyEntryCount)
	binary.BigEndian.PutUint16(buffer[6:8], hdr.AnswerRecordCount)
	binary.BigEndian.PutUint16(buffer[8:10], hdr.AuthorityRecordCount)
	binary.BigEndian.PutUint16(buffer[10:12], hdr.AdditionalRecordCount)

	return buffer
}

func headerBuilder(id, flags uint16) Header {
	header := Header{
		ID:                    id,
		Flags:                 flags,
		BodyEntryCount:        1,
		AnswerRecordCount:     0, // rem. 3 not present
		AuthorityRecordCount:  0,
		AdditionalRecordCount: 0,
	}
	return header
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
	return &newheader
}
