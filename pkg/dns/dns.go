package dns

import (
	"log"
	"net"

	"golang.org/x/net/dns/dnsmessage"
)

type server struct {
	conn         *net.UDPConn
	port         int
	resolverIP   string
	resolverPort int
	buffSize     int
}

const buffSize = 512

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
	log.Printf("Listening on port %d", s.port)
	conn, err := net.ListenUDP("udp", &net.UDPAddr{Port: s.port})
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

// GetConnection get connection struct filled with preliminary data
func GetConnection(port int, resolverIP string, resolverPort int) *server {
	srv := server{port: port, resolverIP: resolverIP, resolverPort: resolverPort, buffSize: buffSize}
	return &srv
}

// Listen start connection and handle incoming queries
func Listen(s *server) {
	ok := s.start()
	if !ok {
		return
	}

	cache := make(map[uint16]*net.UDPAddr)
	defer s.finish()

	for {
		buf := make([]byte, s.buffSize)
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
				log.Printf("Received response from %v, %+v", s.resolverIP, m.Authorities[0].Header.Type)
			} else {
				log.Printf("Received response from %v", s.resolverIP)

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
		resolver := net.UDPAddr{IP: net.ParseIP(s.resolverIP), Port: s.resolverPort}
		s.sendPacket(m, resolver)
		cache[m.ID] = &addr
	}
}
