package link

import (
	"github.com/gg-tools/portmap/util"
	"log"
	"net"
	"strings"
	"time"
)

//TOKEN
const (
	TokenLen      = 4
	C2PConnect    = "C2P0"
	C2PSession    = "C2P1"
	C2PKeepAlive  = "C2P2"
	P2CNewSession = "P2C1"
	SEPS          = "\n"
)

type OnConnectFunc func(net.Conn)

func Listen(port string, onConnect OnConnectFunc) {
	server, err := net.Listen("tcp", net.JoinHostPort("0.0.0.0", port))
	if err != nil {
		log.Fatal("CAN'T LISTEN: ", err)
	}
	log.Println("listen port:", port)

	for {
		conn, err := server.Accept()
		if err != nil {
			log.Println("Can't Accept: ", err)
			continue
		}
		go onConnect(conn)
	}
}

type Server struct {
	pwd         string
	controlLink *controlLink
	dataLink    dataLink
}

func NewServer(pwd string) *Server {
	return &Server{
		pwd:         pwd,
		controlLink: nil,
		dataLink:    dataLink{},
	}
}

func (s *Server) OnConnect(conn net.Conn) {
	_ = conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	msg, err := util.ReadString(conn)
	_ = conn.SetReadDeadline(time.Time{})
	if err != nil {
		log.Println("can't read: ", err)
		return
	}

	msgs := strings.Split(msg, "\n")
	p := msgs[0]
	if p != s.pwd {
		log.Println("bad password")
		return
	}

	token := msgs[1]
	if token == util.C2PConnect {
		// 主连接
		s.controlLink = &controlLink{conn}
		go s.controlLink.keepAlive()
		return
	} else if token == util.C2PSession {
		// 数据连接
		s.dataLink.connCh <- conn
		return
	}
}

func (s *Server) OnUserConnect(conn net.Conn) {
	defer util.CloseConn(conn)
	if s.controlLink == nil {
		_, _ = conn.Write([]byte("no service"))
		return
	}

	if _, err := util.WriteString(s.controlLink.conn, util.P2CNewSession); err != nil {
		_, _ = conn.Write([]byte("SERVICE FAIL"))
		return
	}

	dataConn := s.dataLink.dataLink()
	if dataConn == nil {
		return
	}

	log.Println("Transfer...")
	go util.CopyFromTo(conn, dataConn, nil)
	go util.CopyFromTo(dataConn, conn, nil)
}

type controlLink struct {
	conn net.Conn
}

func (c *controlLink) keepAlive() {
	defer util.CloseConn(c.conn)
	for {
		_, err := util.ReadString(c.conn)
		if err != nil {
			log.Println("UNREG SERVICE")
			break
		}
		// 处理心跳
	}
}

type dataLink struct {
	connCh chan net.Conn
}

func (l *dataLink) dataLink() net.Conn {
	var conn net.Conn = nil
	select {
	case conn = <-l.connCh:
	case <-time.After(time.Second * 5):
		conn = nil
	}
	return conn
}
