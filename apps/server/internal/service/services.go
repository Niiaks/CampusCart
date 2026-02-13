package service

import "github.com/Niiaks/campusCart/internal/server"

type Service struct {
	server *server.Server
}

func NewServices(s *server.Server) *Service { return &Service{server: s} }
