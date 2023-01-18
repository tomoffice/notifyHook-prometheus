package api

import (
	"fmt"
	"net/http"
)

type Server struct {
	ip     string
	port   int
	routes *http.ServeMux
}

func (s *Server) Config(ip string, port int) {
	s.ip = ip
	s.port = port
}

func (s *Server) AddHandle(pattern string, handlefunc http.HandlerFunc) {
	s.routes.Handle(pattern, handlefunc)
}

func (s Server) Run() {
	addr := fmt.Sprintf("%s:%d", s.ip, s.port)
	http.ListenAndServe(addr, s.routes)
}
func New() *Server {
	var routes = http.NewServeMux()
	return &Server{
		routes: routes,
	}
}
