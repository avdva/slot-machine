// Copyright 2018 Aleksandr Demakin. All rights reserved.

package server

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/avdva/slot-machine/machine"

	jwt "github.com/dgrijalva/jwt-go"
	httptransport "github.com/go-kit/kit/transport/http"
	uuid "github.com/satori/go.uuid"
	"github.com/space307/go-utils/api"
	"github.com/space307/go-utils/errors"
)

var (
	secret = []byte("secret")
)

// Config is a server's config.
type Config struct {
	Addr     string
	Machines map[string]machine.Interface
}

// Server is a REST API server.
type Server struct {
	server *api.Server
	config Config
}

// New returns new Server.
func New(config Config) *Server {
	result := &Server{config: config}
	cfg := &api.Config{
		Addr:   config.Addr,
		Prefix: "/api",
		Handlers: []api.PathInfo{
			{ // for debug purposes.
				Path:   "/machines/{machine_id}/new",
				Enc:    api.EncodeJSONResponse,
				Dec:    result.newDecoder,
				E:      result.new,
				Method: "POST",
				O: []httptransport.ServerOption{
					httptransport.ServerErrorEncoder(errors.CtxEncodeError),
				},
			},
			{
				Path:   "/machines/{machine_id}/spins",
				Enc:    api.EncodeJSONResponse,
				Dec:    result.spinsDecoder,
				E:      result.spin,
				Method: "POST",
				O: []httptransport.ServerOption{
					httptransport.ServerErrorEncoder(errors.CtxEncodeError),
				},
			},
		},
	}
	result.server = api.NewServer(cfg)
	return result
}

// Serve starts http listening.
func (s *Server) Serve() error {
	return s.server.Serve()
}

// Stop stops http listening.
func (s *Server) Stop() error {
	return s.server.Stop()
}

func (s *Server) newDecoder(_ context.Context, r *http.Request) (interface{}, error) {
	machine := api.RequestVars(r)["machine_id"]
	if _, found := s.config.Machines[machine]; !found {
		return nil, errors.NewErrorWithCode("bad machine", "", http.StatusBadRequest)
	}
	return machine, nil
}

func (s *Server) new(ctx context.Context, req interface{}) (response interface{}, err error) {
	machine := req.(string)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid":   machine + "_" + uuid.NewV4().String(),
		"bet":   10,
		"chips": 1000,
	})
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return nil, errors.NewErrorWithCode(err.Error(), "", http.StatusInternalServerError)
	}
	return tokenString, nil
}

func (s *Server) spinsDecoder(_ context.Context, r *http.Request) (interface{}, error) {
	claims := struct {
		jwt.StandardClaims
		request
	}{}
	claims.Machine = api.RequestVars(r)["machine_id"]
	if _, found := s.config.Machines[claims.Machine]; !found {
		return nil, errors.NewErrorWithCode("bad machine", "", http.StatusBadRequest)
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, errors.NewErrorWithCode(err.Error(), "", http.StatusInternalServerError)
	}
	token, err := jwt.ParseWithClaims(string(body), &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})
	if !token.Valid || err != nil {
		return nil, errors.NewErrorWithCode(err.Error(), "", http.StatusUnauthorized)
	}
	return claims.request, nil
}

func (s *Server) spin(ctx context.Context, req interface{}) (response interface{}, err error) {
	r := req.(request)
	m := s.config.Machines[r.Machine]
	wager := m.Wager(r.Bet)
	if wager > r.Chips {
		return nil, errors.NewErrorWithCode("insufficient funds", "", http.StatusBadRequest)
	}
	r.Chips -= wager
	result := m.Spin(r.Bet)
	r.Chips += result.Total
	resp := buildResponse(result)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid":   r.UID,
		"bet":   r.Bet,
		"chips": r.Chips,
	})
	if tokenString, err := token.SignedString(secret); err != nil {
		return nil, errors.NewErrorWithCode(err.Error(), "", http.StatusInternalServerError)
	} else {
		resp.JWT = tokenString
	}
	return resp, nil
}

func buildResponse(result machine.Result) response {
	resp := response{Total: result.Total}
	for _, r := range result.Spins {
		resp.Spins = append(resp.Spins, spin{
			Total: r.Total,
			Type:  r.Type,
			Stops: r.Stops,
		})
	}
	return resp
}
