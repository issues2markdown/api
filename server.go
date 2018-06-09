// Copyright 2018 The issues2markdown Authors. All rights reserved.
//
// Licensed to the Apache Software Foundation (ASF) under one or more
// contributor license agreements.  See the NOTICE file distributed with this
// work for additional information regarding copyright ownership.  The ASF
// licenses this file to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.  See the
// License for the specific language governing permissions and limitations
// under the License.

package api

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Server ...
type Server struct {
	Options ServerOptions
	router  *mux.Router
}

// NewServer ...
func NewServer(options ServerOptions) (*Server, error) {
	server := &Server{
		Options: options,
	}
	server.router = mux.NewRouter()

	server.routes()

	return server, nil
}

// Start ...
func (s *Server) Start() error {
	log.Printf("issues2markdown API listening at %s ...\n", s.Options.Address)

	return http.ListenAndServe(s.Options.Address, s.router)
}
