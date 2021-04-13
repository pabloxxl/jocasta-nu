package dns

import (
	"log"
	"net"

	"github.com/pabloxxl/jocasta-nu/pkg/db"
	"golang.org/x/net/dns/dnsmessage"
)

// Server struct containing all dns connection data
type Server struct {
	conn         *net.UDPConn
	port         int
	resolverIP   string
	resolverPort int
	blockedHosts *[]Record
	buffSize     int
	debug        bool
}

const buffSize = 512

func (s *Server) sendPacket(message dnsmessage.Message, addr net.UDPAddr) bool {
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
	if s.debug {
		log.Printf("Sent %d bytes to %v:%v", written, addr.IP, addr.Port)
	}
	return true
}

func (s *Server) read(buf []byte) (net.UDPAddr, bool) {
	_, addr, err := s.conn.ReadFromUDP(buf)
	if err != nil {
		log.Println(err)
		return *addr, false
	}
	return *addr, true
}

func (s *Server) start() bool {
	log.Printf("Listening on port %d", s.port)
	conn, err := net.ListenUDP("udp", &net.UDPAddr{Port: s.port})
	if err != nil {
		log.Println(err)
		return false
	}
	s.conn = conn
	return true
}

func (s *Server) finish() {
	s.conn.Close()
}

// GetConnection get connection struct filled with preliminary data
func GetConnection(port int, resolverIP string, resolverPort int, blockedHosts *[]Record, debug bool) *Server {
	srv := Server{port: port, resolverIP: resolverIP, resolverPort: resolverPort, blockedHosts: blockedHosts, buffSize: buffSize, debug: debug}
	return &srv
}

// Listen start connection and handle incoming queries
func Listen(s *Server) {
	client := db.CreateClient()
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

		m, data, err := createMessageFromBuffer(buf)
		if err != nil {
			continue
		}

		if data.isResponse {
			if s.debug {
				log.Println(responseToString(data, *s))
			}

			_, ok := cache[data.ID]
			if ok {
				ok := s.sendPacket(m, *cache[data.ID])
				delete(cache, data.ID)
				if !ok {
					log.Printf("Failed to remove record from cache")
					continue
				}
			} else {
				log.Printf("No request associated with %d", data.ID)
			}

			continue
		}

		action := ActionNo
		blocked := false
		logged := false
		for _, question := range data.Questions {
			if s.debug {
				log.Println(questionToString(question, addr))
			} else {
				log.Println(questionToStringShort(question))
			}
			action = GetRecordAction(client, question.URL, *s.blockedHosts)
			blocked = action == ActionBlock
			logged = action == ActionLog
		}
		if !blocked {
			resolver := net.UDPAddr{IP: net.ParseIP(s.resolverIP), Port: s.resolverPort}
			s.sendPacket(m, resolver)
			cache[data.ID] = &addr
		} else {
			s.sendPacket(m, addr)
		}
		if logged {
			// TODO log to database
			log.Printf("%d is logged", m.ID)
		}

		if !data.isResponse {
			putStat(client, ActionToString(action), data, addr.IP, addr.Port)
		}
	}
}
