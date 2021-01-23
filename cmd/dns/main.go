package main

import (
	"log"
	"net"

	"golang.org/x/net/dns/dnsmessage"
)

type Server struct {
	conn *net.UDPConn
}

type ConnectionData struct {
	port     int
	buffSize int
}

func sendPacket(conn *net.UDPConn, message dnsmessage.Message, addr net.UDPAddr) {
	packed, err := message.Pack()
	if err != nil {
		log.Println(err)
		return
	}

	_, err = conn.WriteToUDP(packed, &addr)
	if err != nil {
		log.Println(err)
	}
}

type Packet struct {
	addr    net.UDPAddr
	message dnsmessage.Message
}

func main() {

	data := ConnectionData{8090, 512}

	var err error
	var conn *net.UDPConn
	log.Print("Listening on port 8090")
	conn, err = net.ListenUDP("udp", &net.UDPAddr{Port: data.port})
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	for {
		buf := make([]byte, data.buffSize)
		_, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Println(err)
			continue
		}
		var m dnsmessage.Message
		err = m.Unpack(buf)
		if err != nil {
			log.Println(err)
			continue
		}

		packet := Packet{*addr, m}
		log.Printf("Received from %v, %+v", addr, packet.message.Questions[0].Name)
	}
}
