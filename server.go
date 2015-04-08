package server

import (
	"log"
	"net"
	"os"
)

type Server struct {
	Addr     string // address of the server
	Appname  string // application name. mentioned in HELO command etc.
	Hostname string //copycat
	Port     string // server port
}

func StartServer(addr string, appname, string, hostname string, port string) error {
	server := &Server{Addr: addr,
		Appname:  appname,
		Hostname: hostname,
		Port:     port}
	return server.Start()
}

func (server *Server) Start() error {
	if server.Addr == "" {
		server.Addr = server.Port
	}
	if server.Appname == "" {
		server.Appname = "SMTPD SERVER"
	}
	if server.Hostname == "" {
		server.Hostname, _ = os.Hostname()
	}
	ln, err := net.Listen("tcp", server.Port)
	if err != nil {
		return err
	}
	return server.Serve(ln)
}

func (server *Server) Serve(ln net.Listener) error {
	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("%s: Error has occured: %v", server.Appname, err)
			continue
		}
		session, err := server.NewSession(conn)
		if err != nil {
			continue
		}
		go session.HandleConnection()
	}
}
