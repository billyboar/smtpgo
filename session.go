package smtpgo

import (
	"fmt"
	"net"
	"strings"
)

type session struct {
	server     *Server
	conn       net.Conn
	remoteIP   string
	remoteHost string
	remoteName string
}

func (server *Server) NewSession(conn net.Conn) (s *session, err error) {
	s = &session{
		server: server,
		conn:   conn,
	}
	return
}

func (s *session) HandleConnection() {
	defer s.conn.Close()
	/*
		var from string
		var to []string
	*/
	buf := make([]byte, 4096)

	s.remoteIP, _, _ = net.SplitHostPort(s.conn.RemoteAddr().String())
	names, err := net.LookupAddr(s.remoteIP)
	if err == nil && len(names) > 0 {
		s.remoteHost = names[0]
	} else {
		s.remoteHost = "unknown"
	}

	//greet here
	_, err = s.conn.Write([]byte(fmt.Sprintf("220 %s %s SMTP service ready", s.server.Hostname, s.server.Appname)))
	for {
		n, err := s.conn.Read(buf)

		//Print client commands
		fmt.Println(string(buf[0:n]))
		if err != nil || n == 0 {
			s.conn.Close()
			break
		}
		//exclude new line characters
		line := string(buf[0 : n-2])
		command, args := ParseCommand(line)

		switch command {
		case "EHLO", "HELO":
			s.remoteName = args
			//server.conn.
			s.Writef("250 %s greets %s", s.server.Hostname, s.remoteName)
		}
	}
}

func ParseCommand(line string) (command string, args string) {
	if i := strings.Index(line, " "); i != -1 {
		command = strings.ToUpper(line[:i])
		args = strings.TrimSpace(line[i+1 : len(line)])
	} else {
		command = strings.ToUpper(line)
		args = ""
	}
	return command, args
}

func (s *session) Writef(format string, args ...interface{}) {
	s.conn.Write([]byte(fmt.Sprintf(format, args...)))
}
