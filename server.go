package server

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
	if srv.Addr == "" {
		srv.Addr = srv.Port
	}
	if srv.Appname == "" {
		srv.Appname = "SMTPD SERVER"
	}
	if srv.Hostname == "" {
		srv.Hostname, _ = os.Hostname()
	}
}
