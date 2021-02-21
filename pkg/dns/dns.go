package dns

import (
	"log"
	"net"

	"golang.org/x/net/dns/dnsmessage"
)

type server struct {
	conn *net.UDPConn
}

// RunDNSServer Main function, run dns server and start loop
func RunDNSServer() {
	srv := server{}
	srv.start()
	srv.listen()
}

const (
	port         = 8090
	buffSize     = 512
	resolverIP   = "1.1.1.1"
	resolverPort = 53
)

func (s *server) sendPacket(message dnsmessage.Message, addr net.UDPAddr) bool {
	packed, err := message.Pack()
	if err != nil {
		log.Println(err)
		return false
	}

	written, err := s.conn.WriteToUDP(packed, &addr)
	if err != nil {
		log.Println(err)
		return false
	}
	log.Printf("Sent %d bytes to %v", written, addr)
	return true
}

func (s *server) read(buf []byte) (net.UDPAddr, bool) {
	_, addr, err := s.conn.ReadFromUDP(buf)
	if err != nil {
		log.Println(err)
		return *addr, false
	}
	return *addr, true
}

func (s *server) start() bool {
	log.Printf("Listening on port %d", port)
	conn, err := net.ListenUDP("udp", &net.UDPAddr{Port: port})
	if err != nil {
		log.Println(err)
		return false
	}
	s.conn = conn
	return true
}

func (s *server) finish() {
	s.conn.Close()
}

func (s *server) listen() {
	cache := make(map[uint16]*net.UDPAddr)
	defer s.finish()

	for {
		buf := make([]byte, buffSize)
		addr, ok := s.read(buf)

		if !ok {
			continue
		}

		var m dnsmessage.Message
		err := m.Unpack(buf)
		if err != nil {
			log.Println(err)
			continue
		}

		if m.Header.Response {
			if len(m.Authorities) != 0 {
				log.Printf("Received response from %v, %+v", resolverIP, m.Authorities[0].Header.Type)
			} else {
				log.Printf("Received response from %v", resolverIP)

			}

			_, ok := cache[m.ID]
			if ok {
				ok := s.sendPacket(m, *cache[m.ID])
				delete(cache, m.ID)
				if !ok {
					continue
				}
			} else {
				log.Printf("No request associated with %d", m.ID)
			}

			continue
		}

		log.Printf("Received from %v, %+v", addr, m.Questions[0].Name)
		resolver := net.UDPAddr{IP: net.ParseIP(resolverIP), Port: resolverPort}
		s.sendPacket(m, resolver)
		cache[m.ID] = &addr
	}
}
