// Copyright 2018 Aleksandr Demakin. All rights reserved.

package server

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/avdva/slot-machine/machine"

	jwt "github.com/dgrijalva/jwt-go"
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

func (s *Server) decoder(_ context.Context, r *http.Request) (interface{}, error) {
	claims := struct {
		jwt.StandardClaims
		request
	}{}
	machineID := api.RequestVars(r)["machine_id"]
	if _, found := s.config.Machines[machineID]; !found {
		return nil, api.NewErrorWithCode("bad machine", http.StatusBadRequest)
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, api.NewErrorWithCode(err.Error(), http.StatusBadRequest)
	}
	token, err := jwt.ParseWithClaims(string(body), &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("hello"), nil
	})
	if !token.Valid || err != nil {
		return nil, api.NewErrorWithCode(err.Error(), http.StatusUnauthorized)
	}
	return claims.request, nil
}

func (s *Server) spin(ctx context.Context, req interface{}) (response interface{}, err error) {
	r := req.(request)
	result := s.config.Machines[r.Machine].Spin(r.Bet)
	return result, nil
}

func buildResponse(results []machine.SpinResult) response {
	var result response
	for _, r := range results {
		result.Total += r.Total
		result.Spins = append(result.Spins, spin{
			Total: r.Total,
			Type:  r.Type,
			Stops: r.Stops,
		})
	}
	return result
}
