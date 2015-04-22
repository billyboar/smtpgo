package smtpgo

import (
	"bufio"
	"fmt"
	"net"
	"regexp"
	"strings"
)

type session struct {
	server     *Server
	conn       net.Conn
	remoteIP   string
	remoteHost string
	remoteName string
	br         *bufio.Reader
	bw         *bufio.Writer
}

var (
	mailFromRE = regexp.MustCompile(`[Ff][Rr][Oo][Mm]:<(.*)>`) //MAIL FROM: <email address>
	rcptToRe   = regexp.MustCompile(`[Tt][Oo]:<(.+)>`)
)

func (server *Server) NewSession(conn net.Conn) (s *session, err error) {
	s = &session{
		server: server,
		conn:   conn,
		br:     bufio.NewReader(conn),
		bw:     bufio.NewWriter(conn),
	}
	return
}

func (s *session) HandleConnection() {
	defer s.conn.Close()
	/*
		var from string
		var to []string
	*/
	//buf := make([]byte, 4096)

	s.remoteIP, _, _ = net.SplitHostPort(s.conn.RemoteAddr().String())
	names, err := net.LookupAddr(s.remoteIP)
	if err == nil && len(names) > 0 {
		s.remoteHost = names[0]
	} else {
		s.remoteHost = "unknown"
	}

	fmt.Println("Connection established from: %v.", s.remoteIP)

	//greet here
	s.Writef("220 %s %s SMTP service ready\n", s.server.Hostname, s.server.Appname)
	for {
		//n, err := s.conn.Read(buf)
		test_x, err := s.br.ReadSlice('\n')
		if err != nil {
			fmt.Println("Error occured while reading: ", err)
			return
		}

		//Print client commands
		//fmt.Println(string(buf[0:n]))

		/*
			if err != nil || n == 0 {
				s.conn.Close()
				break
			}
		*/
		//exclude new line characters
		//line := string(buf[0 : n-2])
		line := string(test_x[0 : len(test_x)-2])
		command, args := ParseCommand(line)
		fmt.Println("command: ", command)
		fmt.Println("line: ", line)

		switch command {

		//Greeting server
		case "EHLO", "HELO":

			s.remoteName = args
			s.Writef("250 %s greets %s\n", s.server.Hostname, s.remoteName)

		//getting receipt's email address
		case "MAIL":

			mailFrom := mailFromRE.FindStringSubmatch(args)
			if mailFrom == nil {
				s.Writef("501 Syntax error in parameters or arguments (invalid FROM parameter)\n")
			} else {
				fmt.Println(mailFrom)
			}

		case "QUIT":
			//break
			return

		default:
			s.Writef("500 Syntax error, command unrecognized\n")
		}
	}
	fmt.Println("Connection from %v is closed.", s.remoteIP)
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
	_, err := s.conn.Write([]byte(fmt.Sprintf(format, args...)))
	if err != nil {
		fmt.Println("Error: ", err)
	}
}
