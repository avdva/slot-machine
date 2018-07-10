// Copyright 2018 Aleksandr Demakin. All rights reserved.

package server

import (
	"context"
	"net/http"

	"github.com/avdva/slot-machine/machine"

	"github.com/space307/go-utils/api"
)

type Config struct {
	Addr     string
	Machines map[string]machine.Interface
}

type Server struct {
	server *api.Server
	config Config
}

func New(config Config) *Server {
	result := &Server{config: config}
	cfg := &api.Config{
		Addr:   config.Addr,
		Prefix: "/api",
		Handlers: []api.PathInfo{
			{
				Path:   "/machines/{machine_id}/spins",
				Enc:    api.EncodeJSONResponse,
				Dec:    result.decoder,
				E:      result.spin,
				Method: "POST",
			},
		},
	}
	result.server = api.NewServer(cfg)
	return result
}

func (s *Server) Serve() error {
	return s.server.Serve()
}

func (s *Server) Stop() error {
	return s.server.Stop()
}

func (s *Server) decoder(_ context.Context, r *http.Request) (request interface{}, err error) {
	machineID := api.RequestVars(r)["machine_id"]
	if _, found := s.config.Machines[machineID]; !found {
		return nil, api.NewErrorWithCode("bad machine", http.StatusBadRequest)
	}
	return machineID, nil
}

func (s *Server) spin(ctx context.Context, request interface{}) (response interface{}, err error) {
	machineID := request.(string)
	s.config.Machines[machineID].Spin(0)
	return nil, nil
}
