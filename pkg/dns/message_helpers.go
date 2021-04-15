package dns

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strings"

	"golang.org/x/net/dns/dnsmessage"
)

type Question struct {
	URL  string
	Type dnsmessage.Type
}

type MessageData struct {
	ID           uint16
	Questions    []Question
	isResponse   bool
	ResponseType dnsmessage.Type
	ResponseURL  string
	Data         *dnsmessage.Message
}

func createMessageFromBuffer(buf []byte) (MessageData, error) {
	var dnsMessage dnsmessage.Message
	var message MessageData
	err := dnsMessage.Unpack(buf)
	if err != nil {
		log.Println(err)
	}

	message.ID = dnsMessage.ID
	message.isResponse = dnsMessage.Header.Response
	message.ResponseType = dnsmessage.TypeA
	if message.isResponse && len(dnsMessage.Authorities) != 0 {
		message.ResponseType = dnsMessage.Authorities[0].Header.Type
		message.ResponseURL = strings.TrimSuffix(dnsMessage.Authorities[0].Header.Name.String(), ".")
	}

	for _, question := range dnsMessage.Questions {
		var parsedQuestion Question
		parsedQuestion.Type = question.Type
		parsedQuestion.URL = strings.TrimSuffix(question.Name.String(), ".")
		message.Questions = append(message.Questions, parsedQuestion)
	}

	message.Data = &dnsMessage

	return message, err
}

func responseToString(data MessageData, s Server) string {
	return fmt.Sprintf("Received response from %v for %s, %s", s.resolverIP, data.ResponseURL, data.ResponseType)
}

func questionToString(question Question, addr net.UDPAddr) string {
	return fmt.Sprintf("Received question for %s from %v:%d, %+v", question.URL, addr.IP, addr.Port, question.Type)
}

func questionToStringShort(question Question) string {
	return fmt.Sprintf("Received question for %s", question.URL)
}

func TypeToString(mtype dnsmessage.Type) string {
	switch mtype {
	case dnsmessage.TypeA:
		return "A"
	case dnsmessage.TypeAAAA:
		return "AAAA"
	case dnsmessage.TypeCNAME:
		return "CNAME"
	case dnsmessage.TypeMX:
		return "MX"
	case dnsmessage.TypeSOA:
		return "SOA"
	default:
		log.Fatalf("Invalid type: %v", mtype)
		return ""
	}
}

func StringToType(mtype string) (dnsmessage.Type, error) {
	switch mtype {
	case "UNKNOWN":
		return dnsmessage.TypeALL, nil
	case "A":
		return dnsmessage.TypeA, nil
	case "AAAA":
		return dnsmessage.TypeAAAA, nil
	case "CNAME":
		return dnsmessage.TypeCNAME, nil
	case "MX":
		return dnsmessage.TypeMX, nil
	case "SOA":
		return dnsmessage.TypeSOA, nil
	default:
		return dnsmessage.TypeA, errors.New("unsupported type")
	}
}
