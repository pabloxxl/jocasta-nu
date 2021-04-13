package dns

import (
	"fmt"
	"log"
	"net"
	"strings"

	"golang.org/x/net/dns/dnsmessage"
)

type Question struct {
	URL  string
	Type string
}

type MessageData struct {
	ID           uint16
	Questions    []Question
	isResponse   bool
	ResponseType string
	ResponseURL  string
}

func createMessageFromBuffer(buf []byte) (dnsmessage.Message, MessageData, error) {
	var message dnsmessage.Message
	var data MessageData
	err := message.Unpack(buf)
	if err != nil {
		log.Println(err)
	}

	data.ID = message.ID
	data.isResponse = message.Header.Response
	data.ResponseType = ""
	if data.isResponse && len(message.Authorities) != 0 {
		data.ResponseType = message.Authorities[0].Header.Type.String()
		data.ResponseURL = strings.TrimSuffix(message.Authorities[0].Header.Name.String(), ".")
	}

	for _, question := range message.Questions {
		var parsedQuestion Question
		parsedQuestion.Type = question.Type.String()
		parsedQuestion.URL = strings.TrimSuffix(question.Name.String(), ".")
		data.Questions = append(data.Questions, parsedQuestion)
	}

	return message, data, err
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
