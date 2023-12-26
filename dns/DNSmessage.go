package dns

import (
	"fmt"
	"net"
)

// const serverIP = "8.8.8.8"
// const serverIP = "192.203.230.10"
const serverIP = "198.41.0.4"
const DnsPort = "53"
const bufferSize = 512

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
	header := query.Header.encode()
	body, err := query.Body.encode()
	if err != nil {
		return nil, err
	}
	encodedQuery := append(header, body...)
	return encodedQuery, nil
}

func sendMessage(message []byte, IPaddr string) ([]byte, error) {
	conn, err := net.Dial("udp", fmt.Sprintf("%s:53", IPaddr))
	if err != nil {
		return nil, err
	}
	_, err = conn.Write(message)
	if err != nil {
		return nil, err
	}
	buffer := make([]byte, bufferSize)
	_, err = conn.Read(buffer)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	return buffer, nil
}

func Deserialize(message []byte, queryStr string) (string, error) {

	stack := Stack{}
	stack.Push(serverIP)
	visited := make(map[string]int)
	visited[serverIP] = 1

	for len(stack) > 0 {
		IPaddr, err := stack.Pop()
		fmt.Printf("Queueing for %s to %s \n", queryStr, IPaddr)
		if err != nil {
			return "", err
		}
		buffer, err := sendMessage(message, IPaddr)
		if err != nil {
			return "", err
		}

		header := buffer[0:12]
		decodedHeader := decodeHeader(header)
		//todo need to validate id of header

		// decodedQuestion, offset, err := decodeQuestionBody(buffer, 12) // after the header
		_, offset, err := decodeBody(buffer, 12) // after the header

		if err != nil {
			return "", err
		}

		ansCnt := int(decodedHeader.AnswerRecordCount)
		authCnt := int(decodedHeader.AuthorityRecordCount)
		addCnt := int(decodedHeader.AdditionalRecordCount)
		var decodedAnswer *ResourceRecord
		var decodedAuthority *ResourceRecord
		var decodedAdditional *ResourceRecord

		if ansCnt > 0 {
			decodedAnswer, _ = decodeResource(buffer, offset)
			return fmt.Sprintf("%d.%d.%d.%d", decodedAnswer.Data[0], decodedAnswer.Data[1], decodedAnswer.Data[2], decodedAnswer.Data[3]), nil
		}

		responseAuths := make([]*ResourceRecord, 0)
		for i := 0; i < authCnt; i++ {
			decodedAuthority, offset = decodeResource(buffer, offset)
			responseAuths = append(responseAuths, decodedAuthority)
		}

		responseAddis := make([]*ResourceRecord, 0)

		for i := 0; i < addCnt; i++ {
			decodedAdditional, offset = decodeResource(buffer, offset)
			responseAddis = append(responseAddis, decodedAdditional)
		}

		for _, addi := range responseAddis {

			if addi.body.Type == 1 && addi.body.Class == 1 && addi.Length == 4 {
				newserv := fmt.Sprintf("%d.%d.%d.%d", addi.Data[0], addi.Data[1], addi.Data[2], addi.Data[3])
				_, alreadyVisited := visited[newserv]
				if !alreadyVisited {
					stack.Push(newserv)
					visited[newserv] = 1
				}
			}
		}

		if len(stack) == 0 && len(responseAuths) > 0 {
			for _, auth := range responseAuths {
				newserv := fmt.Sprintf("%d.%d.%d.%d", auth.Data[0], auth.Data[1], auth.Data[2], auth.Data[3])
				_, alreadyVisited := visited[newserv]
				if !alreadyVisited {
					stack.Push(newserv)
					visited[newserv] = 1
				}
			}
		}
	}

	return "", fmt.Errorf("no ip found")
}
